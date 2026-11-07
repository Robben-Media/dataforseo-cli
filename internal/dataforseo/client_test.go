package dataforseo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/builtbyrobben/dataforseo-cli/internal/api"
)

func TestSearchVolumeTaskErrors(t *testing.T) {
	for _, tc := range []struct {
		name     string
		response string
		want     string
	}{
		{name: "missing tasks", response: `{"tasks":[]}`, want: "search volume: no tasks in response"},
		{name: "failed task", response: `{"tasks":[{"status_code":40501,"status_message":"Invalid Field"}]}`, want: "search volume: task error (40501): Invalid Field"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tc.response))
			}))
			defer server.Close()

			client := NewClient("", api.WithBaseURL(server.URL))

			result, err := client.SearchVolume(context.Background(), []string{"example"}, 2840, "en")
			if err == nil {
				t.Fatal("expected task error")
			}

			if err.Error() != tc.want {
				t.Fatalf("got %q, want %q", err.Error(), tc.want)
			}

			if result != nil {
				t.Fatalf("expected nil result, got %v", result)
			}
		})
	}
}
