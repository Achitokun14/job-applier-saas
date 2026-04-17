package handlers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"job-applier-backend/internal/auth"
	"job-applier-backend/internal/config"
	"job-applier-backend/internal/middleware"
	"job-applier-backend/internal/models"
	"job-applier-backend/internal/repository"
	"job-applier-backend/internal/services"
	"job-applier-backend/internal/tasks"
)

type Handlers struct {
	db           *gorm.DB
	cfg          *config.Config
	taskClient   *asynq.Client
	blacklist    *auth.TokenBlacklist
	jobRepo      *repository.JobRepository
	appRepo      *repository.ApplicationRepository
	userRepo     *repository.UserRepository
	settingsRepo *repository.SettingsRepository
}

func New(
	db *gorm.DB,
	cfg *config.Config,
	taskClient *asynq.Client,
	blacklist *auth.TokenBlacklist,
	jobRepo *repository.JobRepository,
	appRepo *repository.ApplicationRepository,
	userRepo *repository.UserRepository,
	settingsRepo *repository.SettingsRepository,
) *Handlers {
	return &Handlers{
		db:           db,
		cfg:          cfg,
		taskClient:   taskClient,
		blacklist:    blacklist,
		jobRepo:      jobRepo,
		appRepo:      appRepo,
		userRepo:     userRepo,
		settingsRepo: settingsRepo,
	}
}

type contextKey string

const userIDKey contextKey = "userID"
const userRoleKey contextKey = "userRole"

func (h *Handlers) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handlers) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.Email == "" || input.Password == "" {
		h.respondError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Validate email format
	if !strings.Contains(input.Email, "@") || !strings.Contains(input.Email, ".") {
		h.respondError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Validate password strength
	if len(input.Password) < 8 {
		h.respondError(w, http.StatusBadRequest, "Password must be at least 8 characters")
		return
	}

	var existing models.User
	if err := h.db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		h.respondError(w, http.StatusConflict, "Email already registered")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ERROR: Failed to hash password for %s: %v", input.Email, err)
		h.respondError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := models.User{
		Email:         input.Email,
		Password:      string(hashedPassword),
		Name:          input.Name,
		EmailVerified: false,
	}

	if err := h.db.Create(&user).Error; err != nil {
		log.Printf("ERROR: Failed to create user %s: %v", input.Email, err)
		h.respondError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	accessToken, err := h.generateAccessToken(user.ID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := h.generateRefreshToken(user.ID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]interface{}{
		"token":         accessToken,
		"refresh_token": refreshToken,
		"expires_in":    900,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		h.respondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		h.respondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	accessToken, err := h.generateAccessToken(user.ID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := h.generateRefreshToken(user.ID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"token":         accessToken,
		"refresh_token": refreshToken,
		"expires_in":    900,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func (h *Handlers) generateAccessToken(userID uint) (string, error) {
	jti := uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
		"jti":     jti,
	})

	return token.SignedString([]byte(h.cfg.JWTSecret))
}

func (h *Handlers) generateRefreshToken(userID uint) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	tokenStr := hex.EncodeToString(tokenBytes)

	refreshToken := models.RefreshToken{
		UserID:    userID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := h.db.Create(&refreshToken).Error; err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (h *Handlers) generateRandomToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(tokenBytes), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)

	user, err := h.userRepo.FindByID(r.Context(), userID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "User not found")
		return
	}

	// Flatten response: user fields + resume fields at same level
	response := map[string]interface{}{
		"id":             user.ID,
		"email":          user.Email,
		"name":           user.Name,
		"role":           user.Role,
		"email_verified": user.EmailVerified,
		"personal_info":  user.Resume.PersonalInfo,
		"education":      user.Resume.Education,
		"experience":     user.Resume.Experience,
		"skills":         user.Resume.Skills,
		"projects":       user.Resume.Projects,
		"achievements":   user.Resume.Achievements,
		"certifications": user.Resume.Certifications,
		"languages":      user.Resume.Languages,
		"pdf_path":       user.Resume.PDFPath,
	}
	h.respondJSON(w, http.StatusOK, response)
}

func (h *Handlers) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)

	var input struct {
		Name        string `json:"name"`
		PersonalInfo string `json:"personal_info"`
		Education   string `json:"education"`
		Experience  string `json:"experience"`
		Skills      string `json:"skills"`
		Projects    string `json:"projects"`
		Achievements string `json:"achievements"`
		Certifications string `json:"certifications"`
		Languages   string `json:"languages"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		h.respondError(w, http.StatusNotFound, "User not found")
		return
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	h.db.Save(&user)

	var resume models.Resume
	result := h.db.Where("user_id = ?", userID).First(&resume)
	if result.Error != nil {
		resume = models.Resume{UserID: userID}
	}

	if input.PersonalInfo != "" {
		resume.PersonalInfo = input.PersonalInfo
	}
	if input.Education != "" {
		resume.Education = input.Education
	}
	if input.Experience != "" {
		resume.Experience = input.Experience
	}
	if input.Skills != "" {
		resume.Skills = input.Skills
	}
	if input.Projects != "" {
		resume.Projects = input.Projects
	}
	if input.Achievements != "" {
		resume.Achievements = input.Achievements
	}
	if input.Certifications != "" {
		resume.Certifications = input.Certifications
	}
	if input.Languages != "" {
		resume.Languages = input.Languages
	}

	if resume.ID == 0 {
		h.db.Create(&resume)
	} else {
		h.db.Save(&resume)
	}

	// Invalidate user cache since profile was updated
	if h.userRepo != nil {
		_ = h.userRepo.Update(r.Context(), &user)
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Profile updated"})
}

func (h *Handlers) SearchJobs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	source := r.URL.Query().Get("source")
	pageStr := r.URL.Query().Get("page")
	location := r.URL.Query().Get("location")
	remote := r.URL.Query().Get("remote")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	jobs, total, err := h.jobRepo.Search(r.Context(), query, source, page)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to search jobs")
		return
	}

	// Apply additional filters that the repository doesn't handle
	if location != "" || remote != "" {
		filtered := make([]models.Job, 0, len(jobs))
		for _, job := range jobs {
			if location != "" && !strings.Contains(strings.ToLower(job.Location), strings.ToLower(location)) {
				continue
			}
			if remote == "true" && !job.Remote {
				continue
			}
			if remote == "false" && job.Remote {
				continue
			}
			filtered = append(filtered, job)
		}
		jobs = filtered
		total = int64(len(filtered))
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"jobs":  jobs,
		"total": total,
		"page":  page,
	})
}

func (h *Handlers) ApplyJob(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)
	jobID := chi.URLParam(r, "id")

	// Check usage limit for applications
	allowed, err := services.CheckLimit(h.db, userID, "application")
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to check usage limit")
		return
	}
	if !allowed {
		tier, limit, used := services.GetUserTierAndUsage(h.db, userID, "application")
		h.respondJSON(w, http.StatusPaymentRequired, map[string]interface{}{
			"error": "Usage limit exceeded. Upgrade to Pro.",
			"tier":  tier,
			"limit": limit,
			"used":  used,
		})
		return
	}

	var job models.Job
	if err := h.db.First(&job, jobID).Error; err != nil {
		h.respondError(w, http.StatusNotFound, "Job not found")
		return
	}

	var existing models.Application
	if err := h.db.Where("user_id = ? AND job_id = ?", userID, job.ID).First(&existing).Error; err == nil {
		h.respondError(w, http.StatusConflict, "Already applied to this job")
		return
	}

	application := models.Application{
		UserID:    userID,
		JobID:     job.ID,
		Status:    "applied",
		AppliedAt: time.Now(),
	}

	if err := h.db.Create(&application).Error; err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create application")
		return
	}

	// Increment usage counter after successful application
	_ = services.IncrementUsage(h.db, userID, "application")

	h.respondJSON(w, http.StatusCreated, application)
}

func (h *Handlers) ListApplications(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)
	status := r.URL.Query().Get("status")
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	perPage := 20
	if perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	applications, total, err := h.appRepo.ListByUser(r.Context(), userID, page, perPage, status)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to list applications")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"applications": applications,
		"total":        total,
		"page":         page,
		"per_page":     perPage,
	})
}

func (h *Handlers) GetApplication(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)
	appID := chi.URLParam(r, "id")

	var application models.Application
	if err := h.db.Preload("Job").Where("user_id = ? AND id = ?", userID, appID).First(&application).Error; err != nil {
		h.respondError(w, http.StatusNotFound, "Application not found")
		return
	}

	h.respondJSON(w, http.StatusOK, application)
}

func (h *Handlers) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)
	appID := chi.URLParam(r, "id")

	if err := h.db.Where("user_id = ? AND id = ?", userID, appID).Delete(&models.Application{}).Error; err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to delete application")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Application deleted"})
}

func (h *Handlers) GenerateResume(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)

	// Check usage limit for resume generation
	allowed, err := services.CheckLimit(h.db, userID, "resume_gen")
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to check usage limit")
		return
	}
	if !allowed {
		tier, limit, used := services.GetUserTierAndUsage(h.db, userID, "resume_gen")
		h.respondJSON(w, http.StatusPaymentRequired, map[string]interface{}{
			"error": "Usage limit exceeded. Upgrade to Pro.",
			"tier":  tier,
			"limit": limit,
			"used":  used,
		})
		return
	}

	var input struct {
		Style          string `json:"style"`
		JobDescription string `json:"job_description"`
		ApplicationID  uint   `json:"application_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get user's LLM settings
	var settings models.Settings
	if err := h.db.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		h.respondError(w, http.StatusBadRequest, "User settings not found; please configure LLM settings first")
		return
	}

	// Get user's resume data
	var resume models.Resume
	if err := h.db.Where("user_id = ?", userID).First(&resume).Error; err != nil {
		h.respondError(w, http.StatusBadRequest, "Resume data not found; please fill out your profile first")
		return
	}

	payload := tasks.ResumePayload{
		UserID:         userID,
		ApplicationID:  input.ApplicationID,
		Style:          input.Style,
		JobDescription: input.JobDescription,
		LLMProvider:    settings.LLMProvider,
		LLMModel:       settings.LLMModel,
		LLMAPIKey:      settings.LLMAPIKey,
		PersonalInfo:   resume.PersonalInfo,
		Education:      resume.Education,
		Experience:     resume.Experience,
		Skills:         resume.Skills,
		Projects:       resume.Projects,
		Achievements:   resume.Achievements,
		Certifications: resume.Certifications,
		Languages:      resume.Languages,
	}

	task, err := tasks.NewResumeTask(payload)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create resume task")
		return
	}

	info, err := h.taskClient.Enqueue(task)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to enqueue resume task")
		return
	}

	// Increment usage counter after successful enqueue
	_ = services.IncrementUsage(h.db, userID, "resume_gen")

	h.respondJSON(w, http.StatusAccepted, map[string]string{
		"task_id": info.ID,
		"status":  "queued",
		"queue":   info.Queue,
	})
}

func (h *Handlers) GenerateCoverLetter(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)

	// Check usage limit for cover letter generation
	allowed, err := services.CheckLimit(h.db, userID, "cover_letter_gen")
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to check usage limit")
		return
	}
	if !allowed {
		tier, limit, used := services.GetUserTierAndUsage(h.db, userID, "cover_letter_gen")
		h.respondJSON(w, http.StatusPaymentRequired, map[string]interface{}{
			"error": "Usage limit exceeded. Upgrade to Pro.",
			"tier":  tier,
			"limit": limit,
			"used":  used,
		})
		return
	}

	var input struct {
		JobDescription string `json:"job_description"`
		CompanyName    string `json:"company_name"`
		JobTitle       string `json:"job_title"`
		ApplicationID  uint   `json:"application_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get user's LLM settings
	var settings models.Settings
	if err := h.db.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		h.respondError(w, http.StatusBadRequest, "User settings not found; please configure LLM settings first")
		return
	}

	// Get user's resume data to build resume text for the cover letter
	var resume models.Resume
	if err := h.db.Where("user_id = ?", userID).First(&resume).Error; err != nil {
		h.respondError(w, http.StatusBadRequest, "Resume data not found; please fill out your profile first")
		return
	}

	// Build a combined resume text from the stored fields
	resumeText := fmt.Sprintf("Personal Info: %s\nEducation: %s\nExperience: %s\nSkills: %s\nProjects: %s\nAchievements: %s\nCertifications: %s\nLanguages: %s",
		resume.PersonalInfo, resume.Education, resume.Experience, resume.Skills,
		resume.Projects, resume.Achievements, resume.Certifications, resume.Languages)

	payload := tasks.CoverLetterPayload{
		UserID:         userID,
		ApplicationID:  input.ApplicationID,
		ResumeText:     resumeText,
		JobDescription: input.JobDescription,
		CompanyName:    input.CompanyName,
		JobTitle:       input.JobTitle,
		LLMProvider:    settings.LLMProvider,
		LLMModel:       settings.LLMModel,
		LLMAPIKey:      settings.LLMAPIKey,
	}

	task, err := tasks.NewCoverLetterTask(payload)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create cover letter task")
		return
	}

	info, err := h.taskClient.Enqueue(task)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to enqueue cover letter task")
		return
	}

	// Increment usage counter after successful enqueue
	_ = services.IncrementUsage(h.db, userID, "cover_letter_gen")

	h.respondJSON(w, http.StatusAccepted, map[string]string{
		"task_id": info.ID,
		"status":  "queued",
		"queue":   info.Queue,
	})
}

func (h *Handlers) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)

	settings, err := h.settingsRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		// Settings not found, create defaults
		settings = &models.Settings{
			UserID:          userID,
			LLMProvider:     "openai",
			LLMModel:        "gpt-4o-mini",
			JobSearchRemote: true,
			ExperienceLevel: "mid_senior",
			JobTypes:        "full_time",
			Distance:        50,
		}
		if err := h.settingsRepo.Upsert(r.Context(), settings); err != nil {
			h.respondError(w, http.StatusInternalServerError, "Failed to create default settings")
			return
		}
	}

	h.respondJSON(w, http.StatusOK, settings)
}

func (h *Handlers) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)

	var input models.Settings
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Fetch existing settings to preserve the API key if not provided
	existing, _ := h.settingsRepo.GetByUserID(r.Context(), userID)

	settings := models.Settings{
		UserID:           userID,
		LLMProvider:      input.LLMProvider,
		LLMModel:         input.LLMModel,
		LLMAPIKey:        input.LLMAPIKey,
		JobSearchRemote:  input.JobSearchRemote,
		JobSearchHybrid:  input.JobSearchHybrid,
		JobSearchOnsite:  input.JobSearchOnsite,
		ExperienceLevel:  input.ExperienceLevel,
		JobTypes:         input.JobTypes,
		Positions:        input.Positions,
		Locations:        input.Locations,
		Distance:         input.Distance,
		CompanyBlacklist: input.CompanyBlacklist,
		TitleBlacklist:   input.TitleBlacklist,
	}

	// Preserve existing API key if not provided in the update
	if settings.LLMAPIKey == "" && existing != nil {
		settings.LLMAPIKey = existing.LLMAPIKey
	}

	if err := h.settingsRepo.Upsert(r.Context(), &settings); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to update settings")
		return
	}

	h.respondJSON(w, http.StatusOK, settings)
}

func (h *Handlers) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		h.respondError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: h.cfg.RedisURL})
	defer inspector.Close()

	// Try to find the task across different states
	info, err := inspector.GetTaskInfo("critical", taskID)
	if err != nil {
		info, err = inspector.GetTaskInfo("default", taskID)
	}
	if err != nil {
		info, err = inspector.GetTaskInfo("low", taskID)
	}

	if err != nil {
		h.respondJSON(w, http.StatusOK, map[string]string{
			"task_id": taskID,
			"status":  "unknown",
			"message": "Task not found or already completed",
		})
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"task_id":      info.ID,
		"status":       info.State.String(),
		"queue":        info.Queue,
		"type":         info.Type,
		"max_retry":    info.MaxRetry,
		"retried":      info.Retried,
		"last_err":     info.LastErr,
		"next_process_at": info.NextProcessAt,
	})
}

func (h *Handlers) IngestJobs(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Jobs []struct {
			ExternalID  string `json:"external_id"`
			Source      string `json:"source"`
			Title       string `json:"title"`
			Company     string `json:"company"`
			Location    string `json:"location"`
			Description string `json:"description"`
			URL         string `json:"url"`
			Remote      bool   `json:"remote"`
			Salary      string `json:"salary"`
			PostedAt    string `json:"posted_at"`
		} `json:"jobs"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var jobs []models.Job
	for _, j := range input.Jobs {
		jobs = append(jobs, models.Job{
			ExternalID:  j.ExternalID,
			Source:      j.Source,
			Title:       j.Title,
			Company:     j.Company,
			Location:    j.Location,
			Description: j.Description,
			URL:         j.URL,
			Remote:      j.Remote,
			Salary:      j.Salary,
		})
	}

	inserted, skipped, err := h.jobRepo.BulkUpsert(r.Context(), jobs)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to ingest jobs")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]int{
		"inserted": inserted,
		"skipped":  skipped,
	})
}

func (h *Handlers) TriggerScrape(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)

	// Read user settings for search parameters
	var settings models.Settings
	if err := h.db.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		h.respondError(w, http.StatusBadRequest, "User settings not found; please configure job search settings first")
		return
	}

	positions := settings.Positions
	if positions == "" {
		h.respondError(w, http.StatusBadRequest, "No positions configured in settings")
		return
	}

	// Build scrape request from settings
	scrapeReq := map[string]interface{}{
		"search_term":    positions,
		"results_wanted": 50,
		"hours_old":      72,
		"distance":       settings.Distance,
		"sites":          []string{"indeed", "linkedin", "glassdoor", "google"},
	}

	if settings.Locations != "" {
		scrapeReq["location"] = settings.Locations
	}
	if settings.JobSearchRemote {
		scrapeReq["is_remote"] = true
	}

	// Allow triggering via Asynq queue
	var input struct {
		Async bool `json:"async"`
	}
	// Attempt to decode body but don't fail if empty
	json.NewDecoder(r.Body).Decode(&input)

	if input.Async {
		// Enqueue as Asynq task
		payload := tasks.ScrapePayload{
			SearchTerm: positions,
			Location:   settings.Locations,
			IsRemote:   settings.JobSearchRemote,
			Distance:   settings.Distance,
			Sites:      []string{"indeed", "linkedin", "glassdoor", "google"},
		}

		task, err := tasks.NewScrapeTask(payload)
		if err != nil {
			h.respondError(w, http.StatusInternalServerError, "Failed to create scrape task")
			return
		}

		info, err := h.taskClient.Enqueue(task)
		if err != nil {
			h.respondError(w, http.StatusInternalServerError, "Failed to enqueue scrape task")
			return
		}

		h.respondJSON(w, http.StatusAccepted, map[string]string{
			"task_id": info.ID,
			"status":  "queued",
			"queue":   info.Queue,
		})
		return
	}

	// Synchronous: call Python service directly and ingest results
	scrapeBody, err := json.Marshal(scrapeReq)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to marshal scrape request")
		return
	}

	pythonURL := h.cfg.PythonServiceURL + "/scrape-jobs"
	httpClient := &http.Client{Timeout: 300 * time.Second}
	resp, err := httpClient.Post(pythonURL, "application/json", strings.NewReader(string(scrapeBody)))
	if err != nil {
		h.respondError(w, http.StatusBadGateway, fmt.Sprintf("Failed to call Python scrape service: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.respondError(w, http.StatusBadGateway, fmt.Sprintf("Python scrape service returned status %d", resp.StatusCode))
		return
	}

	var scrapeResult struct {
		Jobs []struct {
			ExternalID  string `json:"external_id"`
			Source      string `json:"source"`
			Title       string `json:"title"`
			Company     string `json:"company"`
			Location    string `json:"location"`
			Description string `json:"description"`
			URL         string `json:"url"`
			Remote      bool   `json:"remote"`
			Salary      string `json:"salary"`
			PostedAt    string `json:"posted_at"`
		} `json:"jobs"`
		Total int `json:"total"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&scrapeResult); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to decode scrape response")
		return
	}

	// Ingest results into database via repository
	var jobs []models.Job
	for _, j := range scrapeResult.Jobs {
		jobs = append(jobs, models.Job{
			ExternalID:  j.ExternalID,
			Source:      j.Source,
			Title:       j.Title,
			Company:     j.Company,
			Location:    j.Location,
			Description: j.Description,
			URL:         j.URL,
			Remote:      j.Remote,
			Salary:      j.Salary,
		})
	}

	inserted, skipped, _ := h.jobRepo.BulkUpsert(r.Context(), jobs)

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"scraped":  scrapeResult.Total,
		"inserted": inserted,
		"skipped":  skipped,
	})
}

func (h *Handlers) AutoApply(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)
	jobID := chi.URLParam(r, "id")

	// 1. Get the job from the database
	var job models.Job
	if err := h.db.First(&job, jobID).Error; err != nil {
		h.respondError(w, http.StatusNotFound, "Job not found")
		return
	}

	// 2. Check for existing application
	var existing models.Application
	if err := h.db.Where("user_id = ? AND job_id = ?", userID, job.ID).First(&existing).Error; err == nil {
		h.respondError(w, http.StatusConflict, "Already applied to this job")
		return
	}

	// 3. Get user profile (name, email, phone)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		h.respondError(w, http.StatusNotFound, "User not found")
		return
	}

	// 4. Get user's latest resume PDF path
	var resume models.Resume
	if err := h.db.Where("user_id = ?", userID).First(&resume).Error; err != nil {
		h.respondError(w, http.StatusBadRequest, "Resume not found; please generate a resume first")
		return
	}

	resumePDFPath := resume.PDFPath
	if resumePDFPath == "" {
		h.respondError(w, http.StatusBadRequest, "No resume PDF available; please generate a resume first")
		return
	}

	// 5. Build the auto-apply request for the Python service
	applyURL := job.URL
	if applyURL == "" {
		h.respondError(w, http.StatusBadRequest, "Job has no application URL")
		return
	}

	autoApplyReq := map[string]interface{}{
		"job_url":          job.URL,
		"apply_url":        applyURL,
		"source":           job.Source,
		"resume_pdf_path":  resumePDFPath,
		"user_name":        user.Name,
		"user_email":       user.Email,
	}

	// 6. Call the Python service /api/auto-apply
	reqBody, err := json.Marshal(autoApplyReq)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to build auto-apply request")
		return
	}

	pythonURL := h.cfg.PythonServiceURL + "/api/auto-apply"
	httpClient := &http.Client{Timeout: 120 * time.Second}
	resp, err := httpClient.Post(pythonURL, "application/json", strings.NewReader(string(reqBody)))
	if err != nil {
		h.respondError(w, http.StatusBadGateway, fmt.Sprintf("Failed to call auto-apply service: %v", err))
		return
	}
	defer resp.Body.Close()

	var autoApplyResult struct {
		Success              bool   `json:"success"`
		Method               string `json:"method"`
		ConfirmationID       string `json:"confirmation_id"`
		ScreenshotPath       string `json:"screenshot_path"`
		Error                string `json:"error"`
		RequiresConfirmation bool   `json:"requires_confirmation"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&autoApplyResult); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to decode auto-apply response")
		return
	}

	// 7. Determine application status
	status := "preparing"
	if autoApplyResult.RequiresConfirmation {
		status = "ready_to_submit"
	} else if autoApplyResult.Success {
		status = "submitted"
	} else if autoApplyResult.Error != "" {
		status = "failed"
	}

	// 8. Create the Application record
	application := models.Application{
		UserID:          userID,
		JobID:           job.ID,
		Status:          status,
		ResumePDF:       resumePDFPath,
		AppliedAt:       time.Now(),
		AutoApplyMethod: autoApplyResult.Method,
		ScreenshotPath:  autoApplyResult.ScreenshotPath,
		ConfirmationID:  autoApplyResult.ConfirmationID,
	}

	if err := h.db.Create(&application).Error; err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create application record")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]interface{}{
		"application":          application,
		"auto_apply_success":   autoApplyResult.Success,
		"method":               autoApplyResult.Method,
		"requires_confirmation": autoApplyResult.RequiresConfirmation,
		"screenshot_path":      autoApplyResult.ScreenshotPath,
		"error":                autoApplyResult.Error,
	})
}

func (h *Handlers) BulkApply(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)

	var input struct {
		JobIDs []uint `json:"job_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if len(input.JobIDs) == 0 {
		h.respondError(w, http.StatusBadRequest, "No job IDs provided")
		return
	}
	if len(input.JobIDs) > 50 {
		h.respondError(w, http.StatusBadRequest, "Maximum 50 jobs per batch")
		return
	}

	applied := 0
	skipped := 0
	errors := []string{}

	for _, jobID := range input.JobIDs {
		// Check job exists
		var job models.Job
		if err := h.db.First(&job, jobID).Error; err != nil {
			errors = append(errors, fmt.Sprintf("Job %d not found", jobID))
			continue
		}

		// Check not already applied
		var existing models.Application
		if err := h.db.Where("user_id = ? AND job_id = ?", userID, jobID).First(&existing).Error; err == nil {
			skipped++
			continue
		}

		// Create application
		app := models.Application{
			UserID:    userID,
			JobID:     jobID,
			Status:    "applied",
			AppliedAt: time.Now(),
		}
		if err := h.db.Create(&app).Error; err != nil {
			errors = append(errors, fmt.Sprintf("Failed to apply to job %d", jobID))
			continue
		}
		applied++
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"applied": applied,
		"skipped": skipped,
		"errors":  errors,
	})
}

func (h *Handlers) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.respondError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			h.respondError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(h.cfg.JWTSecret), nil
		})

		if err != nil {
			h.respondError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check if token has been blacklisted
			if jti, ok := claims["jti"].(string); ok && jti != "" {
				if h.blacklist.IsBlacklisted(jti) {
					h.respondError(w, http.StatusUnauthorized, "Token has been revoked")
					return
				}
			}

			userID := uint(claims["user_id"].(float64))

			// Load user role from database
			var user models.User
			if err := h.db.Select("role").First(&user, userID).Error; err != nil {
				h.respondError(w, http.StatusUnauthorized, "User not found")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, middleware.ContextKey("userID"), userID)
			ctx = context.WithValue(ctx, userRoleKey, user.Role)
			ctx = context.WithValue(ctx, middleware.UserRoleKey, user.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			h.respondError(w, http.StatusUnauthorized, "Invalid token claims")
		}
	})
}
