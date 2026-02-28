package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type capturedIssueRequest struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

func TestCreateIssue_Success(t *testing.T) {
	var got capturedIssueRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") == "" {
			t.Fatalf("expected Authorization header")
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"html_url":"https://example.com/issue/123"}`))
	}))
	defer server.Close()

	originalBaseURL := baseURL
	originalClient := httpClient
	baseURL = server.URL
	httpClient = server.Client()
	defer func() {
		baseURL = originalBaseURL
		httpClient = originalClient
	}()

	t.Setenv("GITHUB_TOKEN", "token")
	t.Setenv("GITHUB_OWNER", "octo")
	t.Setenv("GITHUB_REPO", "repo")

	url, err := CreateIssue("Bug title", "Bug body")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if url != "https://example.com/issue/123" {
		t.Fatalf("unexpected issue url: %s", url)
	}
	if got.Title != "Bug title" || got.Body != "Bug body" {
		t.Fatalf("unexpected payload: %#v", got)
	}
	if len(got.Labels) != 1 || got.Labels[0] != "bug" {
		t.Fatalf("unexpected labels: %#v", got.Labels)
	}
}

func TestCreateIssue_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"bad"}`))
	}))
	defer server.Close()

	originalBaseURL := baseURL
	originalClient := httpClient
	baseURL = server.URL
	httpClient = server.Client()
	defer func() {
		baseURL = originalBaseURL
		httpClient = originalClient
	}()

	t.Setenv("GITHUB_TOKEN", "token")
	t.Setenv("GITHUB_OWNER", "octo")
	t.Setenv("GITHUB_REPO", "repo")

	_, err := CreateIssue("Bug title", "Bug body")
	if err == nil {
		t.Fatal("expected error")
	}
}
