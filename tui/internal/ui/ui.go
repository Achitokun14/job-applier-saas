package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"job-applier-tui/internal/api"
)

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

type model struct {
	state       state
	stateStack  []state
	api         *api.Client
	inputs      []textinput.Model
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

func New() model {
	m := model{
		state: stateMenu,
		api:   api.NewClient(""),
		inputs: make([]textinput.Model, 4),
	}

	for i := range m.inputs {
		m.inputs[i] = textinput.New()
		m.inputs[i].CharLimit = 150
		m.inputs[i].Width = 50
	}

	m.inputs[0].Placeholder = "Email"
	m.inputs[1].Placeholder = "Password"
	m.inputs[1].EchoMode = textinput.EchoPassword
	m.inputs[2].Placeholder = "Name"
	m.inputs[3].Placeholder = "Search query"

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if len(m.stateStack) > 0 {
				m.state = m.stateStack[len(m.stateStack)-1]
				m.stateStack = m.stateStack[:len(m.stateStack)-1]
			} else if m.state != stateMenu {
				m.state = stateMenu
			}
			return m, nil
		case "tab":
			m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
			return m, m.updateFocus()
		case "enter":
			return m.handleEnter()
		}
	}

	m.inputs[m.focusIndex], _ = m.inputs[m.focusIndex].Update(msg)
	return m, nil
}

func (m model) updateFocus() tea.Cmd {
	for i := range m.inputs {
		if i == m.focusIndex {
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
			m.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		} else {
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = lipgloss.NewStyle()
			m.inputs[i].TextStyle = lipgloss.NewStyle()
		}
	}
	return textinput.Blink
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	m.err = ""
	m.success = ""

	switch m.state {
	case stateMenu:
		return m.handleMenuSelection()
	case stateLogin:
		return m.handleLogin()
	case stateRegister:
		return m.handleRegister()
	case stateSearchJob:
		return m.handleSearchJob()
	case stateApplyJob:
		return m.handleApplyJob()
	}

	return m, nil
}

func (m model) handleMenuSelection() (tea.Model, tea.Cmd) {
	menuItems := []string{"Login", "Register", "Search Jobs", "View Applications", "Profile", "Settings"}
	if m.focusIndex < len(menuItems) {
		switch menuItems[m.focusIndex] {
		case "Login":
			m.stateStack = append(m.stateStack, stateMenu)
			m.state = stateLogin
			m.focusIndex = 0
			m.resetInputs("email", "password")
		case "Register":
			m.stateStack = append(m.stateStack, stateMenu)
			m.state = stateRegister
			m.focusIndex = 0
			m.resetInputs("email", "password", "name")
		case "Search Jobs":
			m.stateStack = append(m.stateStack, stateMenu)
			m.state = stateSearchJob
			m.focusIndex = 0
			m.resetInputs("query")
		case "View Applications":
			m.stateStack = append(m.stateStack, stateMenu)
			m.state = stateApplications
			m.loadApplications()
		case "Profile":
			m.stateStack = append(m.stateStack, stateMenu)
			m.state = stateProfile
			m.loadProfile()
		case "Settings":
			m.stateStack = append(m.stateStack, stateMenu)
			m.state = stateSettings
			m.loadSettings()
		}
	}
	return m, nil
}

func (m model) handleLogin() (tea.Model, tea.Cmd) {
	email := m.inputs[0].Value()
	password := m.inputs[1].Value()

	if email == "" || password == "" {
		m.err = "Email and password are required"
		return m, nil
	}

	_, err := m.api.Login(email, password)
	if err != nil {
		m.err = fmt.Sprintf("Login failed: %v", err)
		return m, nil
	}

	m.success = "Login successful!"
	m.state = stateMenu
	m.stateStack = nil
	return m, nil
}

func (m model) handleRegister() (tea.Model, tea.Cmd) {
	email := m.inputs[0].Value()
	password := m.inputs[1].Value()
	name := m.inputs[2].Value()

	if email == "" || password == "" {
		m.err = "Email and password are required"
		return m, nil
	}

	_, err := m.api.Register(email, password, name)
	if err != nil {
		m.err = fmt.Sprintf("Registration failed: %v", err)
		return m, nil
	}

	m.success = "Registration successful!"
	m.state = stateMenu
	m.stateStack = nil
	return m, nil
}

func (m model) handleSearchJob() (tea.Model, tea.Cmd) {
	query := m.inputs[3].Value()
	if query == "" {
		query = "software engineer"
	}

	jobs, err := m.api.GetJobs(query)
	if err != nil {
		m.err = fmt.Sprintf("Failed to search jobs: %v", err)
		return m, nil
	}

	m.jobs = jobs
	m.state = stateJobs
	m.stateStack = append(m.stateStack, stateSearchJob)
	return m, nil
}

func (m model) handleApplyJob() (tea.Model, tea.Cmd) {
	if m.selectedJob < len(m.jobs) {
		job := m.jobs[m.selectedJob]
		jobID := fmt.Sprintf("%v", job["id"])

		err := m.api.ApplyJob(jobID)
		if err != nil {
			m.err = fmt.Sprintf("Failed to apply: %v", err)
			return m, nil
		}

		m.success = "Application submitted!"
		m.state = stateJobs
	}
	return m, nil
}

func (m model) loadApplications() {
	apps, err := m.api.GetApplications()
	if err != nil {
		m.err = fmt.Sprintf("Failed to load applications: %v", err)
		return
	}
	m.apps = apps
}

func (m model) loadProfile() {
	profile, err := m.api.GetProfile()
	if err != nil {
		m.err = fmt.Sprintf("Failed to load profile: %v", err)
		return
	}
	m.profile = profile
}

func (m model) loadSettings() {
	settings, err := m.api.GetSettings()
	if err != nil {
		m.err = fmt.Sprintf("Failed to load settings: %v", err)
		return
	}
	m.settings = settings
}

func (m model) resetInputs(fields ...string) {
	for i := range m.inputs {
		m.inputs[i].SetValue("")
		m.inputs[i].Blur()
	}

	for i, field := range fields {
		if i < len(m.inputs) {
			m.inputs[i].Placeholder = field
			if i == 0 {
				m.inputs[i].Focus()
			}
		}
	}
}

func (m model) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00"))

	b.WriteString(titleStyle.Render("Job Applier TUI"))
	b.WriteString("\n\n")

	if m.err != "" {
		b.WriteString(errorStyle.Render("Error: "+m.err))
		b.WriteString("\n\n")
	}
	if m.success != "" {
		b.WriteString(successStyle.Render(m.success))
		b.WriteString("\n\n")
	}

	switch m.state {
	case stateMenu:
		b.WriteString(m.renderMenu())
	case stateLogin:
		b.WriteString(m.renderLogin())
	case stateRegister:
		b.WriteString(m.renderRegister())
	case stateSearchJob:
		b.WriteString(m.renderSearchJob())
	case stateJobs:
		b.WriteString(m.renderJobs())
	case stateApplications:
		b.WriteString(m.renderApplications())
	case stateProfile:
		b.WriteString(m.renderProfile())
	case stateSettings:
		b.WriteString(m.renderSettings())
	case stateApplyJob:
		b.WriteString(m.renderApplyJob())
	}

	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("Press ESC to go back | Ctrl+C to quit"))

	return b.String()
}

func (m model) renderMenu() string {
	menuStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("#7D56F4"))

	items := []string{
		"1. Login",
		"2. Register",
		"3. Search Jobs",
		"4. View Applications",
		"5. Profile",
		"6. Settings",
	}

	for i := range items {
		if i == m.focusIndex {
			items[i] = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true).
				Render("> " + items[i])
		} else {
			items[i] = "  " + items[i]
		}
	}

	return menuStyle.Render(strings.Join(items, "\n"))
}

func (m model) renderLogin() string {
	formStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("#7D56F4"))

	var b strings.Builder
	b.WriteString("Login\n\n")
	b.WriteString(m.inputs[0].View() + "\n")
	b.WriteString(m.inputs[1].View() + "\n\n")
	b.WriteString("Press Enter to login")

	return formStyle.Render(b.String())
}

func (m model) renderRegister() string {
	formStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("#7D56F4"))

	var b strings.Builder
	b.WriteString("Register\n\n")
	b.WriteString(m.inputs[0].View() + "\n")
	b.WriteString(m.inputs[1].View() + "\n")
	b.WriteString(m.inputs[2].View() + "\n\n")
	b.WriteString("Press Enter to register")

	return formStyle.Render(b.String())
}

func (m model) renderSearchJob() string {
	formStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("#7D56F4"))

	var b strings.Builder
	b.WriteString("Search Jobs\n\n")
	b.WriteString(m.inputs[3].View() + "\n\n")
	b.WriteString("Press Enter to search")

	return formStyle.Render(b.String())
}

func (m model) renderJobs() string {
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("#7D56F4"))

	var b strings.Builder
	b.WriteString("Jobs Found\n\n")

	if len(m.jobs) == 0 {
		b.WriteString("No jobs found. Try searching with different keywords.")
	} else {
		for i, job := range m.jobs {
			title := fmt.Sprintf("%v", job["title"])
			company := fmt.Sprintf("%v", job["company"])
			location := fmt.Sprintf("%v", job["location"])

			if i == m.selectedJob {
				b.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color("#7D56F4")).
					Bold(true).
					Render(fmt.Sprintf("> %s\n   %s - %s\n", title, company, location)))
			} else {
				b.WriteString(fmt.Sprintf("  %s\n   %s - %s\n", title, company, location))
			}
			b.WriteString("\n")
		}
		b.WriteString("Press Enter to apply to selected job | j/k to navigate")
	}

	return listStyle.Render(b.String())
}

func (m model) renderApplications() string {
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("#7D56F4"))

	var b strings.Builder
	b.WriteString("Applications\n\n")

	if len(m.apps) == 0 {
		b.WriteString("No applications yet.")
	} else {
		for _, app := range m.apps {
			job := app["job"].(map[string]interface{})
			title := fmt.Sprintf("%v", job["title"])
			company := fmt.Sprintf("%v", job["company"])
			status := fmt.Sprintf("%v", app["status"])
			b.WriteString(fmt.Sprintf("- %s at %s [%s]\n", title, company, status))
		}
	}

	return listStyle.Render(b.String())
}

func (m model) renderProfile() string {
	profileStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("#7D56F4"))

	var b strings.Builder
	b.WriteString("Profile\n\n")

	if m.profile != nil {
		b.WriteString(fmt.Sprintf("Name: %v\n", m.profile["name"]))
		b.WriteString(fmt.Sprintf("Email: %v\n", m.profile["email"]))
	} else {
		b.WriteString("No profile data loaded.")
	}

	return profileStyle.Render(b.String())
}

func (m model) renderSettings() string {
	settingsStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("#7D56F4"))

	var b strings.Builder
	b.WriteString("Settings\n\n")

	if m.settings != nil {
		b.WriteString(fmt.Sprintf("LLM Provider: %v\n", m.settings["llm_provider"]))
		b.WriteString(fmt.Sprintf("LLM Model: %v\n", m.settings["llm_model"]))
		b.WriteString(fmt.Sprintf("Remote Jobs: %v\n", m.settings["job_search_remote"]))
		b.WriteString(fmt.Sprintf("Experience Level: %v\n", m.settings["experience_level"]))
	} else {
		b.WriteString("No settings loaded.")
	}

	return settingsStyle.Render(b.String())
}

func (m model) renderApplyJob() string {
	applyStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("#7D56F4"))

	var b strings.Builder
	b.WriteString("Apply to Job\n\n")

	if m.selectedJob < len(m.jobs) {
		job := m.jobs[m.selectedJob]
		b.WriteString(fmt.Sprintf("Title: %v\n", job["title"]))
		b.WriteString(fmt.Sprintf("Company: %v\n", job["company"]))
		b.WriteString("\nPress Enter to confirm application")
	}

	return applyStyle.Render(b.String())
}
