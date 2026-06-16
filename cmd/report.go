package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/cache"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/discovery"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/matcher"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/parser"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/reporter"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

var reportServicesFilter string

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a gap analysis report cross-referencing ACK and Terraform data",
	Long: `Runs both the ACK controller scan and Terraform documentation scan, then
cross-references the results to produce a prioritized gap analysis report.

The report identifies Document_String_Fields missing is_document or is_iam_policy
annotations, validated against Terraform AWS provider documentation. Fields are
categorized as: gap confirmed by Terraform, gap without Terraform confirmation,
already annotated, or Terraform-only.

If no cached scan data exists, both scans are executed automatically.
The default output format for the report command is markdown.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Default output format to "markdown" for report command if user hasn't
		// explicitly set --output.
		if !cmd.Flags().Changed("output") {
			outputFormat = "markdown"
		}
		return nil
	},
	RunE: runReport,
}

func init() {
	reportCmd.Flags().StringVar(&reportServicesFilter, "services", "", "comma-separated list of service names to include in the report (default: all)")
	rootCmd.AddCommand(reportCmd)
}

func runReport(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	start := time.Now()

	// Step 1: Run the ACK scan.
	if verbose {
		fmt.Fprintln(os.Stderr, "Running ACK controller scan...")
	}

	ackResults, err := executeACKScan(ctx)
	if err != nil {
		return fmt.Errorf("ACK scan failed (%w) — no report generated", err)
	}

	// Step 2: Run the Terraform scan.
	if verbose {
		fmt.Fprintln(os.Stderr, "Running Terraform documentation scan...")
	}

	tfFields, err := executeTerraformScan(ctx)
	if err != nil {
		return fmt.Errorf("Terraform scan failed (%w) — no report generated", err)
	}

	// Step 3: Match ACK fields against Terraform fields.
	if verbose {
		fmt.Fprintln(os.Stderr, "Matching ACK fields with Terraform documentation...")
	}

	m := &matcher.Matcher{}
	matchResults := m.Match(ackResults, tfFields)

	// Step 4: Apply service filter if specified.
	if reportServicesFilter != "" {
		services := parseServicesList(reportServicesFilter)
		matchResults = m.FilterByServices(matchResults, services)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Report generated in %s\n", time.Since(start).Round(time.Millisecond))
	}

	// Step 5: Generate report output.
	rep := reporter.NewReporter(outputFormat)
	return rep.GenerateReport(matchResults, os.Stdout)
}

// executeACKScan executes the full ACK controller scan pipeline and returns the
// scan results. It discovers controllers, clones/fetches repos, parses CRDs
// and generator configs, and classifies fields.
func executeACKScan(ctx context.Context) ([]types.ScanResult, error) {
	ghClient := GetGitHubClient(ctx)

	disc := discovery.NewGitHubDiscoverer(ghClient, "aws-controllers-k8s")
	controllers, err := disc.DiscoverControllers(ctx)
	if err != nil {
		return nil, fmt.Errorf("discovering controllers (%w)", err)
	}

	repoCache := cache.NewRepoCache(cacheDir)
	crdParser := &parser.CRDParser{}
	goParser := &parser.GoASTParser{}
	genParser := &parser.GeneratorParser{}

	var allResults []types.ScanResult

	for _, ctrl := range controllers {
		if verbose {
			fmt.Fprintf(os.Stderr, "  Scanning %s...\n", ctrl.RepoName)
		}

		repoDir, err := repoCache.EnsureRepo(ctx, "aws-controllers-k8s", ctrl.RepoName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to clone/fetch %s: %v\n", ctrl.RepoName, err)
			continue
		}

		// Parse CRD YAML files for string fields under spec (preferred).
		// Falls back to Go AST parsing if no CRDs found.
		fields, err := crdParser.ParseAllCRDs(repoDir)
		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "    No CRDs found, falling back to types.go\n")
			}
			typesFile, err := findTypesFile(repoDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %s — skipping %s\n", err.Error(), ctrl.RepoName)
				continue
			}
			fields, err = goParser.ParseTypesFile(typesFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to parse types.go in %s: %v\n", ctrl.RepoName, err)
				continue
			}
		}

		generatorPath := filepath.Join(repoDir, "generator.yaml")
		genConfig, err := genParser.ParseGeneratorConfig(generatorPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: no generator.yaml found for %s — treating all fields as unannotated\n", ctrl.RepoName)
			genConfig = &types.GeneratorConfig{
				Resources: make(map[string]types.ResourceConfig),
			}
		}

		results := classifyFields(ctrl, fields, genConfig)
		allResults = append(allResults, results...)
	}

	return allResults, nil
}

// executeTerraformScan executes the Terraform documentation scan pipeline and
// returns the identified JSON-accepting fields.
func executeTerraformScan(ctx context.Context) ([]types.TerraformField, error) {
	repoCache := cache.NewRepoCache(cacheDir)

	repoDir, err := repoCache.EnsureSparseRepo(ctx, "hashicorp", "terraform-provider-aws", []string{"website/docs/r/"})
	if err != nil {
		return nil, fmt.Errorf("failed to clone Terraform provider repository (%w)", err)
	}

	docsDir := filepath.Join(repoDir, "website", "docs", "r")
	tfParser := &parser.TerraformParser{}

	fields, err := tfParser.ParseAllDocs(docsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Terraform documentation (%w)", err)
	}

	return fields, nil
}

// parseServicesList splits a comma-separated services string into a trimmed slice.
func parseServicesList(filter string) []string {
	parts := strings.Split(filter, ",")
	var services []string
	for _, s := range parts {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			services = append(services, trimmed)
		}
	}
	return services
}
