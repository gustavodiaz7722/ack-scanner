package cmd

import (
	"os"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// Feature: ack-scanner, Property 18: Invalid output flag rejection
//
// For any string value provided to the --output flag that is not one of
// "json", "table", or "markdown", the CLI should reject it with an error
// message that lists all valid options.
//
// **Validates: Requirements 6.10**

func TestProperty18_InvalidOutputFlagRejection(t *testing.T) {
	validFormats := []string{"json", "table", "markdown"}

	tmpDir := t.TempDir()

	rapid.Check(t, func(rt *rapid.T) {
		// Generate an arbitrary string that is NOT one of the valid formats.
		candidate := rapid.String().Draw(rt, "outputFormat")

		isValid := false
		for _, v := range validFormats {
			if candidate == v {
				isValid = true
				break
			}
		}

		if isValid {
			// Skip valid values in this property — they are tested separately below.
			return
		}

		// Set up globals for persistentPreRun.
		githubToken = "test-token"
		cacheDir = tmpDir
		outputFormat = candidate

		err := persistentPreRun(rootCmd, nil)

		// Invalid output format must produce an error.
		if err == nil {
			rt.Fatalf("expected error for invalid output format %q, got nil", candidate)
		}

		// The error message must list all valid options.
		errMsg := err.Error()
		for _, v := range validFormats {
			if !strings.Contains(errMsg, v) {
				rt.Fatalf("error message %q does not contain valid option %q", errMsg, v)
			}
		}
	})
}

func TestProperty18_ValidOutputFormatsAccepted_Property(t *testing.T) {
	validFormats := []string{"json", "table", "markdown"}

	tmpDir := t.TempDir()

	rapid.Check(t, func(rt *rapid.T) {
		// Pick one of the valid formats at random.
		format := rapid.SampledFrom(validFormats).Draw(rt, "validFormat")

		githubToken = "test-token"
		cacheDir = tmpDir
		outputFormat = format

		// Unset GITHUB_TOKEN to avoid interference.
		os.Unsetenv("GITHUB_TOKEN")

		err := persistentPreRun(rootCmd, nil)
		if err != nil {
			rt.Fatalf("expected no error for valid output format %q, got: %v", format, err)
		}
	})
}
