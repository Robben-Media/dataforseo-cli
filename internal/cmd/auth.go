package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/builtbyrobben/dataforseo-cli/internal/outfmt"
	"github.com/builtbyrobben/dataforseo-cli/internal/secrets"
)

type AuthCmd struct {
	SetKey AuthSetKeyCmd `cmd:"" help:"Set API credentials (uses --stdin by default)"`
	Status AuthStatusCmd `cmd:"" help:"Show authentication status"`
	Remove AuthRemoveCmd `cmd:"" help:"Remove stored credentials"`
}

type AuthSetKeyCmd struct {
	Stdin bool   `help:"Read auth from stdin (default: true)" default:"true"`
	Key   string `arg:"" optional:"" help:"Base64-encoded auth (discouraged; exposes in shell history)"`
}

func (cmd *AuthSetKeyCmd) Run(ctx context.Context) error {
	var apiKey string

	switch {
	case cmd.Key != "":
		fmt.Fprintln(os.Stderr, "Warning: passing credentials as arguments exposes them in shell history. Use --stdin instead.")
		apiKey = strings.TrimSpace(cmd.Key)
	case term.IsTerminal(int(os.Stdin.Fd())):
		fmt.Fprint(os.Stderr, "Enter base64-encoded auth (login:password): ")

		byteKey, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr)

		if err != nil {
			return fmt.Errorf("read auth: %w", err)
		}

		apiKey = strings.TrimSpace(string(byteKey))
	default:
		byteKey, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("read auth from stdin: %w", err)
		}

		apiKey = strings.TrimSpace(string(byteKey))
	}

	if apiKey == "" {
		return fmt.Errorf("auth credential cannot be empty")
	}

	store, err := secrets.OpenDefault()
	if err != nil {
		return fmt.Errorf("open credential store: %w", err)
	}

	if err := store.SetAPIKey(apiKey); err != nil {
		return fmt.Errorf("store credential: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "Credentials stored in keyring",
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "MESSAGE"}, [][]string{{"success", "Credentials stored in keyring"}})
	}

	fmt.Fprintln(os.Stderr, "Credentials stored in keyring")

	return nil
}

type AuthStatusCmd struct{}

func (cmd *AuthStatusCmd) Run(ctx context.Context) error {
	store, err := secrets.OpenDefault()
	if err != nil {
		return fmt.Errorf("open credential store: %w", err)
	}

	hasKey, err := store.HasKey()
	if err != nil {
		return fmt.Errorf("check credentials: %w", err)
	}

	envKey := os.Getenv("DATAFORSEO_AUTH")
	envOverride := envKey != ""

	status := map[string]any{
		"has_key":         hasKey,
		"env_override":    envOverride,
		"storage_backend": "keyring",
	}

	if hasKey && !envOverride {
		key, err := store.GetAPIKey()
		if err == nil && len(key) > 8 {
			status["key_redacted"] = key[:4] + "..." + key[len(key)-4:]
		}
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, status)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"HAS_KEY", "ENV_OVERRIDE", "STORAGE"}
		rows := [][]string{{
			fmt.Sprintf("%t", hasKey),
			fmt.Sprintf("%t", envOverride),
			"keyring",
		}}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	fmt.Fprintf(os.Stdout, "Storage: %s\n", status["storage_backend"])

	switch {
	case envOverride:
		fmt.Fprintln(os.Stdout, "Credentials: Using DATAFORSEO_AUTH environment variable")
	case hasKey:
		fmt.Fprintln(os.Stdout, "Credentials: Authenticated")

		if redacted, ok := status["key_redacted"].(string); ok {
			fmt.Fprintf(os.Stdout, "Key: %s\n", redacted)
		}
	default:
		fmt.Fprintln(os.Stdout, "Credentials: Not authenticated")
		fmt.Fprintln(os.Stderr, "Run: dataforseo-cli auth set-key --stdin")
	}

	return nil
}

type AuthRemoveCmd struct{}

func (cmd *AuthRemoveCmd) Run(ctx context.Context) error {
	store, err := secrets.OpenDefault()
	if err != nil {
		return fmt.Errorf("open credential store: %w", err)
	}

	if err := store.DeleteAPIKey(); err != nil {
		return fmt.Errorf("remove credentials: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "Credentials removed",
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "MESSAGE"}, [][]string{{"success", "Credentials removed"}})
	}

	fmt.Fprintln(os.Stderr, "Credentials removed")

	return nil
}
