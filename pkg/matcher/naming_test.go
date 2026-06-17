package matcher

import "testing"

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic CamelCase
		{"PolicyDocument", "policy_document"},
		{"AssumeRolePolicyDocument", "assume_role_policy_document"},
		{"FilterPolicy", "filter_policy"},
		{"ConfigurationValues", "configuration_values"},

		// Known acronyms
		{"IAMRole", "iam_role"},
		{"IAMPolicyDocument", "iam_policy_document"},
		{"JSONSchema", "json_schema"},
		{"HTTPEndpoint", "http_endpoint"},
		{"HTTPSEnabled", "https_enabled"},
		{"AWSRegion", "aws_region"},
		{"APIVersion", "api_version"},
		{"VPCId", "vpc_id"},
		{"ARNString", "arn_string"},

		// Acronyms at end
		{"ResourceARN", "resource_arn"},
		{"HostedZoneVPC", "hosted_zone_vpc"},
		{"ClusterIAM", "cluster_iam"},

		// Digits
		{"Ec2Instance", "ec2_instance"},
		{"S3Bucket", "s3_bucket"},
		{"Route53Record", "route53_record"},

		// Single word
		{"Policy", "policy"},
		{"Name", "name"},

		// Empty string
		{"", ""},

		// Multiple consecutive uppercase (unknown acronym)
		{"HTMLParser", "html_parser"},
		{"XMLDocument", "xml_document"},

		// Mixed
		{"EKSAddonConfiguration", "eks_addon_configuration"},
		{"SNSFilterPolicy", "sns_filter_policy"},
		{"KMSKeyId", "kms_key_id"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := CamelToSnake(tc.input)
			if result != tc.expected {
				t.Errorf("CamelToSnake(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestNormalizeFieldName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// CamelCase input
		{"AssumeRolePolicyDocument", "assume_role_policy_document"},
		{"PolicyDocument", "policy_document"},
		{"IAMRole", "iam_role"},

		// Already snake_case
		{"assume_role_policy_document", "assume_role_policy_document"},
		{"policy_document", "policy_document"},
		{"filter_policy", "filter_policy"},

		// With common prefixes to strip
		{"spec_policy_document", "policy_document"},
		{"status_policy_document", "policy_document"},

		// With common suffixes to strip
		{"policy_document_input", "policy_document"},
		{"policy_document_output", "policy_document"},

		// CamelCase that normalizes same as snake_case
		{"FilterPolicy", "filter_policy"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := NormalizeFieldName(tc.input)
			if result != tc.expected {
				t.Errorf("NormalizeFieldName(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"assume_role_policy", true},
		{"policy_document", true},
		{"AssumeRolePolicyDocument", false},
		{"PolicyDocument", false},
		{"policy", false}, // no underscore, not snake_case by our definition
		{"UPPER_CASE", false},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := isSnakeCase(tc.input)
			if result != tc.expected {
				t.Errorf("isSnakeCase(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsSimilarField(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		// Exact match
		{
			name:     "exact match",
			a:        "policy",
			b:        "policy",
			expected: true,
		},

		// Rule 1: prefix at underscore boundary
		{
			name:     "prefix match - assume_role_policy is prefix of assume_role_policy_document",
			a:        "assume_role_policy_document",
			b:        "assume_role_policy",
			expected: true,
		},
		{
			name:     "prefix match reversed args",
			a:        "assume_role_policy",
			b:        "assume_role_policy_document",
			expected: true,
		},

		// Rule 2: suffix at underscore boundary (multi-segment only)
		{
			name:     "multi-segment suffix match - policy_document suffix of assume_role_policy_document",
			a:        "assume_role_policy_document",
			b:        "policy_document",
			expected: true,
		},
		{
			name:     "single-segment suffix should NOT match - scope vs match_scope",
			a:        "match_scope",
			b:        "scope",
			expected: false,
		},
		{
			name:     "single-segment suffix should NOT match - name vs resource_name",
			a:        "resource_name",
			b:        "name",
			expected: false,
		},
		{
			name:     "single-segment suffix should NOT match - type vs content_type",
			a:        "content_type",
			b:        "type",
			expected: false,
		},
		{
			name:     "single-segment suffix should NOT match - policy vs redrive_policy",
			a:        "redrive_policy",
			b:        "policy",
			expected: false,
		},
		{
			name:     "single-segment suffix should NOT match - format vs output_format",
			a:        "output_format",
			b:        "format",
			expected: false,
		},

		// Rule 3: shared suffix of at least 2 segments
		{
			name:     "shared 2-segment suffix - inline_policy_document and role_policy_document",
			a:        "inline_policy_document",
			b:        "role_policy_document",
			expected: true,
		},
		{
			name:     "shared 2-segment suffix - filter_policy_scope and subscription_policy_scope",
			a:        "filter_policy_scope",
			b:        "subscription_policy_scope",
			expected: true,
		},

		// Should NOT match - completely different names
		{
			name:     "no match - completely different",
			a:        "bucket_name",
			b:        "policy_document",
			expected: false,
		},
		{
			name:     "no match - single shared segment suffix",
			a:        "delivery_policy",
			b:        "filter_policy",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isSimilarField(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("isSimilarField(%q, %q) = %v, want %v", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}
