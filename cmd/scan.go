package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/cache"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/discovery"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/parser"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/reporter"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

var servicesFilter string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan ACK controllers for document-string fields and annotation gaps",
	Long: `Scan clones or fetches each ACK controller repository, parses Go types to
identify *string fields matching document-string heuristics, and checks
generator.yaml for existing is_document/is_iam_policy annotations.

Outputs a summary of identified fields, annotated fields, and annotation gaps.`,
	RunE: runScan,
}

func init() {
	scanCmd.Flags().StringVar(&servicesFilter, "services", "", "comma-separated list of service names to scan (default: all)")
	rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	client := GetGitHubClient(ctx)

	start := time.Now()

	// Step 1: Discover controllers
	disc := discovery.NewGitHubDiscoverer(client, "aws-controllers-k8s")
	controllers, err := disc.DiscoverControllers(ctx)
	if err != nil {
		return fmt.Errorf("discovering controllers (%w)", err)
	}

	// Step 2: Apply service filter if specified
	if servicesFilter != "" {
		controllers = filterControllers(controllers, servicesFilter)
	}

	// Step 3: Set up cache and parsers
	repoCache := cache.NewRepoCache(cacheDir)
	goParser := &parser.GoASTParser{}
	genParser := &parser.GeneratorParser{}

	var allResults []types.ScanResult

	// Step 4: Process each controller
	for _, ctrl := range controllers {
		if verbose {
			fmt.Fprintf(os.Stderr, "Scanning %s...\n", ctrl.RepoName)
		}

		// Clone/fetch the repository
		repoDir, err := repoCache.EnsureRepo(ctx, "aws-controllers-k8s", ctrl.RepoName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to clone/fetch %s: %v\n", ctrl.RepoName, err)
			continue
		}

		// Find the latest API version directory
		typesFile, err := findTypesFile(repoDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %s — skipping %s\n", err.Error(), ctrl.RepoName)
			continue
		}

		// Parse types.go for *string fields matching heuristics
		fields, err := goParser.ParseTypesFile(typesFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse types.go in %s: %v\n", ctrl.RepoName, err)
			continue
		}

		// Parse generator.yaml for existing annotations
		generatorPath := filepath.Join(repoDir, "generator.yaml")
		genConfig, err := genParser.ParseGeneratorConfig(generatorPath)
		if err != nil {
			// Missing generator.yaml: treat all fields as unannotated
			fmt.Fprintf(os.Stderr, "Warning: no generator.yaml found for %s — treating all fields as unannotated\n", ctrl.RepoName)
			genConfig = &types.GeneratorConfig{
				Resources: make(map[string]types.ResourceConfig),
			}
		}

		// Classify each field against generator config annotations
		results := classifyFields(ctrl, fields, genConfig)
		allResults = append(allResults, results...)
	}

	// Step 5: Output summary to stderr
	annotated := 0
	gaps := 0
	for _, r := range allResults {
		if r.AnnotationType != types.AnnotationNone {
			annotated++
		} else {
			gaps++
		}
	}
	fmt.Fprintf(os.Stderr, "Scanned %d controller(s): %d fields identified, %d annotated, %d gaps\n",
		len(controllers), len(allResults), annotated, gaps)
	if verbose {
		fmt.Fprintf(os.Stderr, "Scan completed in %s\n", time.Since(start).Round(time.Millisecond))
	}

	// Step 6: Output full results via Reporter to stdout
	rep := reporter.NewReporter(outputFormat)
	// Convert ScanResults to MatchResults for the reporter
	matchResults := scanResultsToMatchResults(allResults)
	return rep.GenerateReport(matchResults, os.Stdout)
}

// filterControllers filters the controller list to only include services
// whose name exactly matches one of the comma-separated values.
func filterControllers(controllers []types.ControllerRepo, filter string) []types.ControllerRepo {
	services := strings.Split(filter, ",")
	allowed := make(map[string]bool, len(services))
	for _, s := range services {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			allowed[trimmed] = true
		}
	}

	var filtered []types.ControllerRepo
	for _, ctrl := range controllers {
		if allowed[ctrl.ServiceName] {
			filtered = append(filtered, ctrl)
		}
	}
	return filtered
}

// findTypesFile locates the API types file by finding the latest API version
// directory under apis/ sorted alphabetically.
func findTypesFile(repoDir string) (string, error) {
	apisDir := filepath.Join(repoDir, "apis")
	entries, err := os.ReadDir(apisDir)
	if err != nil {
		return "", fmt.Errorf("no apis/ directory found in repository")
	}

	// Collect version directories
	var versions []string
	for _, entry := range entries {
		if entry.IsDir() {
			versions = append(versions, entry.Name())
		}
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no version directories found under apis/")
	}

	// Sort alphabetically and pick the last (latest) version
	sort.Strings(versions)
	latestVersion := versions[len(versions)-1]

	typesFile := filepath.Join(apisDir, latestVersion, "types.go")
	if _, err := os.Stat(typesFile); os.IsNotExist(err) {
		return "", fmt.Errorf("types.go not found at %s", typesFile)
	}

	return typesFile, nil
}

// classifyFields classifies each identified field as annotated or unannotated
// by checking against the generator config.
func classifyFields(ctrl types.ControllerRepo, fields []types.GoField, config *types.GeneratorConfig) []types.ScanResult {
	var results []types.ScanResult

	for _, field := range fields {
		annotation := lookupAnnotation(field, config)
		results = append(results, types.ScanResult{
			ServiceName:    ctrl.ServiceName,
			RepoName:       ctrl.RepoName,
			ResourceName:   field.StructName,
			FieldName:      field.FieldName,
			FieldPath:      field.FieldPath,
			GoType:         field.GoType,
			AnnotationType: annotation,
		})
	}

	return results
}

// lookupAnnotation checks whether a field is annotated in the generator config.
// It looks up the field by struct name (resource) and field name, checking
// for is_document or is_iam_policy annotations.
func lookupAnnotation(field types.GoField, config *types.GeneratorConfig) types.AnnotationType {
	if config == nil {
		return types.AnnotationNone
	}

	// Try matching by struct name as resource name
	resource, ok := config.Resources[field.StructName]
	if !ok {
		// Also try with "Spec" suffix stripped (e.g., "BucketSpec" → "Bucket")
		baseName := strings.TrimSuffix(field.StructName, "Spec")
		baseName = strings.TrimSuffix(baseName, "_SDK")
		resource, ok = config.Resources[baseName]
		if !ok {
			return types.AnnotationNone
		}
	}

	// Try exact field name match
	fieldConfig, ok := resource.Fields[field.FieldName]
	if !ok {
		// Try with dot-path notation (e.g., "Spec.PolicyDocument")
		fieldConfig, ok = resource.Fields[field.FieldPath]
		if !ok {
			return types.AnnotationNone
		}
	}

	if fieldConfig.IsIAMPolicy {
		return types.AnnotationIAMPolicy
	}
	if fieldConfig.IsDocument {
		return types.AnnotationDocument
	}
	return types.AnnotationNone
}

// scanResultsToMatchResults converts ScanResult entries to MatchResult entries
// for use by the reporter. Without Terraform data, all unannotated fields are
// classified as gap_without_terraform_confirmation.
func scanResultsToMatchResults(results []types.ScanResult) []types.MatchResult {
	var matchResults []types.MatchResult
	for _, r := range results {
		var category types.Category
		var tfConfirmation types.TerraformConfirmation

		if r.AnnotationType != types.AnnotationNone {
			category = types.CategoryAnnotated
			tfConfirmation = types.TFNotApplicable
		} else {
			category = types.CategoryGapUnconfirmed
			tfConfirmation = types.TFUnconfirmed
		}

		matchResults = append(matchResults, types.MatchResult{
			ServiceName:      r.ServiceName,
			ResourceName:     r.ResourceName,
			FieldName:        r.FieldName,
			FieldPath:        r.FieldPath,
			AnnotationStatus: r.AnnotationType,
			TFConfirmation:   tfConfirmation,
			Category:         category,
		})
	}
	return matchResults
}
