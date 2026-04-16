# TUI Technical Documentation

## Overview

The TUI (Terminal User Interface) is a Go application using Bubble Tea for a local CLI interface to test the job application workflow.

## Directory Structure

```
tui/
├── cmd/
│   └── main.go                     # Entry point
├── internal/
│   ├── api/
│   │   └── client.go               # API client
│   ├── models/
│   │   └── models.go               # Data models
│   └── ui/
│       └── ui.go                   # TUI components
└── go.mod                          # Go module
```

## Entry Point (`cmd/main.go`)

```go
func main() {
    p := tea.NewProgram(ui.New(), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}
```

## UI Architecture (`internal/ui/ui.go`)

### State Machine

```go
type state int

const (
    stateMenu state = iota
    stateLogin
    stateRegister
    stateJobs
    stateApplications
    stateProfile
    stateSettings
    stateSearchJob
    stateApplyJob
)
```

### Model Structure

```go
type model struct {
    state       state
    stateStack  []state           // For back navigation
    api         *api.Client
    inputs      []textinput.Model // Form inputs
    focusIndex  int
    jobs        []map[string]interface{}
    apps        []map[string]interface{}
    profile     map[string]interface{}
    settings    map[string]interface{}
    selectedJob int
    err         string
    success     string
    width       int
    height      int
}
```

### Bubble Tea Implementation

Implements the Elm architecture:

```go
func (m model) Init() tea.Cmd {
    return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    case tea.KeyMsg:
        // Handle keyboard input
    }
    return m, nil
}

func (m model) View() string {
    // Render current state
}
```

### Navigation

- **ESC**: Go back to previous state (using stateStack)
- **Enter**: Submit current form/selection
- **Tab**: Cycle through form inputs
- **Ctrl+C**: Quit application
- **1-6**: Quick menu selection

## API Client (`internal/api/client.go`)

### Client Structure

```go
type Client struct {
    BaseURL    string
    Token      string
    HTTPClient *http.Client
}
```

### Methods

```go
func NewClient(baseURL string) *Client
func (c *Client) Register(email, password, name string) (string, error)
func (c *Client) Login(email, password string) (string, error)
func (c *Client) GetJobs(query string) ([]map[string]interface{}, error)
func (c *Client) ApplyJob(jobID string) error
func (c *Client) GetApplications() ([]map[string]interface{}, error)
func (c *Client) GetProfile() (map[string]interface{}, error)
func (c *Client) UpdateProfile(data map[string]interface{}) error
func (c *Client) GetSettings() (map[string]interface{}, error)
func (c *Client) UpdateSettings(data map[string]interface{}) error
```

### Request Pattern

```go
func (c *Client) request(method, path string, body interface{}, result interface{}) error {
    // 1. Marshal body to JSON
    // 2. Create HTTP request
    // 3. Add Authorization header if token exists
    // 4. Execute request
    // 5. Decode response
    // 6. Return result
}
```

## Views

### Menu View

```
┌─────────────────────────────────────┐
│        Job Applier TUI              │
│                                     │
│  > 1. Login                         │
│    2. Register                      │
│    3. Search Jobs                   │
│    4. View Applications             │
│    5. Profile                       │
│    6. Settings                      │
│                                     │
│  Press ESC to go back | Ctrl+C quit │
└─────────────────────────────────────┘
```

### Login View

```
┌─────────────────────────────────────┐
│        Login                        │
│                                     │
│  Email                              │
│  ┌─────────────────────────────────┐│
│  │ user@example.com                ││
│  └─────────────────────────────────┘│
│  Password                           │
│  ┌─────────────────────────────────┐│
│  │ ********                        ││
│  └─────────────────────────────────┘│
│                                     │
│  Press Enter to login               │
└─────────────────────────────────────┘
```

### Jobs View

```
┌─────────────────────────────────────┐
│        Jobs Found                   │
│                                     │
│  > Software Engineer                │
│    Tech Corp - San Francisco        │
│                                     │
│    Backend Developer                │
│    Startup Inc - Remote             │
│                                     │
│  Press Enter to apply | j/k nav     │
└─────────────────────────────────────┘
```

## Styling

Uses Lipgloss for terminal styling:

```go
titleStyle := lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    Padding(0, 1)

menuStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    Padding(1, 2).
    BorderForeground(lipgloss.Color("#7D56F4"))
```

## Error Handling

- API errors displayed in red
- Success messages displayed in green
- Errors clear on next action
- Network timeouts handled gracefully

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `BACKEND_URL` | `http://localhost:8080` | Backend API URL |

## Dependencies

```go
require (
    github.com/charmbracelet/bubbles v0.18.0  // UI components
    github.com/charmbracelet/bubbletea v0.26.6 // TUI framework
    github.com/charmbracelet/lipgloss v0.9.1   // Styling
)
```

## Build

```bash
# Development
go run cmd/main.go

# Production build
go build -o bin/tui cmd/main.go

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o bin/tui-linux cmd/main.go
GOOS=darwin GOARCH=amd64 go build -o bin/tui-mac cmd/main.go
GOOS=windows GOARCH=amd64 go build -o bin/tui.exe cmd/main.go
```

## Usage

1. Start the backend server
2. Run the TUI: `go run cmd/main.go`
3. Login or register
4. Search for jobs
5. Apply to jobs
6. View applications
7. Manage profile and settings

## Future Enhancements

- [ ] Keyboard shortcuts help screen
- [ ] Job details preview
- [ ] Application status updates
- [ ] Export applications to CSV
- [ ] Theme customization
- [ ] Multi-language support
