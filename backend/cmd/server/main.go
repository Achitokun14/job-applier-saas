package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hibiken/asynq"
	"github.com/stripe/stripe-go/v81"

	"job-applier-backend/internal/auth"
	"job-applier-backend/internal/cache"
	"job-applier-backend/internal/config"
	"job-applier-backend/internal/database"
	apperrors "job-applier-backend/internal/errors"
	"job-applier-backend/internal/handlers"
	"job-applier-backend/internal/logger"
	"job-applier-backend/internal/metrics"
	"job-applier-backend/internal/middleware"
	"job-applier-backend/internal/repository"
	"job-applier-backend/internal/services"
	"job-applier-backend/internal/tasks"
)

func main() {
	cfg := config.Load()

	// --- Structured logger ---
	lg := logger.New(cfg.AppEnv)

	// --- Sentry error tracking ---
	if cfg.SentryDSN != "" {
		if err := apperrors.InitSentry(cfg.SentryDSN, cfg.AppEnv); err != nil {
			lg.Fatal().Err(err).Msg("Failed to initialize Sentry")
		}
		defer apperrors.FlushSentry()
		lg.Info().Msg("Sentry initialized")
	}

	// Initialize Stripe
	if cfg.StripeSecretKey != "" {
		stripe.Key = cfg.StripeSecretKey
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		lg.Fatal().Err(err).Msg("Failed to connect to database")
	}

	if err := database.AutoMigrate(db, cfg.AppEnv); err != nil {
		lg.Fatal().Err(err).Msg("Failed to run migrations")
	}

	// Set up Asynq (Redis-based task queue)
	redisOpt := asynq.RedisClientOpt{Addr: cfg.RedisURL}
	asynqClient := asynq.NewClient(redisOpt)
	defer asynqClient.Close()

	// Initialize token blacklist with Redis
	blacklist := auth.NewTokenBlacklist(cfg.RedisURL)

	// Initialize Redis cache
	redisCache := cache.NewRedisCache(cfg.RedisURL)

	// Initialize repositories
	jobRepo := repository.NewJobRepository(db, redisCache)
	appRepo := repository.NewApplicationRepository(db, redisCache)
	userRepo := repository.NewUserRepository(db, redisCache)
	settingsRepo := repository.NewSettingsRepository(db, redisCache)

	h := handlers.New(db, cfg, asynqClient, blacklist, jobRepo, appRepo, userRepo, settingsRepo)

	// Set up Asynq worker server
	pythonClient := services.NewPythonClient(cfg.PythonServiceURL)

	asynqServer := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})

	asynqMux := asynq.NewServeMux()
	asynqMux.Handle(tasks.TypeResumeGenerate, tasks.NewResumeHandler(db, pythonClient))
	asynqMux.Handle(tasks.TypeCoverLetterGenerate, tasks.NewCoverLetterHandler(db, pythonClient))
	asynqMux.Handle(tasks.TypeScrapeJobs, tasks.NewScrapeHandler(db, pythonClient))
	asynqMux.Handle(tasks.TypeAutoApply, tasks.NewAutoApplyHandler(db, pythonClient))

	go func() {
		lg.Info().Msg("Starting Asynq worker server...")
		if err := asynqServer.Run(asynqMux); err != nil {
			lg.Error().Err(err).Msg("Asynq worker server error")
		}
	}()

	r := chi.NewRouter()

	// --- Middleware stack ---
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(logger.RequestLogger(lg))   // structured request logging (replaces chimiddleware.Logger)
	r.Use(apperrors.SentryMiddleware()) // capture panics to Sentry
	r.Use(chimiddleware.Recoverer)     // recover from panics after Sentry sees them
	r.Use(metrics.MetricsMiddleware()) // Prometheus request metrics
	r.Use(middleware.GlobalRateLimit())
	r.Use(middleware.SecurityHeaders(cfg))

	// CORS configuration: use configured origins, fall back to ["*"] for dev
	allowedOrigins := cfg.CORSAllowedOrigins
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// --- Public observability endpoints (no auth) ---
	r.Get("/health", handlers.HealthCheck(db, cfg.RedisURL, cfg.PythonServiceURL))
	r.Handle("/metrics", metrics.MetricsHandler())

	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes (no JWT required) with stricter rate limiting
		r.With(middleware.AuthRateLimit()).Post("/auth/register", h.Register)
		r.With(middleware.AuthRateLimit()).Post("/auth/login", h.Login)
		r.With(middleware.AuthRateLimit()).Post("/auth/refresh", h.RefreshTokenHandler)
		r.With(middleware.AuthRateLimit()).Post("/auth/forgot-password", h.RequestPasswordReset)
		r.With(middleware.AuthRateLimit()).Post("/auth/reset-password", h.ResetPassword)

		// Stripe webhook -- PUBLIC (no auth middleware). Stripe sends webhooks without JWT.
		r.Post("/payments/webhook", h.HandleStripeWebhook)

		// Protected routes (requires valid JWT)
		r.Group(func(r chi.Router) {
			r.Use(h.AuthMiddleware)
			r.Use(middleware.SubscriptionMiddleware(db, redisCache))

			r.Post("/auth/logout", h.Logout)
			r.Get("/jobs", h.SearchJobs)
			r.Post("/jobs/{id}/apply", h.ApplyJob)
			r.Get("/applications", h.ListApplications)
			r.Get("/applications/{id}", h.GetApplication)
			r.Delete("/applications/{id}", h.DeleteApplication)
			r.Post("/resume/generate", h.GenerateResume)
			r.Post("/cover-letter/generate", h.GenerateCoverLetter)
			r.Get("/tasks/{id}", h.GetTaskStatus)
			r.Get("/profile", h.GetProfile)
			r.Put("/profile", h.UpdateProfile)
			r.Get("/settings", h.GetSettings)
			r.Put("/settings", h.UpdateSettings)
			r.Post("/jobs/ingest", h.IngestJobs)
			r.Post("/scrape/trigger", h.TriggerScrape)
			r.Post("/jobs/{id}/auto-apply", h.AutoApply)

			// Payment routes (protected)
			r.Post("/payments/checkout", h.CreateCheckoutSession)
			r.Get("/payments/subscription", h.GetSubscription)
			r.Get("/payments/portal", h.CreateBillingPortal)
			r.Post("/payments/cancel", h.CancelSubscription)
		})
	})

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		lg.Info().Str("port", port).Msg("Server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	<-done
	lg.Info().Msg("Server shutting down gracefully...")

	// Shutdown the Asynq worker server
	asynqServer.Shutdown()
	lg.Info().Msg("Asynq worker server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		lg.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	lg.Info().Msg("Server stopped")
}
