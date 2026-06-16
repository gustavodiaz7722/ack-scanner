package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-github/v60/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// Package-level variables accessible by subcommands.
var (
	githubToken  string
	cacheDir     string
	verbose      bool
	outputFormat string
)

var rootCmd = &cobra.Command{
	Use:   "ack-scanner",
	Short: "ACK Scanner identifies missing is_document annotations across ACK controllers",
	Long: `ACK Scanner proactively identifies missing is_document: true and is_iam_policy: true
annotations across AWS Controllers for Kubernetes (ACK) controller repositories.

It cross-references candidate document-string fields with Terraform AWS provider
documentation to produce a prioritized gap analysis report.`,
	SilenceUsage:      true,
	PersistentPreRunE: persistentPreRun,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&githubToken, "github-token", "", "GitHub personal access token for API authentication")
	rootCmd.PersistentFlags().StringVar(&cacheDir, "cache-dir", "", "directory for caching cloned repositories (default $HOME/.ack-scanner/cache)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable detailed progress logging to stderr")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "output", "table", "output format: json, table, or markdown")
}

// persistentPreRun resolves configuration before any subcommand runs.
func persistentPreRun(cmd *cobra.Command, args []string) error {
	// Resolve GitHub token: flag → env var → unauthenticated with warning.
	if githubToken == "" {
		githubToken = os.Getenv("GITHUB_TOKEN")
	}
	if githubToken == "" {
		fmt.Fprintln(os.Stderr, "Warning: no GitHub token provided (--github-token or GITHUB_TOKEN env var). Rate limits may apply.")
	}

	// Resolve cache directory: flag → default ($HOME/.ack-scanner/cache).
	if cacheDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to determine home directory (%w)", err)
		}
		cacheDir = filepath.Join(home, ".ack-scanner", "cache")
	}

	// Create cache directory if it doesn't exist.
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return fmt.Errorf("failed to create cache directory (%s: %w)", cacheDir, err)
	}

	// Validate --output flag.
	switch outputFormat {
	case "json", "table", "markdown":
		// valid
	default:
		return fmt.Errorf("invalid --output value %q (must be one of json, table, markdown)", outputFormat)
	}

	return nil
}

// GetGitHubClient returns a GitHub client configured with the resolved token.
// If no token was provided, the client is unauthenticated.
func GetGitHubClient(ctx context.Context) *github.Client {
	if githubToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		return github.NewClient(tc)
	}
	return github.NewClient(nil)
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
