package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/discovery"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/reporter"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all ACK controller repositories",
	Long: `Discover and list all ACK controller repositories in the aws-controllers-k8s
GitHub organization. Displays non-archived, non-fork repositories matching the
*-controller naming pattern.`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	start := time.Now()

	client := GetGitHubClient(ctx)
	discoverer := discovery.NewGitHubDiscoverer(client, "aws-controllers-k8s")

	controllers, err := discoverer.DiscoverControllers(ctx)
	if err != nil {
		return fmt.Errorf("failed to discover controllers (%w)", err)
	}

	if len(controllers) == 0 {
		fmt.Fprintln(os.Stderr, "No controllers found.")
	} else {
		fmt.Fprintf(os.Stderr, "Found %d controller(s)\n", len(controllers))
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Discovery completed in %s\n", time.Since(start).Round(time.Millisecond))
	}

	rep := reporter.NewReporter(outputFormat)
	return rep.FormatControllerList(controllers, os.Stdout)
}
