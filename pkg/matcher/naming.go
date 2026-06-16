package matcher

import (
	"strings"
	"unicode"
)

// knownAcronyms is a set of common AWS/tech acronyms that should be treated
// as single words during CamelCase to snake_case conversion.
var knownAcronyms = map[string]bool{
	"IAM":   true,
	"JSON":  true,
	"HTTP":  true,
	"HTTPS": true,
	"AWS":   true,
	"API":   true,
	"URL":   true,
	"URI":   true,
	"ARN":   true,
	"VPC":   true,
	"EC2":   true,
	"S3":    true,
	"SNS":   true,
	"SQS":   true,
	"KMS":   true,
	"ECS":   true,
	"EKS":   true,
	"RDS":   true,
	"SSH":   true,
	"SSL":   true,
	"TLS":   true,
	"TCP":   true,
	"UDP":   true,
	"DNS":   true,
	"IP":    true,
	"ID":    true,
	"TTL":   true,
	"ACL":   true,
	"CIDR":  true,
	"XML":   true,
	"YAML":  true,
	"HTML":  true,
	"CSS":   true,
	"SQL":   true,
	"CPU":   true,
	"RAM":   true,
	"IO":    true,
	"OS":    true,
}

// CamelToSnake converts a Go CamelCase field name to snake_case.
// It handles:
//   - Standard CamelCase boundaries (e.g., "PolicyDocument" → "policy_document")
//   - Consecutive uppercase / acronyms (e.g., "IAMRole" → "iam_role")
//   - Numeric characters (e.g., "Ec2Instance" → "ec2_instance")
//   - Known acronyms remain grouped (e.g., "JSONSchema" → "json_schema")
func CamelToSnake(s string) string {
	if s == "" {
		return ""
	}

	// Split the string into word tokens using acronym-aware logic
	words := splitCamelWords(s)

	// Join with underscores and lowercase
	var parts []string
	for _, w := range words {
		if w != "" {
			parts = append(parts, strings.ToLower(w))
		}
	}
	return strings.Join(parts, "_")
}

// splitCamelWords splits a CamelCase string into its constituent words,
// properly handling known acronyms, consecutive uppercase letters, and digits.
// Digits are kept attached to the preceding word segment when they appear
// between lowercase and uppercase boundaries (e.g., "Ec2Instance" → ["Ec2", "Instance"],
// "Route53Record" → ["Route53", "Record"]).
func splitCamelWords(s string) []string {
	var words []string
	runes := []rune(s)
	n := len(runes)
	i := 0

	for i < n {
		// Try to match a known acronym at this position
		if unicode.IsUpper(runes[i]) {
			acronym := matchAcronym(runes, i)
			if acronym != "" {
				words = append(words, acronym)
				i += len([]rune(acronym))
				continue
			}
		}

		// Collect a word
		start := i

		if unicode.IsUpper(runes[i]) {
			// Check for a run of uppercase letters (potential unknown acronym)
			j := i + 1
			for j < n && unicode.IsUpper(runes[j]) {
				j++
			}

			if j-i > 1 {
				// Multiple uppercase letters in a row
				if j < n && unicode.IsLower(runes[j]) {
					// The last uppercase letter starts the next word
					// e.g., "HTMLParser" → ["HTML", "Parser"]
					words = append(words, string(runes[start:j-1]))
					i = j - 1
				} else {
					// All uppercase until end or non-letter
					words = append(words, string(runes[start:j]))
					i = j
				}
			} else {
				// Single uppercase letter starts a normal word
				// Consume lowercase letters and trailing digits as part of this word
				j = i + 1
				for j < n && (unicode.IsLower(runes[j]) || unicode.IsDigit(runes[j])) {
					j++
				}
				// If we ended on digits, also consume any additional digits
				words = append(words, string(runes[start:j]))
				i = j
			}
		} else if unicode.IsDigit(runes[i]) {
			// Digits at the start of a segment — collect digits and any
			// following lowercase letters until next uppercase or end
			j := i + 1
			for j < n && (unicode.IsDigit(runes[j]) || unicode.IsLower(runes[j])) {
				j++
			}
			words = append(words, string(runes[start:j]))
			i = j
		} else {
			// Lowercase start (unusual for CamelCase but handle gracefully)
			// Consume lowercase and digits until next uppercase
			j := i + 1
			for j < n && !unicode.IsUpper(runes[j]) {
				j++
			}
			words = append(words, string(runes[start:j]))
			i = j
		}
	}

	return words
}

// matchAcronym checks if a known acronym starts at position i in the rune slice.
// It returns the acronym string if found (and the character following it is not lowercase,
// or it's at the end of the string), otherwise returns "".
func matchAcronym(runes []rune, i int) string {
	remaining := string(runes[i:])

	// Try longer acronyms first to avoid partial matches
	for length := 5; length >= 2; length-- {
		if len(remaining) < length {
			continue
		}
		candidate := remaining[:length]
		if knownAcronyms[candidate] {
			// Verify this is a proper boundary:
			// The acronym must be followed by:
			// - end of string, OR
			// - an uppercase letter (start of next word), OR
			// - a digit, OR
			// - a lowercase letter (in which case this is still valid as the
			//   acronym stands alone before the next camelCase word)
			afterIdx := i + length
			if afterIdx >= len(runes) {
				return candidate
			}
			next := runes[afterIdx]
			if unicode.IsUpper(next) || unicode.IsDigit(next) {
				return candidate
			}
			// If followed by lowercase, the acronym is still valid
			// e.g., "IAMrole" is unusual but "IAMRole" is standard
			// For "HTMLParser" → we match "HTML" because next is 'P' (upper)
			// For "HTTPSOnly" → we should match "HTTPS" because next is 'O' (upper)
			// For a case like "IDs" we want "ID" + "s"
			if unicode.IsLower(next) {
				return candidate
			}
		}
	}
	return ""
}

// commonPrefixes are prefixes that may be stripped during normalization
// for cross-reference matching.
var commonPrefixes = []string{
	"spec_",
	"status_",
}

// commonSuffixes are suffixes that may be stripped during normalization
// for cross-reference matching.
var commonSuffixes = []string{
	"_input",
	"_output",
}

// NormalizeFieldName normalizes a field name for comparison by converting
// to lowercase snake_case and stripping common prefixes/suffixes.
// This is used when matching ACK CamelCase fields against Terraform snake_case fields.
func NormalizeFieldName(name string) string {
	// If already snake_case (contains underscore and is all lowercase), just lowercase it
	// Otherwise convert from CamelCase
	var normalized string
	if isSnakeCase(name) {
		normalized = strings.ToLower(name)
	} else {
		normalized = CamelToSnake(name)
	}

	// Strip common prefixes
	for _, prefix := range commonPrefixes {
		normalized = strings.TrimPrefix(normalized, prefix)
	}

	// Strip common suffixes
	for _, suffix := range commonSuffixes {
		normalized = strings.TrimSuffix(normalized, suffix)
	}

	return normalized
}

// isSnakeCase returns true if a string appears to already be in snake_case
// (contains underscores and no uppercase letters).
func isSnakeCase(s string) bool {
	hasUnderscore := false
	for _, r := range s {
		if r == '_' {
			hasUnderscore = true
		}
		if unicode.IsUpper(r) {
			return false
		}
	}
	return hasUnderscore
}
