package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type issueRequest struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

type issueResponse struct {
	HTMLURL string `json:"html_url"`
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Labels  []struct {
		Name string `json:"name"`
	} `json:"labels"`
}

var httpClient = http.DefaultClient
var baseURL = "https://api.github.com"

func CreateIssue(title, body string, labels []string) (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	owner := os.Getenv("GITHUB_OWNER")
	repo := os.Getenv("GITHUB_REPO")

	url := fmt.Sprintf("%s/repos/%s/%s/issues", baseURL, owner, repo)

	payload, err := json.Marshal(issueRequest{
		Title:  title,
		Body:   body,
		Labels: labels,
	})
	if err != nil {
		return "", fmt.Errorf("marshal issue request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+v3+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("github API responded with %d: %s", resp.StatusCode, string(respBody))
	}

	var issue issueResponse
	if err := json.Unmarshal(respBody, &issue); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	return issue.HTMLURL, nil
}

type Issue struct {
	Number int
	Title  string
	Labels []string
}

func ListIssues(limit int) ([]Issue, error) {
	token := os.Getenv("GITHUB_TOKEN")
	owner := os.Getenv("GITHUB_OWNER")
	repo := os.Getenv("GITHUB_REPO")

	url := fmt.Sprintf("%s/repos/%s/%s/issues?state=open&per_page=%d&sort=created&direction=desc", baseURL, owner, repo, limit)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+v3+json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API responded with %d: %s", resp.StatusCode, string(respBody))
	}

	var issues []issueResponse
	if err := json.Unmarshal(respBody, &issues); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	result := make([]Issue, len(issues))
	for i, issue := range issues {
		labels := make([]string, len(issue.Labels))
		for j, label := range issue.Labels {
			labels[j] = label.Name
		}
		result[i] = Issue{
			Number: issue.Number,
			Title:  issue.Title,
			Labels: labels,
		}
	}

	return result, nil
}
