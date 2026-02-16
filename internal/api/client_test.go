package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPost_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.Header.Get("Authorization") != "Basic dGVzdDpwYXNz" {
			t.Errorf("expected Authorization header 'Basic dGVzdDpwYXNz', got %s", r.Header.Get("Authorization"))
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		resp := map[string]any{
			"tasks": []map[string]any{
				{
					"result": []map[string]any{
						{
							"keyword": "seo tools",
							"search_volume": 12000,
						},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAuthFn(func(r *http.Request) {
		r.Header.Set("Authorization", "Basic dGVzdDpwYXNz")
	}))

	var result struct {
		Tasks []struct {
			Result []struct {
				Keyword      string `json:"keyword"`
				SearchVolume int    `json:"search_volume"`
			} `json:"result"`
		} `json:"tasks"`
	}

	err := client.Post(context.Background(), "/keywords_data/google_ads/search_volume/live", []map[string]any{
		{"keywords": []string{"seo tools"}, "location_code": 2840, "language_code": "en"},
	}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Tasks) == 0 {
		t.Fatal("expected at least one task")
	}

	if len(result.Tasks[0].Result) == 0 {
		t.Fatal("expected at least one result")
	}

	if result.Tasks[0].Result[0].Keyword != "seo tools" {
		t.Errorf("expected keyword 'seo tools', got %s", result.Tasks[0].Result[0].Keyword)
	}

	if result.Tasks[0].Result[0].SearchVolume != 12000 {
		t.Errorf("expected search_volume 12000, got %d", result.Tasks[0].Result[0].SearchVolume)
	}
}

func TestPost_HTTPError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAuthFn(func(r *http.Request) {
		r.Header.Set("Authorization", "Basic YmFkOmNyZWRz")
	}))

	var result struct{}

	err := client.Post(context.Background(), "/keywords_data/google_ads/search_volume/live", []map[string]any{
		{"keywords": []string{"test"}},
	}, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.StatusCode != 401 {
		t.Errorf("expected status 401, got %d", apiErr.StatusCode)
	}
}

func TestGet_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		if r.URL.Query().Get("page") != "1" {
			t.Errorf("expected query param page=1, got %s", r.URL.Query().Get("page"))
		}

		resp := map[string]string{"status": "ok"}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)

	var result map[string]string

	err := client.Get(context.Background(), "/status", map[string][]string{"page": {"1"}}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("expected status ok, got %s", result["status"])
	}
}

func TestPost_ErrorMessage(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid keywords parameter"})
	}))
	defer server.Close()

	client := NewClient(server.URL)

	var result struct{}

	err := client.Post(context.Background(), "/test", nil, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.Message != "Invalid keywords parameter" {
		t.Errorf("expected message 'Invalid keywords parameter', got %s", apiErr.Message)
	}
}
