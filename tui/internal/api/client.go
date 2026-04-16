package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = os.Getenv("BACKEND_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
	}

	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) request(method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}

func (c *Client) Register(email, password, name string) (string, error) {
	var resp struct {
		Token string `json:"token"`
		User  struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"user"`
	}

	err := c.request("POST", "/api/v1/auth/register", map[string]string{
		"email":    email,
		"password": password,
		"name":     name,
	}, &resp)

	if err != nil {
		return "", err
	}

	c.Token = resp.Token
	return resp.Token, nil
}

func (c *Client) Login(email, password string) (string, error) {
	var resp struct {
		Token string `json:"token"`
		User  struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"user"`
	}

	err := c.request("POST", "/api/v1/auth/login", map[string]string{
		"email":    email,
		"password": password,
	}, &resp)

	if err != nil {
		return "", err
	}

	c.Token = resp.Token
	return resp.Token, nil
}

func (c *Client) GetJobs(query string) ([]map[string]interface{}, error) {
	var jobs []map[string]interface{}
	err := c.request("GET", "/api/v1/jobs?q="+query, nil, &jobs)
	return jobs, err
}

func (c *Client) ApplyJob(jobID string) error {
	return c.request("POST", fmt.Sprintf("/api/v1/jobs/%s/apply", jobID), nil, nil)
}

func (c *Client) GetApplications() ([]map[string]interface{}, error) {
	var apps []map[string]interface{}
	err := c.request("GET", "/api/v1/applications", nil, &apps)
	return apps, err
}

func (c *Client) GetProfile() (map[string]interface{}, error) {
	var profile map[string]interface{}
	err := c.request("GET", "/api/v1/profile", nil, &profile)
	return profile, err
}

func (c *Client) UpdateProfile(data map[string]interface{}) error {
	return c.request("PUT", "/api/v1/profile", data, nil)
}

func (c *Client) GetSettings() (map[string]interface{}, error) {
	var settings map[string]interface{}
	err := c.request("GET", "/api/v1/settings", nil, &settings)
	return settings, err
}

func (c *Client) UpdateSettings(data map[string]interface{}) error {
	return c.request("PUT", "/api/v1/settings", data, nil)
}
