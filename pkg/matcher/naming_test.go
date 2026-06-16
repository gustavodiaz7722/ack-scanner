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
