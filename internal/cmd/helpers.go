package cmd

import (
	"fmt"
	"os"

	"github.com/builtbyrobben/dataforseo-cli/internal/dataforseo"
	"github.com/builtbyrobben/dataforseo-cli/internal/secrets"
)

func getDataForSEOClient() (*dataforseo.Client, error) {
	// 1. Check env var
	if auth := os.Getenv("DATAFORSEO_AUTH"); auth != "" {
		return dataforseo.NewClient(auth), nil
	}

	// 2. Check keyring
	store, err := secrets.OpenDefault()
	if err != nil {
		return nil, fmt.Errorf("open credential store: %w", err)
	}

	key, err := store.GetAPIKey()
	if err != nil {
		return nil, fmt.Errorf("no credentials found; run: dataforseo-cli auth set-key --stdin")
	}

	return dataforseo.NewClient(key), nil
}
