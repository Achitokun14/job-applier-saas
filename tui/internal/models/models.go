package models

type Job struct {
	ID          uint   `json:"id"`
	ExternalID  string `json:"external_id"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Source      string `json:"source"`
	Remote      bool   `json:"remote"`
	Salary      string `json:"salary"`
}

type Application struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	JobID     uint   `json:"job_id"`
	Job       Job    `json:"job"`
	Status    string `json:"status"`
	AppliedAt string `json:"applied_at"`
}

type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Resume struct {
	UserID       uint   `json:"user_id"`
	PersonalInfo string `json:"personal_info"`
	Education    string `json:"education"`
	Experience   string `json:"experience"`
	Skills       string `json:"skills"`
	Projects     string `json:"projects"`
}

type Settings struct {
	LLMProvider     string `json:"llm_provider"`
	LLMModel        string `json:"llm_model"`
	JobSearchRemote bool   `json:"job_search_remote"`
	ExperienceLevel string `json:"experience_level"`
	JobTypes        string `json:"job_types"`
	Positions       string `json:"positions"`
	Locations       string `json:"locations"`
}
