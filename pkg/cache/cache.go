// Package cache manages local clones of git repositories.
package cache

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RepoCache manages local clones of git repositories.
type RepoCache struct {
	baseDir      string
	forceRefresh bool
}

// NewRepoCache creates a new RepoCache that stores clones under baseDir.
func NewRepoCache(baseDir string) *RepoCache {
	return &RepoCache{baseDir: baseDir}
}

// SetForceRefresh enables force-refresh mode, which always fetches latest
// changes from remote even if a local clone already exists.
func (c *RepoCache) SetForceRefresh(refresh bool) {
	c.forceRefresh = refresh
}

// EnsureRepo clones or fetches a repository, returning the local path.
// If the repo already exists locally (has a .git directory), it fetches
// the latest changes and resets to origin/HEAD. Otherwise it performs
// a fresh clone.
func (c *RepoCache) EnsureRepo(ctx context.Context, owner, repo string) (string, error) {
	repoDir := filepath.Join(c.baseDir, owner, repo)

	if isGitRepo(repoDir) {
		// Existing clone: fetch and reset to latest
		if err := runGit(ctx, repoDir, "fetch", "origin"); err != nil {
			return "", fmt.Errorf("git fetch failed for %s/%s: %w", owner, repo, err)
		}
		branch, err := defaultBranch(ctx, repoDir)
		if err != nil {
			return "", fmt.Errorf("failed to determine default branch for %s/%s: %w", owner, repo, err)
		}
		if err := runGit(ctx, repoDir, "reset", "--hard", "origin/"+branch); err != nil {
			return "", fmt.Errorf("git reset failed for %s/%s: %w", owner, repo, err)
		}
		return repoDir, nil
	}

	// Fresh clone
	if err := os.MkdirAll(filepath.Dir(repoDir), 0o755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	cloneURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
	if err := runGit(ctx, filepath.Dir(repoDir), "clone", cloneURL, repo); err != nil {
		return "", fmt.Errorf("git clone failed for %s/%s: %w", owner, repo, err)
	}

	return repoDir, nil
}

// EnsureSparseRepo performs a sparse clone checking out only the specified paths.
// Used for large repos like terraform-provider-aws where we only need specific
// directories (e.g., website/docs/r/).
func (c *RepoCache) EnsureSparseRepo(ctx context.Context, owner, repo string, paths []string) (string, error) {
	repoDir := filepath.Join(c.baseDir, owner, repo)

	if isGitRepo(repoDir) {
		// Existing clone: update sparse-checkout paths and fetch latest
		if err := runGit(ctx, repoDir, sparseCheckoutArgs(paths)...); err != nil {
			return "", fmt.Errorf("git sparse-checkout set failed for %s/%s: %w", owner, repo, err)
		}
		if err := runGit(ctx, repoDir, "fetch", "origin"); err != nil {
			return "", fmt.Errorf("git fetch failed for %s/%s: %w", owner, repo, err)
		}
		branch, err := defaultBranch(ctx, repoDir)
		if err != nil {
			return "", fmt.Errorf("failed to determine default branch for %s/%s: %w", owner, repo, err)
		}
		if err := runGit(ctx, repoDir, "reset", "--hard", "origin/"+branch); err != nil {
			return "", fmt.Errorf("git reset failed for %s/%s: %w", owner, repo, err)
		}
		return repoDir, nil
	}

	// Fresh sparse clone
	if err := os.MkdirAll(filepath.Dir(repoDir), 0o755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	cloneURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
	if err := runGit(ctx, filepath.Dir(repoDir), "clone", "--filter=blob:none", "--sparse", cloneURL, repo); err != nil {
		return "", fmt.Errorf("git sparse clone failed for %s/%s: %w", owner, repo, err)
	}

	// Set the sparse-checkout paths
	if err := runGit(ctx, repoDir, sparseCheckoutArgs(paths)...); err != nil {
		return "", fmt.Errorf("git sparse-checkout set failed for %s/%s: %w", owner, repo, err)
	}

	return repoDir, nil
}

// sparseCheckoutArgs builds the argument list for git sparse-checkout set.
func sparseCheckoutArgs(paths []string) []string {
	args := make([]string, 0, 2+len(paths))
	args = append(args, "sparse-checkout", "set")
	args = append(args, paths...)
	return args
}

// defaultBranch returns the name of the remote's default branch (e.g., "main" or "master").
// It first tries origin/HEAD, then falls back to checking remote show output.
func defaultBranch(ctx context.Context, dir string) (string, error) {
	// Try to read the symbolic ref for origin/HEAD
	cmd := exec.CommandContext(ctx, "git", "symbolic-ref", "refs/remotes/origin/HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err == nil {
		// Output is like "refs/remotes/origin/main"
		ref := strings.TrimSpace(string(out))
		parts := strings.Split(ref, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1], nil
		}
	}

	// Fallback: look at which remote branches exist
	cmd = exec.CommandContext(ctx, "git", "branch", "-r")
	cmd.Dir = dir
	out, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to list remote branches: %w", err)
	}

	// Prefer main, then master, then first available
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		if branch == "origin/main" {
			return "main", nil
		}
	}
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		if branch == "origin/master" {
			return "master", nil
		}
	}
	// Use the first remote branch found
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		if strings.HasPrefix(branch, "origin/") && !strings.Contains(branch, "->") {
			return strings.TrimPrefix(branch, "origin/"), nil
		}
	}

	return "", fmt.Errorf("no remote branches found")
}

// isGitRepo checks if the given directory is an existing git repository.
func isGitRepo(dir string) bool {
	gitDir := filepath.Join(dir, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// runGit executes a git command with the given arguments in the specified directory.
func runGit(ctx context.Context, dir string, args ...string) error {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stderr // git output goes to stderr for verbose logging
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
