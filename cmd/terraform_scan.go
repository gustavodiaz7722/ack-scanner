package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/cache"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/parser"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

var terraformScanCmd = &cobra.Command{
	Use:   "terraform-scan",
	Short: "Scan Terraform AWS provider documentation for JSON-accepting fields",
	Long: `Performs a sparse clone of hashicorp/terraform-provider-aws (website/docs/r/ only)
and parses resource documentation to identify fields that accept JSON strings.

Detection uses two signals:
  1. Description phrases in Argument Reference (e.g., "JSON String", "JSON formatted")
  2. jsonencode() usage in Example Usage code blocks

Results are output as a JSON array (--output json) or a formatted table (default).`,
	RunE: runTerraformScanCmd,
}

func init() {
	rootCmd.AddCommand(terraformScanCmd)
}

func runTerraformScanCmd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	start := time.Now()

	// Create the repo cache.
	repoCache := cache.NewRepoCache(cacheDir)

	// Sparse clone hashicorp/terraform-provider-aws with only website/docs/r/.
	if verbose {
		fmt.Fprintln(os.Stderr, "Cloning hashicorp/terraform-provider-aws (sparse: website/docs/r/)...")
	}
	repoDir, err := repoCache.EnsureSparseRepo(ctx, "hashicorp", "terraform-provider-aws", []string{"website/docs/r/"})
	if err != nil {
		return fmt.Errorf("failed to clone Terraform provider repository (%w)", err)
	}

	// Parse all docs in the website/docs/r/ directory.
	docsDir := filepath.Join(repoDir, "website", "docs", "r")
	tfParser := &parser.TerraformParser{}

	if verbose {
		fmt.Fprintln(os.Stderr, "Parsing Terraform resource documentation...")
	}
	fields, err := tfParser.ParseAllDocs(docsDir)
	if err != nil {
		return fmt.Errorf("failed to parse Terraform documentation (%w)", err)
	}

	// Group results by detection method and count.
	var descCount, jsonencodeCount, bothCount int
	for _, f := range fields {
		switch f.DetectionMethod {
		case types.DetectDescriptionPhrase:
			descCount++
		case types.DetectJsonEncodeExample:
			jsonencodeCount++
		case types.DetectBoth:
			bothCount++
		}
	}

	// Output summary to stderr.
	fmt.Fprintf(os.Stderr, "Terraform scan complete:\n")
	fmt.Fprintf(os.Stderr, "  Total JSON fields: %d\n", len(fields))
	fmt.Fprintf(os.Stderr, "  Detected by description phrase: %d\n", descCount)
	fmt.Fprintf(os.Stderr, "  Detected by jsonencode example: %d\n", jsonencodeCount)
	fmt.Fprintf(os.Stderr, "  Detected by both: %d\n", bothCount)
	if verbose {
		fmt.Fprintf(os.Stderr, "  Completed in %s\n", time.Since(start).Round(time.Millisecond))
	}

	// Output full results to stdout.
	switch outputFormat {
	case "json":
		return formatTerraformFieldsJSON(fields, os.Stdout)
	default:
		return formatTerraformFieldsTable(fields, os.Stdout)
	}
}

// formatTerraformFieldsJSON outputs TerraformField results as a JSON array.
func formatTerraformFieldsJSON(fields []types.TerraformField, w io.Writer) error {
	if fields == nil {
		fields = []types.TerraformField{}
	}

	data, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))
	return err
}

// formatTerraformFieldsTable outputs TerraformField results as a column-aligned table.
func formatTerraformFieldsTable(fields []types.TerraformField, w io.Writer) error {
	table := tablewriter.NewTable(w,
		tablewriter.WithHeaderAutoFormat(tw.Off),
	)
	table.Header([]string{"SERVICE", "RESOURCE", "FIELD", "DETECTION METHOD"})

	for _, f := range fields {
		table.Append([]string{
			f.ServiceName,
			f.ResourceType,
			f.FieldName,
			string(f.DetectionMethod),
		})
	}

	return table.Render()
}
