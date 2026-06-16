// Package types contains shared data types used across the ACK Scanner.
package types

// ControllerRepo represents a discovered ACK controller repository.
type ControllerRepo struct {
	RepoName    string `json:"repository_name"`
	ServiceName string `json:"service_name"`
}

// GoField represents a parsed *string field from a Go types.go file.
type GoField struct {
	StructName string `json:"struct_name"`
	FieldName  string `json:"field_name"`
	FieldPath  string `json:"field_path"` // dot-separated from struct root
	GoType     string `json:"go_type"`    // always "*string" for our purposes
	JSONTag    string `json:"json_tag"`
}

// AnnotationType represents the type of annotation on a field.
type AnnotationType string

const (
	AnnotationNone      AnnotationType = "none"
	AnnotationDocument  AnnotationType = "document"
	AnnotationIAMPolicy AnnotationType = "iam_policy"
)

// AnnotatedField represents a field's annotation status from generator.yaml.
type AnnotatedField struct {
	ResourceName   string         `json:"resource_name"`
	FieldName      string         `json:"field_name"`
	AnnotationType AnnotationType `json:"annotation_type"`
}

// GeneratorConfig represents parsed generator.yaml content relevant to scanning.
type GeneratorConfig struct {
	Resources map[string]ResourceConfig `yaml:"resources"`
}

// ResourceConfig holds field configurations for a resource.
type ResourceConfig struct {
	Fields map[string]FieldConfig `yaml:"fields"`
}

// FieldConfig holds annotation settings for a single field.
type FieldConfig struct {
	IsDocument  bool `yaml:"is_document"`
	IsIAMPolicy bool `yaml:"is_iam_policy"`
}

// ScanResult represents the scan result for a single field in a controller.
type ScanResult struct {
	ServiceName    string         `json:"service_name"`
	RepoName       string         `json:"repository_name"`
	ResourceName   string         `json:"resource_name"`
	FieldName      string         `json:"field_name"`
	FieldPath      string         `json:"field_path"`
	GoType         string         `json:"go_type"`
	AnnotationType AnnotationType `json:"annotation_type"`
}

// TerraformField represents a JSON-accepting field from Terraform docs.
type TerraformField struct {
	ServiceName     string          `json:"service_name"`
	ResourceType    string          `json:"resource_type"`
	FieldName       string          `json:"field_name"`
	Description     string          `json:"description"`
	DetectionMethod DetectionMethod `json:"detection_method"`
}

// DetectionMethod indicates how a Terraform JSON field was identified.
type DetectionMethod string

const (
	DetectDescriptionPhrase DetectionMethod = "description_phrase"
	DetectJsonEncodeExample DetectionMethod = "jsonencode_example"
	DetectBoth              DetectionMethod = "both"
)

// TerraformConfirmation represents whether a field was confirmed by Terraform.
type TerraformConfirmation string

const (
	TFConfirmed     TerraformConfirmation = "confirmed"
	TFUnconfirmed   TerraformConfirmation = "unconfirmed"
	TFNotApplicable TerraformConfirmation = "not-applicable"
)

// Category represents the report classification of a field.
type Category string

const (
	CategoryGapConfirmed   Category = "gap_confirmed_by_terraform"
	CategoryGapUnconfirmed Category = "gap_without_terraform_confirmation"
	CategoryAnnotated      Category = "already_annotated"
	CategoryTerraformOnly  Category = "terraform_only"
)

// MatchResult represents a fully classified field in the gap report.
type MatchResult struct {
	ServiceName      string                `json:"service_name"`
	ResourceName     string                `json:"resource_name"`
	FieldName        string                `json:"field_name"`
	FieldPath        string                `json:"field_path"`
	AnnotationStatus AnnotationType        `json:"annotation_status"`
	TFConfirmation   TerraformConfirmation `json:"terraform_confirmation"`
	Category         Category              `json:"category"`
}

// ReportSummary contains aggregate statistics for the report.
type ReportSummary struct {
	TotalFields         int               `json:"total_fields"`
	AnnotatedCount      int               `json:"annotated_count"`
	GapConfirmedCount   int               `json:"gap_confirmed_count"`
	GapUnconfirmedCount int               `json:"gap_unconfirmed_count"`
	TerraformOnlyCount  int               `json:"terraform_only_count"`
	GapsPerService      map[string]int    `json:"gaps_per_service"`
	ServicesByPriority  []ServicePriority `json:"services_by_priority"`
}

// ServicePriority represents a service ranked by Terraform-confirmed gaps.
type ServicePriority struct {
	ServiceName       string `json:"service_name"`
	ConfirmedGapCount int    `json:"confirmed_gap_count"`
}
