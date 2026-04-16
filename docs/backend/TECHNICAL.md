# Backend Technical Documentation

## Overview

The backend is a Go application using the Chi router for HTTP handling and GORM for database operations.

## Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration loading
│   ├── database/
│   │   └── database.go             # Database connection & migration
│   ├── handlers/
│   │   └── handlers.go             # HTTP request handlers
│   ├── models/
│   │   └── models.go               # Data models (GORM)
│   └── middleware/                  # HTTP middleware (future)
├── go.mod                           # Go module definition
└── go.sum                           # Dependency checksums
```

## Entry Point (`cmd/server/main.go`)

The main function:
1. Loads configuration from environment variables
2. Connects to the database (SQLite or PostgreSQL)
3. Runs auto-migrations
4. Sets up the Chi router with middleware
5. Registers all API routes
6. Starts the HTTP server

```go
func main() {
    cfg := config.Load()
    db, _ := database.Connect(cfg.DatabaseURL)
    database.AutoMigrate(db)
    
    r := chi.NewRouter()
    // ... middleware and routes
    http.ListenAndServe(":8080", r)
}
```

## Configuration (`internal/config`)

### Config Struct

```go
type Config struct {
    DatabaseURL      string   // SQLite or PostgreSQL connection string
    JWTSecret        string   // JWT signing secret
    JWTExpiry        string   // Token expiry duration
    PythonServiceURL string   // Python service endpoint
}
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `sqlite:./data/jobapplier.db` | Database connection |
| `JWT_SECRET` | `default-secret-...` | JWT secret key |
| `JWT_EXPIRY` | `24h` | Token validity |
| `PYTHON_SERVICE_URL` | `http://localhost:8001` | Python service |

## Database (`internal/database`)

### Connection

Supports two database backends:
- **SQLite**: `sqlite:./data/jobapplier.db`
- **PostgreSQL**: `postgres://user:pass@host:5432/dbname`

```go
func Connect(databaseURL string) (*gorm.DB, error) {
    if strings.HasPrefix(databaseURL, "sqlite:") {
        return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
    }
    return gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
}
```

### Auto-Migration

GORM auto-migrates these models:
- `User`
- `Resume`
- `Job`
- `Application`
- `Settings`

## Models (`internal/models`)

### User

```go
type User struct {
    ID        uint           `gorm:"primarykey"`
    Email     string         `gorm:"uniqueIndex;not null"`
    Password  string         `gorm:"not null"`  // bcrypt hash
    Name      string
    Resume    Resume         `gorm:"foreignKey:UserID"`
}
```

### Resume

```go
type Resume struct {
    ID            uint   `gorm:"primarykey"`
    UserID        uint   `gorm:"uniqueIndex"`
    PersonalInfo  string `gorm:"type:text"`  // JSON
    Education     string `gorm:"type:text"`  // JSON array
    Experience    string `gorm:"type:text"`  // JSON array
    Skills        string `gorm:"type:text"`  // JSON array
    Projects      string `gorm:"type:text"`  // JSON array
    PDFPath       string
}
```

### Job

```go
type Job struct {
    ID          uint   `gorm:"primarykey"`
    ExternalID  string `gorm:"index"`
    Title       string `gorm:"index"`
    Company     string `gorm:"index"`
    Location    string
    Description string `gorm:"type:text"`
    URL         string
    Source      string `gorm:"index"`  // linkedin, indeed, etc.
    Remote      bool
    Salary      string
}
```

### Application

```go
type Application struct {
    ID        uint      `gorm:"primarykey"`
    UserID    uint      `gorm:"index"`
    JobID     uint      `gorm:"index"`
    Job       Job       `gorm:"foreignKey:JobID"`
    Status    string    // applied, interview, offer, rejected
    ResumePDF string
    CoverPDF  string
    AppliedAt time.Time
}
```

### Settings

```go
type Settings struct {
    ID              uint   `gorm:"primarykey"`
    UserID          uint   `gorm:"uniqueIndex"`
    LLMProvider     string // openai, anthropic, google, ollama
    LLMModel        string // gpt-4o-mini, claude-3, etc.
    LLMAPIKey       string
    JobSearchRemote bool
    ExperienceLevel string // entry, mid_senior, director, etc.
    JobTypes        string // full_time, contract, etc.
    Positions       string // comma-separated
    Locations       string // comma-separated
    Distance        int    // miles
}
```

## Handlers (`internal/handlers`)

### Handler Struct

```go
type Handlers struct {
    db  *gorm.DB
    cfg *config.Config
}
```

### Authentication Handlers

#### Register (`POST /api/v1/auth/register`)
- Validates input
- Checks for existing email
- Hashes password with bcrypt
- Creates user
- Generates JWT token
- Returns token + user data

#### Login (`POST /api/v1/auth/login`)
- Finds user by email
- Verifies password with bcrypt
- Generates JWT token
- Returns token + user data

### JWT Token Structure

```json
{
  "user_id": 1,
  "exp": 1711612800,
  "iat": 1711526400
}
```

### Auth Middleware

Validates JWT token from Authorization header:
1. Extracts Bearer token
2. Parses and validates JWT
3. Extracts user_id from claims
4. Adds user_id to request context

### Job Handlers

#### Search Jobs (`GET /api/v1/jobs`)
- Queries database with optional filters
- Supports text search (title, company, description)
- Supports source filter
- Pagination (20 per page)

#### Apply Job (`POST /api/v1/jobs/{id}/apply`)
- Creates Application record
- Sets status to "applied"
- Records timestamp

### Application Handlers

#### List Applications (`GET /api/v1/applications`)
- Returns all user's applications
- Preloads Job data
- Ordered by applied_at DESC

#### Get Application (`GET /api/v1/applications/{id}`)
- Returns single application with job details
- Verifies ownership

#### Delete Application (`DELETE /api/v1/applications/{id}`)
- Deletes application record
- Verifies ownership

### Profile Handlers

#### Get Profile (`GET /api/v1/profile`)
- Returns user data with resume

#### Update Profile (`PUT /api/v1/profile`)
- Updates user name
- Creates or updates Resume record
- Stores JSON-structured data

### Settings Handlers

#### Get Settings (`GET /api/v1/settings`)
- Returns user's settings
- Creates default if not exists

#### Update Settings (`PUT /api/v1/settings`)
- Updates all settings fields
- LLMAPIKey only updated if provided

## Middleware

### CORS

```go
cors.Handler(cors.Options{
    AllowedOrigins:   []string{"*"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
    AllowCredentials: true,
    MaxAge:           300,
})
```

### Logger

Chi's built-in request logger.

### Recoverer

Catches panics and returns 500.

## Error Handling

All errors return JSON:
```json
{
  "error": "Error message here"
}
```

## Dependencies

```go
require (
    github.com/go-chi/chi/v5 v5.1.0      // HTTP router
    github.com/go-chi/cors v1.2.1         // CORS middleware
    github.com/golang-jwt/jwt/v5 v5.2.1   // JWT handling
    github.com/glebarez/sqlite v1.11.0     // SQLite driver
    gorm.io/driver/postgres v1.5.9         // PostgreSQL driver
    gorm.io/gorm v1.25.11                  // ORM
    golang.org/x/crypto v0.27.0            // bcrypt
)
```

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test ./internal/handlers -run TestLogin
```

## Build

```bash
# Development
go run cmd/server/main.go

# Production build
go build -o bin/server cmd/server/main.go

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o bin/server-linux cmd/server/main.go
```

## Performance Considerations

- Database connection pooling via GORM
- Stateless handlers (no in-memory sessions)
- JWT validation on each request (fast with HS256)
- SQLite suitable for development (<1000 users)
- PostgreSQL recommended for production
