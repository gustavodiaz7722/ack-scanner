package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCRDFile_NestedStringFields(t *testing.T) {
	// CRD with deeply nested string fields (like lambda EventSourceMapping)
	crdYAML := `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
spec:
  group: lambda.services.k8s.aws
  names:
    kind: EventSourceMapping
  versions:
    - name: v1alpha1
      schema:
        openAPIV3Schema:
          properties:
            spec:
              type: object
              properties:
                functionName:
                  type: string
                batchSize:
                  type: integer
                amazonManagedKafkaEventSourceConfig:
                  type: object
                  properties:
                    consumerGroupID:
                      type: string
                    schemaRegistryConfig:
                      type: object
                      properties:
                        eventRecordFormat:
                          type: string
                        schemaRegistryURI:
                          type: string
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "lambda.services.k8s.aws_eventsourcemappings.yaml")
	if err := os.WriteFile(filePath, []byte(crdYAML), 0644); err != nil {
		t.Fatalf("Failed to write test CRD: %v", err)
	}

	p := &CRDParser{}
	fields, resourceName, err := p.ParseCRDFile(filePath)
	if err != nil {
		t.Fatalf("ParseCRDFile returned error: %v", err)
	}

	if resourceName != "EventSourceMapping" {
		t.Errorf("expected resource name EventSourceMapping, got %s", resourceName)
	}

	// Should find all string fields at all levels:
	// - functionName (top-level)
	// - consumerGroupID (nested 1 level)
	// - eventRecordFormat (nested 2 levels)
	// - schemaRegistryURI (nested 2 levels)
	// Should NOT include batchSize (integer)
	expectedFields := map[string]bool{
		"FunctionName":      false,
		"ConsumerGroupID":   false,
		"EventRecordFormat": false,
		"SchemaRegistryURI": false,
	}

	for _, f := range fields {
		if _, ok := expectedFields[f.FieldName]; ok {
			expectedFields[f.FieldName] = true
		}
		if f.StructName != "EventSourceMapping" {
			t.Errorf("field %s: expected StructName EventSourceMapping, got %s", f.FieldName, f.StructName)
		}
		if f.GoType != "*string" {
			t.Errorf("field %s: expected GoType *string, got %s", f.FieldName, f.GoType)
		}
	}

	for fieldName, found := range expectedFields {
		if !found {
			t.Errorf("expected field %s not found in parsed results", fieldName)
		}
	}

	if len(fields) != 4 {
		t.Errorf("expected 4 fields, got %d", len(fields))
		for _, f := range fields {
			t.Logf("  %s (path: %s)", f.FieldName, f.FieldPath)
		}
	}
}

func TestParseCRDFile_DeduplicatesSameFieldAtMultiplePaths(t *testing.T) {
	// CRD where the same field name appears at multiple nesting levels
	// (like sagemaker ModelPackage with ContentType at 18 different paths)
	crdYAML := `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
spec:
  group: sagemaker.services.k8s.aws
  names:
    kind: ModelPackage
  versions:
    - name: v1alpha1
      schema:
        openAPIV3Schema:
          properties:
            spec:
              type: object
              properties:
                contentType:
                  type: string
                inputConfig:
                  type: object
                  properties:
                    contentType:
                      type: string
                    dataSource:
                      type: string
                outputConfig:
                  type: object
                  properties:
                    contentType:
                      type: string
                    s3URI:
                      type: string
                validationConfig:
                  type: object
                  properties:
                    contentType:
                      type: string
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "sagemaker.services.k8s.aws_modelpackages.yaml")
	if err := os.WriteFile(filePath, []byte(crdYAML), 0644); err != nil {
		t.Fatalf("Failed to write test CRD: %v", err)
	}

	p := &CRDParser{}
	fields, _, err := p.ParseCRDFile(filePath)
	if err != nil {
		t.Fatalf("ParseCRDFile returned error: %v", err)
	}

	// Count how many times ContentType appears
	contentTypeCount := 0
	for _, f := range fields {
		if f.FieldName == "ContentType" {
			contentTypeCount++
		}
	}

	// Should be deduplicated to exactly 1
	if contentTypeCount != 1 {
		t.Errorf("expected ContentType to appear exactly once (deduplicated), got %d", contentTypeCount)
		for _, f := range fields {
			if f.FieldName == "ContentType" {
				t.Logf("  ContentType at path: %s", f.FieldPath)
			}
		}
	}

	// The kept entry should have the shortest path (top-level)
	for _, f := range fields {
		if f.FieldName == "ContentType" {
			expected := "ModelPackage.ContentType"
			if f.FieldPath != expected {
				t.Errorf("expected shortest path %q, got %q", expected, f.FieldPath)
			}
			break
		}
	}

	// Other unique fields should still be present
	uniqueFields := map[string]bool{
		"ContentType": false,
		"DataSource":  false,
		"S3URI":       false,
	}
	for _, f := range fields {
		if _, ok := uniqueFields[f.FieldName]; ok {
			uniqueFields[f.FieldName] = true
		}
	}
	for name, found := range uniqueFields {
		if !found {
			t.Errorf("expected unique field %s not found", name)
		}
	}
}

func TestParseCRDFile_TopLevelOnlyFields(t *testing.T) {
	// Simple CRD with only top-level string fields (backward compatibility)
	crdYAML := `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
spec:
  group: sqs.services.k8s.aws
  names:
    kind: Queue
  versions:
    - name: v1alpha1
      schema:
        openAPIV3Schema:
          properties:
            spec:
              type: object
              properties:
                queueName:
                  type: string
                redrivePolicy:
                  type: string
                redriveAllowPolicy:
                  type: string
                visibilityTimeout:
                  type: integer
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "sqs.services.k8s.aws_queues.yaml")
	if err := os.WriteFile(filePath, []byte(crdYAML), 0644); err != nil {
		t.Fatalf("Failed to write test CRD: %v", err)
	}

	p := &CRDParser{}
	fields, resourceName, err := p.ParseCRDFile(filePath)
	if err != nil {
		t.Fatalf("ParseCRDFile returned error: %v", err)
	}

	if resourceName != "Queue" {
		t.Errorf("expected Queue, got %s", resourceName)
	}

	// Should find 3 string fields (not visibilityTimeout which is integer)
	if len(fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(fields))
		for _, f := range fields {
			t.Logf("  %s", f.FieldName)
		}
	}
}

func TestParseCRDFile_NoSpec(t *testing.T) {
	crdYAML := `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
spec:
  group: test.services.k8s.aws
  names:
    kind: NoSpec
  versions:
    - name: v1alpha1
      schema:
        openAPIV3Schema:
          properties:
            status:
              type: object
              properties:
                arn:
                  type: string
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(filePath, []byte(crdYAML), 0644); err != nil {
		t.Fatalf("Failed to write test CRD: %v", err)
	}

	p := &CRDParser{}
	fields, resourceName, err := p.ParseCRDFile(filePath)
	if err != nil {
		t.Fatalf("ParseCRDFile returned error: %v", err)
	}

	if resourceName != "NoSpec" {
		t.Errorf("expected NoSpec, got %s", resourceName)
	}
	if len(fields) != 0 {
		t.Errorf("expected 0 fields (no spec), got %d", len(fields))
	}
}

func TestParseCRDFile_NestedFieldPaths(t *testing.T) {
	// Verify the full dot-separated path is correctly constructed
	crdYAML := `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
spec:
  group: test.services.k8s.aws
  names:
    kind: MyResource
  versions:
    - name: v1alpha1
      schema:
        openAPIV3Schema:
          properties:
            spec:
              type: object
              properties:
                topLevel:
                  type: string
                nested:
                  type: object
                  properties:
                    deepField:
                      type: string
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(filePath, []byte(crdYAML), 0644); err != nil {
		t.Fatalf("Failed to write test CRD: %v", err)
	}

	p := &CRDParser{}
	fields, _, err := p.ParseCRDFile(filePath)
	if err != nil {
		t.Fatalf("ParseCRDFile returned error: %v", err)
	}

	pathMap := make(map[string]string)
	for _, f := range fields {
		pathMap[f.FieldName] = f.FieldPath
	}

	if path, ok := pathMap["TopLevel"]; !ok {
		t.Error("TopLevel field not found")
	} else if path != "MyResource.TopLevel" {
		t.Errorf("TopLevel path: expected MyResource.TopLevel, got %s", path)
	}

	if path, ok := pathMap["DeepField"]; !ok {
		t.Error("DeepField field not found")
	} else if path != "MyResource.Nested.DeepField" {
		t.Errorf("DeepField path: expected MyResource.Nested.DeepField, got %s", path)
	}
}

func TestParseAllCRDs_SkipsServicesCRDs(t *testing.T) {
	tmpDir := t.TempDir()
	basesDir := filepath.Join(tmpDir, "config", "crd", "bases")
	if err := os.MkdirAll(basesDir, 0755); err != nil {
		t.Fatalf("Failed to create bases dir: %v", err)
	}

	// Service CRD that should be skipped
	serviceCRD := `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
spec:
  group: services.k8s.aws
  names:
    kind: FieldExport
  versions:
    - name: v1alpha1
      schema:
        openAPIV3Schema:
          properties:
            spec:
              type: object
              properties:
                fieldPath:
                  type: string
`
	if err := os.WriteFile(filepath.Join(basesDir, "services.k8s.aws_fieldexports.yaml"), []byte(serviceCRD), 0644); err != nil {
		t.Fatal(err)
	}

	// Real resource CRD
	resourceCRD := `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
spec:
  group: sqs.services.k8s.aws
  names:
    kind: Queue
  versions:
    - name: v1alpha1
      schema:
        openAPIV3Schema:
          properties:
            spec:
              type: object
              properties:
                policy:
                  type: string
`
	if err := os.WriteFile(filepath.Join(basesDir, "sqs.services.k8s.aws_queues.yaml"), []byte(resourceCRD), 0644); err != nil {
		t.Fatal(err)
	}

	p := &CRDParser{}
	fields, err := p.ParseAllCRDs(tmpDir)
	if err != nil {
		t.Fatalf("ParseAllCRDs returned error: %v", err)
	}

	// Should only find the Queue CRD's field, not FieldExport
	if len(fields) != 1 {
		t.Fatalf("expected 1 field (services.k8s.aws_ skipped), got %d", len(fields))
	}
	if fields[0].StructName != "Queue" {
		t.Errorf("expected Queue, got %s", fields[0].StructName)
	}
}
