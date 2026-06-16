package cache

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewRepoCache(t *testing.T) {
	c := NewRepoCache("/tmp/test-cache")
	if c == nil {
		t.Fatal("NewRepoCache returned nil")
	}
	if c.baseDir != "/tmp/test-cache" {
		t.Errorf("expected baseDir /tmp/test-cache, got %s", c.baseDir)
	}
}

func TestIsGitRepo(t *testing.T) {
	// Non-existent directory
	if isGitRepo("/nonexistent/path") {
		t.Error("expected false for non-existent directory")
	}

	// Directory without .git
	tmpDir := t.TempDir()
	if isGitRepo(tmpDir) {
		t.Error("expected false for directory without .git")
	}

	// Directory with .git
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.Mkdir(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if !isGitRepo(tmpDir) {
		t.Error("expected true for directory with .git")
	}
}

func TestSparseCheckoutArgs(t *testing.T) {
	args := sparseCheckoutArgs([]string{"website/docs/r/"})
	expected := []string{"sparse-checkout", "set", "website/docs/r/"}
	if len(args) != len(expected) {
		t.Fatalf("expected %d args, got %d", len(expected), len(args))
	}
	for i, v := range expected {
		if args[i] != v {
			t.Errorf("args[%d]: expected %q, got %q", i, v, args[i])
		}
	}
}

func TestSparseCheckoutArgsMultiplePaths(t *testing.T) {
	args := sparseCheckoutArgs([]string{"path/a", "path/b", "path/c"})
	expected := []string{"sparse-checkout", "set", "path/a", "path/b", "path/c"}
	if len(args) != len(expected) {
		t.Fatalf("expected %d args, got %d", len(expected), len(args))
	}
	for i, v := range expected {
		if args[i] != v {
			t.Errorf("args[%d]: expected %q, got %q", i, v, args[i])
		}
	}
}

func TestEnsureRepoCreatesDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	c := NewRepoCache(tmpDir)

	// EnsureRepo with a non-existent repo URL will fail at git clone,
	// but we can verify it creates the parent directories.
	ctx := context.Background()
	_, err := c.EnsureRepo(ctx, "testowner", "testrepo")
	// We expect a failure because there's no real git repo to clone,
	// but verify the parent directory was created.
	if err == nil {
		t.Fatal("expected error for non-existent remote repo")
	}

	parentDir := filepath.Join(tmpDir, "testowner")
	info, statErr := os.Stat(parentDir)
	if statErr != nil {
		t.Fatalf("expected parent directory to be created: %v", statErr)
	}
	if !info.IsDir() {
		t.Error("expected parent path to be a directory")
	}
}

func TestEnsureRepoFetchesExisting(t *testing.T) {
	tmpDir := t.TempDir()
	c := NewRepoCache(tmpDir)

	// Set up a local git repo to simulate an existing cached clone
	repoDir := filepath.Join(tmpDir, "owner", "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Initialize a bare-minimum git repo
	ctx := context.Background()
	if err := runGit(ctx, repoDir, "init"); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	// Create a commit so we have a branch
	dummyFile := filepath.Join(repoDir, "README.md")
	if err := os.WriteFile(dummyFile, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runGit(ctx, repoDir, "add", "."); err != nil {
		t.Fatal(err)
	}
	if err := runGit(ctx, repoDir, "commit", "-m", "init"); err != nil {
		t.Fatal(err)
	}

	// Add a fake origin remote pointing to itself
	if err := runGit(ctx, repoDir, "remote", "add", "origin", repoDir); err != nil {
		t.Fatal(err)
	}

	// Now EnsureRepo should succeed (fetch from itself, reset to origin/HEAD)
	path, err := c.EnsureRepo(ctx, "owner", "repo")
	if err != nil {
		t.Fatalf("EnsureRepo failed on existing repo: %v", err)
	}
	if path != repoDir {
		t.Errorf("expected path %s, got %s", repoDir, path)
	}
}
