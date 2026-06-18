package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/reporter"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

// allACKControllers is the definitive list of all ACK controllers in the
// aws-controllers-k8s org as of June 2026. This list should be updated
// when new controllers are added or removed.
var allACKControllers = []string{
	"acm",
	"acmpca",
	"amg",
	"apigateway",
	"apigatewayv2",
	"applicationautoscaling",
	"athena",
	"autoscaling",
	"backup",
	"batch",
	"bedrock",
	"bedrockagent",
	"bedrockagentcorecontrol",
	"billing",
	"budgets",
	"cloudfront",
	"cloudtrail",
	"cloudwatch",
	"cloudwatchevents",
	"cloudwatchlogs",
	"codeartifact",
	"cognitoidentity",
	"cognitoidentityprovider",
	"dms",
	"documentdb",
	"drs",
	"dsql",
	"dynamodb",
	"ebs",
	"ec2",
	"ecr",
	"ecrpublic",
	"ecs",
	"efs",
	"eks",
	"elasticache",
	"elbv2",
	"emrcontainers",
	"emrserverless",
	"eventbridge",
	"firehose",
	"globalaccelerator",
	"glue",
	"iam",
	"kafka",
	"keyspaces",
	"kinesis",
	"kms",
	"lambda",
	"memorydb",
	"mq",
	"mwaa",
	"mwaaserverless",
	"networkfirewall",
	"networkmanager",
	"opensearchserverless",
	"opensearchservice",
	"organizations",
	"pipes",
	"prometheusservice",
	"quicksight",
	"ram",
	"rds",
	"recyclebin",
	"rolesanywhere",
	"route53",
	"route53resolver",
	"s3",
	"s3control",
	"s3files",
	"s3tables",
	"s3vectors",
	"sagemaker",
	"secretsmanager",
	"servicecatalog",
	"ses",
	"sfn",
	"sns",
	"sqs",
	"ssm",
	"ssoadmin",
	"transfer",
	"wafv2",
}

// TestReportIncludesAllControllersWithResources verifies that the report output
// includes a section for every controller that has CRD resources, and excludes
// controllers without CRD resources.
func TestReportIncludesAllControllersWithResources(t *testing.T) {
	// Build controllerResources map — controllers with resources should appear,
	// controllers with empty resources should be excluded.
	controllerResources := make(map[string][]string)
	for _, svc := range allACKControllers {
		// Give every controller at least one resource
		controllerResources[svc] = []string{"SomeResource"}
	}
	// Override a few with actual resources
	controllerResources["eks"] = []string{"Addon", "Cluster", "FargateProfile", "Nodegroup"}
	controllerResources["s3"] = []string{"Bucket"}
	controllerResources["sqs"] = []string{"Queue"}

	// Create some sample results for a subset of controllers
	results := []types.MatchResult{
		{
			ServiceName:      "eks",
			ResourceName:     "Addon",
			FieldName:        "ConfigurationValues",
			FieldPath:        "Addon.ConfigurationValues",
			AnnotationStatus: types.AnnotationDocument,
			TFConfirmation:   types.TFNotApplicable,
			Category:         types.CategoryAnnotated,
		},
		{
			ServiceName:      "sqs",
			ResourceName:     "Queue",
			FieldName:        "Policy",
			FieldPath:        "Queue.Policy",
			AnnotationStatus: types.AnnotationNone,
			TFConfirmation:   types.TFConfirmed,
			Category:         types.CategoryGapConfirmed,
		},
	}

	// Generate report
	var buf bytes.Buffer
	rep := reporter.NewReporter("markdown")
	err := rep.GenerateReportWithResources(results, controllerResources, &buf)
	if err != nil {
		t.Fatalf("GenerateReportWithResources failed: %v", err)
	}

	output := buf.String()

	// Verify every controller with resources appears as a ## heading
	for _, svc := range allACKControllers {
		heading := "## " + svc
		if !strings.Contains(output, heading) {
			t.Errorf("report missing controller section: %s", svc)
		}
	}
}

// TestReportExcludesControllersWithoutCRDs verifies that controllers with
// no CRD resources are excluded from the report entirely.
func TestReportExcludesControllersWithoutCRDs(t *testing.T) {
	controllerResources := map[string][]string{
		"eks":            {"Addon", "Cluster"},
		"ssoadmin":       {}, // no CRDs — should be excluded
		"mwaaserverless": {}, // no CRDs — should be excluded
	}

	results := []types.MatchResult{
		{
			ServiceName:      "eks",
			ResourceName:     "Addon",
			FieldName:        "ConfigurationValues",
			FieldPath:        "Addon.ConfigurationValues",
			AnnotationStatus: types.AnnotationDocument,
			TFConfirmation:   types.TFNotApplicable,
			Category:         types.CategoryAnnotated,
		},
	}

	var buf bytes.Buffer
	rep := reporter.NewReporter("markdown")
	err := rep.GenerateReportWithResources(results, controllerResources, &buf)
	if err != nil {
		t.Fatalf("GenerateReportWithResources failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "## eks") {
		t.Error("report missing eks section (has CRD resources)")
	}
	if strings.Contains(output, "## ssoadmin") {
		t.Error("report should NOT contain ssoadmin section (no CRD resources)")
	}
	if strings.Contains(output, "## mwaaserverless") {
		t.Error("report should NOT contain mwaaserverless section (no CRD resources)")
	}
}

// TestReportEmptyTableForNoFindings verifies that controllers with no findings
// output an empty table with the "No findings" message.
func TestReportEmptyTableForNoFindings(t *testing.T) {
	controllerResources := map[string][]string{
		"dynamodb": {"Table", "GlobalTable", "Backup"},
		"sqs":      {"Queue"},
	}

	// Only sqs has findings
	results := []types.MatchResult{
		{
			ServiceName:      "sqs",
			ResourceName:     "Queue",
			FieldName:        "Policy",
			FieldPath:        "Queue.Policy",
			AnnotationStatus: types.AnnotationNone,
			TFConfirmation:   types.TFConfirmed,
			Category:         types.CategoryGapConfirmed,
		},
	}

	var buf bytes.Buffer
	rep := reporter.NewReporter("markdown")
	err := rep.GenerateReportWithResources(results, controllerResources, &buf)
	if err != nil {
		t.Fatalf("GenerateReportWithResources failed: %v", err)
	}

	output := buf.String()

	// dynamodb should have "No findings" message
	if !strings.Contains(output, "## dynamodb") {
		t.Error("report missing dynamodb section")
	}
	if !strings.Contains(output, "_No findings for this controller._") {
		t.Error("report missing 'No findings' message for dynamodb")
	}

	// dynamodb should list its resources
	if !strings.Contains(output, "Table, GlobalTable, Backup") {
		t.Error("report missing resource list for dynamodb")
	}
}

// TestReportListsResources verifies that the report includes the list of
// CRD resources for each controller.
func TestReportListsResources(t *testing.T) {
	controllerResources := map[string][]string{
		"eks": {"AccessEntry", "Addon", "Cluster", "FargateProfile", "Nodegroup"},
	}

	results := []types.MatchResult{
		{
			ServiceName:      "eks",
			ResourceName:     "Addon",
			FieldName:        "ConfigurationValues",
			FieldPath:        "Addon.ConfigurationValues",
			AnnotationStatus: types.AnnotationDocument,
			TFConfirmation:   types.TFNotApplicable,
			Category:         types.CategoryAnnotated,
		},
	}

	var buf bytes.Buffer
	rep := reporter.NewReporter("markdown")
	err := rep.GenerateReportWithResources(results, controllerResources, &buf)
	if err != nil {
		t.Fatalf("GenerateReportWithResources failed: %v", err)
	}

	output := buf.String()

	// Should contain the resources list
	if !strings.Contains(output, "**Resources:** AccessEntry, Addon, Cluster, FargateProfile, Nodegroup") {
		t.Error("report missing resources list for eks")
	}
}
