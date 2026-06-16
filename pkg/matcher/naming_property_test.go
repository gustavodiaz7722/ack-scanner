package matcher

import (
	"strings"
	"testing"
	"unicode"

	"pgregory.net/rapid"
)

// Feature: ack-scanner, Property 13: CamelCase to snake_case field matching
// **Validates: Requirements 5.1**

// genCamelCaseWord generates a single CamelCase word starting with an ASCII uppercase
// letter followed by zero or more ASCII lowercase letters. This matches Go field naming
// conventions used in ACK controllers.
func genCamelCaseWord() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		upper := rapid.ByteRange('A', 'Z').Draw(t, "upper")
		lowerCount := rapid.IntRange(0, 8).Draw(t, "lowerCount")
		var sb strings.Builder
		sb.WriteByte(upper)
		for i := 0; i < lowerCount; i++ {
			lower := rapid.ByteRange('a', 'z').Draw(t, "lower")
			sb.WriteByte(lower)
		}
		return sb.String()
	})
}

// genCamelCaseString generates a CamelCase string composed of 1-5 words,
// each starting with an uppercase letter followed by lowercase letters.
func genCamelCaseString() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		wordCount := rapid.IntRange(1, 5).Draw(t, "wordCount")
		var sb strings.Builder
		for i := 0; i < wordCount; i++ {
			word := genCamelCaseWord().Draw(t, "word")
			sb.WriteString(word)
		}
		return sb.String()
	})
}

// TestCamelToSnake_OutputContainsOnlyValidChars verifies that the output of
// CamelToSnake contains only lowercase letters, digits, and underscores.
func TestCamelToSnake_OutputContainsOnlyValidChars(t *testing.T) {
	// Feature: ack-scanner, Property 13: CamelCase to snake_case field matching
	rapid.Check(t, func(t *rapid.T) {
		input := genCamelCaseString().Draw(t, "camelCase")
		result := CamelToSnake(input)

		for i, r := range result {
			if !unicode.IsLower(r) && !unicode.IsDigit(r) && r != '_' {
				t.Fatalf("CamelToSnake(%q) = %q, contains invalid char %q at position %d",
					input, result, string(r), i)
			}
		}
	})
}

// TestCamelToSnake_NoConsecutiveUnderscores verifies that CamelToSnake output
// never has consecutive underscores.
func TestCamelToSnake_NoConsecutiveUnderscores(t *testing.T) {
	// Feature: ack-scanner, Property 13: CamelCase to snake_case field matching
	rapid.Check(t, func(t *rapid.T) {
		input := genCamelCaseString().Draw(t, "camelCase")
		result := CamelToSnake(input)

		if strings.Contains(result, "__") {
			t.Fatalf("CamelToSnake(%q) = %q, contains consecutive underscores",
				input, result)
		}
	})
}

// TestCamelToSnake_NoLeadingOrTrailingUnderscore verifies that CamelToSnake output
// never starts or ends with an underscore for non-empty input.
func TestCamelToSnake_NoLeadingOrTrailingUnderscore(t *testing.T) {
	// Feature: ack-scanner, Property 13: CamelCase to snake_case field matching
	rapid.Check(t, func(t *rapid.T) {
		input := genCamelCaseString().Draw(t, "camelCase")
		result := CamelToSnake(input)

		if result == "" {
			return // empty output is fine for edge cases
		}
		if strings.HasPrefix(result, "_") {
			t.Fatalf("CamelToSnake(%q) = %q, starts with underscore",
				input, result)
		}
		if strings.HasSuffix(result, "_") {
			t.Fatalf("CamelToSnake(%q) = %q, ends with underscore",
				input, result)
		}
	})
}

// TestCamelToSnake_Deterministic verifies that CamelToSnake is a pure function—
// the same input always produces the same output.
func TestCamelToSnake_Deterministic(t *testing.T) {
	// Feature: ack-scanner, Property 13: CamelCase to snake_case field matching
	rapid.Check(t, func(t *rapid.T) {
		input := genCamelCaseString().Draw(t, "camelCase")
		result1 := CamelToSnake(input)
		result2 := CamelToSnake(input)

		if result1 != result2 {
			t.Fatalf("CamelToSnake(%q) is non-deterministic: %q vs %q",
				input, result1, result2)
		}
	})
}

// TestNormalizeFieldName_IdempotentOnSnakeCase verifies that NormalizeFieldName
// is idempotent when applied to snake_case input (output that already went through
// CamelToSnake and has at least one underscore).
func TestNormalizeFieldName_IdempotentOnSnakeCase(t *testing.T) {
	// Feature: ack-scanner, Property 13: CamelCase to snake_case field matching
	rapid.Check(t, func(t *rapid.T) {
		// Generate a CamelCase string with at least 2 words to ensure underscore in output
		wordCount := rapid.IntRange(2, 5).Draw(t, "wordCount")
		var sb strings.Builder
		for i := 0; i < wordCount; i++ {
			word := genCamelCaseWord().Draw(t, "word")
			sb.WriteString(word)
		}
		camel := sb.String()

		// Convert to snake_case first
		snake := CamelToSnake(camel)

		// NormalizeFieldName on snake_case input should strip prefixes/suffixes
		// but applying it again should be stable (idempotent)
		normalized := NormalizeFieldName(snake)
		normalizedAgain := NormalizeFieldName(normalized)

		if normalized != normalizedAgain {
			t.Fatalf("NormalizeFieldName is not idempotent on snake_case input: "+
				"NormalizeFieldName(%q) = %q, but NormalizeFieldName(%q) = %q",
				snake, normalized, normalized, normalizedAgain)
		}
	})
}

// TestCamelToSnake_KnownMappings verifies specific known CamelCase to snake_case
// mappings that are critical for ACK field matching.
func TestCamelToSnake_KnownMappings(t *testing.T) {
	// Feature: ack-scanner, Property 13: CamelCase to snake_case field matching
	knownMappings := map[string]string{
		"AssumeRolePolicyDocument": "assume_role_policy_document",
		"PolicyDocument":           "policy_document",
		"FilterPolicy":             "filter_policy",
		"ConfigurationValues":      "configuration_values",
		"IAMRole":                  "iam_role",
		"JSONSchema":               "json_schema",
		"HTTPEndpoint":             "http_endpoint",
		"VPCId":                    "vpc_id",
		"ResourceARN":              "resource_arn",
		"Ec2Instance":              "ec2_instance",
		"S3Bucket":                 "s3_bucket",
		"Route53Record":            "route53_record",
	}

	for input, expected := range knownMappings {
		result := CamelToSnake(input)
		if result != expected {
			t.Errorf("CamelToSnake(%q) = %q, want %q", input, result, expected)
		}
	}
}
