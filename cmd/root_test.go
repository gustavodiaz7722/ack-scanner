package cmd

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPersistentPreRun_TokenFromEnv(t *testing.T) {
	// Reset globals.
	githubToken = ""
	cacheDir = t.TempDir()
	outputFormat = "json"

	t.Setenv("GITHUB_TOKEN", "test-token-from-env")

	err := persistentPreRun(rootCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if githubToken != "test-token-from-env" {
		t.Errorf("expected githubToken=%q, got %q", "test-token-from-env", githubToken)
	}
}

func TestPersistentPreRun_TokenFromFlag(t *testing.T) {
	// Flag takes priority over env var.
	githubToken = "flag-token"
	cacheDir = t.TempDir()
	outputFormat = "table"

	t.Setenv("GITHUB_TOKEN", "env-token")

	err := persistentPreRun(rootCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if githubToken != "flag-token" {
		t.Errorf("expected githubToken=%q, got %q", "flag-token", githubToken)
	}
}

func TestPersistentPreRun_NoToken(t *testing.T) {
	// No token should not be an error — just a warning.
	githubToken = ""
	cacheDir = t.TempDir()
	outputFormat = "markdown"

	t.Setenv("GITHUB_TOKEN", "")

	err := persistentPreRun(rootCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if githubToken != "" {
		t.Errorf("expected empty githubToken, got %q", githubToken)
	}
}

func TestPersistentPreRun_DefaultCacheDir(t *testing.T) {
	githubToken = "token"
	cacheDir = ""
	outputFormat = "json"

	err := persistentPreRun(rootCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".ack-scanner", "cache")
	if cacheDir != expected {
		t.Errorf("expected cacheDir=%q, got %q", expected, cacheDir)
	}
}

func TestPersistentPreRun_CacheDirCreation(t *testing.T) {
	githubToken = "token"
	dir := filepath.Join(t.TempDir(), "nested", "cache")
	cacheDir = dir
	outputFormat = "json"

	err := persistentPreRun(rootCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("cache dir not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("cache dir is not a directory")
	}
}

func TestPersistentPreRun_InvalidOutputFormat(t *testing.T) {
	githubToken = "token"
	cacheDir = t.TempDir()
	outputFormat = "yaml"

	err := persistentPreRun(rootCmd, nil)
	if err == nil {
		t.Fatal("expected error for invalid output format")
	}

	expected := `invalid --output value "yaml" (must be one of json, table, markdown)`
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestPersistentPreRun_ValidOutputFormats(t *testing.T) {
	for _, format := range []string{"json", "table", "markdown"} {
		t.Run(format, func(t *testing.T) {
			githubToken = "token"
			cacheDir = t.TempDir()
			outputFormat = format

			err := persistentPreRun(rootCmd, nil)
			if err != nil {
				t.Fatalf("unexpected error for format %q: %v", format, err)
			}
		})
	}
}

func TestGetGitHubClient_WithToken(t *testing.T) {
	githubToken = "test-token"
	ctx := context.Background()

	client := GetGitHubClient(ctx)
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestGetGitHubClient_WithoutToken(t *testing.T) {
	githubToken = ""
	ctx := context.Background()

	client := GetGitHubClient(ctx)
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}
