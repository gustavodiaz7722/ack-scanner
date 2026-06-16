package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMatchesHeuristic(t *testing.T) {
	p := &GoASTParser{}

	tests := []struct {
		fieldName string
		want      bool
	}{
		// Should match
		{"AssumeRolePolicyDocument", true},
		{"FilterPolicy", true},
		{"PolicyDocument", true},
		{"ConfigurationValues", true},
		{"LaunchTemplate", true},
		{"SchemaDefinition", true},
		{"JSONSchema", true},
		{"TaskDefinition", true},
		// Case-insensitive matching
		{"policydocument", true},
		{"CONFIGURATION", true},
		// Should NOT match
		{"Name", false},
		{"ARN", false},
		{"Endpoint", false},
		{"Tags", false},
		{"CreationTime", false},
		{"Status", false},
		{"ID", false},
		{"VpcID", false},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			got := p.MatchesHeuristic(tt.fieldName)
			if got != tt.want {
				t.Errorf("MatchesHeuristic(%q) = %v, want %v", tt.fieldName, got, tt.want)
			}
		})
	}
}

func TestParseTypesFile(t *testing.T) {
	// Create a temporary Go source file for testing
	content := `package types

// RoleSpec defines the desired state of Role.
type RoleSpec struct {
	// The trust relationship policy document.
	AssumeRolePolicyDocument *string ` + "`" + `json:"assumeRolePolicyDocument,omitempty"` + "`" + `

	// A description of the role.
	Description *string ` + "`" + `json:"description,omitempty"` + "`" + `

	// The name of the role.
	RoleName *string ` + "`" + `json:"roleName,omitempty"` + "`" + `

	// The JSON policy document.
	PolicyDocument *string ` + "`" + `json:"policyDocument,omitempty"` + "`" + `

	// Maximum session duration.
	MaxSessionDuration *int64 ` + "`" + `json:"maxSessionDuration,omitempty"` + "`" + `

	// Tags for the role.
	Tags []*Tag ` + "`" + `json:"tags,omitempty"` + "`" + `
}

// TopicSpec defines the desired state of Topic.
type TopicSpec struct {
	// The display name.
	DisplayName *string ` + "`" + `json:"displayName,omitempty"` + "`" + `

	// The filter policy for the subscription.
	FilterPolicy *string ` + "`" + `json:"filterPolicy,omitempty"` + "`" + `

	// The schema for validation.
	ValidationSchema *string ` + "`" + `json:"validationSchema,omitempty"` + "`" + `

	// unexported field should be skipped
	internalPolicy *string
}

// EmbeddedStruct has an anonymous embed that should be skipped.
type EmbeddedStruct struct {
	*RoleSpec

	// Template for notifications.
	NotificationTemplate *string ` + "`" + `json:"notificationTemplate,omitempty"` + "`" + `
}
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "types.go")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	p := &GoASTParser{}
	fields, err := p.ParseTypesFile(filePath)
	if err != nil {
		t.Fatalf("ParseTypesFile returned error: %v", err)
	}

	// Expected fields:
	// - RoleSpec.AssumeRolePolicyDocument (*string, matches "Policy" and "Document")
	// - RoleSpec.PolicyDocument (*string, matches "Policy" and "Document")
	// - TopicSpec.FilterPolicy (*string, matches "Policy")
	// - TopicSpec.ValidationSchema (*string, matches "Schema")
	// - EmbeddedStruct.NotificationTemplate (*string, matches "Template")

	// Should NOT include:
	// - RoleSpec.Description (doesn't match heuristic)
	// - RoleSpec.RoleName (doesn't match heuristic)
	// - RoleSpec.MaxSessionDuration (not *string)
	// - RoleSpec.Tags (not *string)
	// - TopicSpec.DisplayName (doesn't match heuristic)
	// - TopicSpec.internalPolicy (unexported)
	// - EmbeddedStruct.*RoleSpec (anonymous/embedded)

	expectedFields := map[string]struct {
		structName string
		fieldName  string
		fieldPath  string
		goType     string
		jsonTag    string
	}{
		"RoleSpec.AssumeRolePolicyDocument": {
			structName: "RoleSpec",
			fieldName:  "AssumeRolePolicyDocument",
			fieldPath:  "RoleSpec.AssumeRolePolicyDocument",
			goType:     "*string",
			jsonTag:    "assumeRolePolicyDocument",
		},
		"RoleSpec.PolicyDocument": {
			structName: "RoleSpec",
			fieldName:  "PolicyDocument",
			fieldPath:  "RoleSpec.PolicyDocument",
			goType:     "*string",
			jsonTag:    "policyDocument",
		},
		"TopicSpec.FilterPolicy": {
			structName: "TopicSpec",
			fieldName:  "FilterPolicy",
			fieldPath:  "TopicSpec.FilterPolicy",
			goType:     "*string",
			jsonTag:    "filterPolicy",
		},
		"TopicSpec.ValidationSchema": {
			structName: "TopicSpec",
			fieldName:  "ValidationSchema",
			fieldPath:  "TopicSpec.ValidationSchema",
			goType:     "*string",
			jsonTag:    "validationSchema",
		},
		"EmbeddedStruct.NotificationTemplate": {
			structName: "EmbeddedStruct",
			fieldName:  "NotificationTemplate",
			fieldPath:  "EmbeddedStruct.NotificationTemplate",
			goType:     "*string",
			jsonTag:    "notificationTemplate",
		},
	}

	if len(fields) != len(expectedFields) {
		t.Errorf("Expected %d fields, got %d", len(expectedFields), len(fields))
		for _, f := range fields {
			t.Logf("  Got: %s.%s (%s, json:%q)", f.StructName, f.FieldName, f.GoType, f.JSONTag)
		}
	}

	for _, f := range fields {
		key := f.FieldPath
		expected, ok := expectedFields[key]
		if !ok {
			t.Errorf("Unexpected field: %s", key)
			continue
		}
		if f.StructName != expected.structName {
			t.Errorf("Field %s: StructName = %q, want %q", key, f.StructName, expected.structName)
		}
		if f.FieldName != expected.fieldName {
			t.Errorf("Field %s: FieldName = %q, want %q", key, f.FieldName, expected.fieldName)
		}
		if f.GoType != expected.goType {
			t.Errorf("Field %s: GoType = %q, want %q", key, f.GoType, expected.goType)
		}
		if f.JSONTag != expected.jsonTag {
			t.Errorf("Field %s: JSONTag = %q, want %q", key, f.JSONTag, expected.jsonTag)
		}
	}
}

func TestParseTypesFile_NoMatchingFields(t *testing.T) {
	content := `package types

type SimpleStruct struct {
	Name *string ` + "`" + `json:"name"` + "`" + `
	ARN  *string ` + "`" + `json:"arn"` + "`" + `
	ID   *string ` + "`" + `json:"id"` + "`" + `
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "types.go")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	p := &GoASTParser{}
	fields, err := p.ParseTypesFile(filePath)
	if err != nil {
		t.Fatalf("ParseTypesFile returned error: %v", err)
	}

	if len(fields) != 0 {
		t.Errorf("Expected 0 fields, got %d", len(fields))
	}
}

func TestParseTypesFile_InvalidFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "invalid.go")
	if err := os.WriteFile(filePath, []byte("this is not valid Go"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	p := &GoASTParser{}
	_, err := p.ParseTypesFile(filePath)
	if err == nil {
		t.Error("Expected error for invalid Go file, got nil")
	}
}

func TestParseTypesFile_NonexistentFile(t *testing.T) {
	p := &GoASTParser{}
	_, err := p.ParseTypesFile("/nonexistent/path/types.go")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestParseTypesFile_FieldWithoutTag(t *testing.T) {
	content := `package types

type MyStruct struct {
	PolicyDocument *string
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "types.go")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	p := &GoASTParser{}
	fields, err := p.ParseTypesFile(filePath)
	if err != nil {
		t.Fatalf("ParseTypesFile returned error: %v", err)
	}

	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	if fields[0].JSONTag != "" {
		t.Errorf("Expected empty JSONTag, got %q", fields[0].JSONTag)
	}
}
