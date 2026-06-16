# ACK Scanner

ACK Scanner is a CLI tool that proactively identifies missing `is_document: true` and `is_iam_policy: true` annotations across [AWS Controllers for Kubernetes (ACK)](https://github.com/aws-controllers-k8s) controller repositories.

## The Problem

ACK controllers compare desired state against observed state to detect drift. For fields containing JSON or YAML documents (e.g., IAM policies, SNS filter policies), a simple string comparison produces false positives when semantic content is identical but formatting differs. This causes **infinite reconciliation loops** where controllers continuously update resources that are already in the desired state.

The fix is marking these fields with `is_document: true` or `is_iam_policy: true` in `generator.yaml`, but this is done manually and reactively — only after customers report issues.

## What ACK Scanner Does

1. **Discovers** all `*-controller` repos in the `aws-controllers-k8s` GitHub organization
2. **Parses** Go source files (`apis/{version}/types.go`) to find `*string` fields matching document-string heuristics
3. **Checks** `generator.yaml` files to determine which fields already have annotations
4. **Cross-references** with Terraform AWS provider documentation to validate classifications
5. **Produces** a prioritized gap analysis report

## Installation

```bash
go install github.com/aws-controllers-k8s/ack-scanner@latest
```

Or build from source:

```bash
git clone https://github.com/aws-controllers-k8s/ack-scanner.git
cd ack-scanner
go build -o ack-scanner .
```

## Usage

### List all ACK controllers

```bash
# List controllers in table format (default)
ack-scanner list

# List in JSON format
ack-scanner list --output json
```

### Scan controllers for annotation gaps

```bash
# Scan all controllers
ack-scanner scan

# Scan specific services
ack-scanner scan --services iam,eks,sns

# Scan with verbose logging
ack-scanner scan --verbose
```

### Scan Terraform provider documentation

```bash
# Identify JSON-accepting fields in Terraform docs
ack-scanner terraform-scan

# Output as JSON
ack-scanner terraform-scan --output json
```

### Generate a gap analysis report

```bash
# Full report (default: markdown format)
ack-scanner report

# Report for specific services in JSON
ack-scanner report --services iam,s3,eks --output json

# Report in table format
ack-scanner report --output table
```

## Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--github-token` | GitHub personal access token | `$GITHUB_TOKEN` env var |
| `--cache-dir` | Directory for caching cloned repositories | `$HOME/.ack-scanner/cache` |
| `--output` | Output format: `json`, `table`, or `markdown` | `table` (commands vary) |
| `--verbose` | Enable detailed progress logging to stderr | `false` |

## Authentication

For higher GitHub API rate limits, provide a token:

```bash
# Via environment variable (recommended)
export GITHUB_TOKEN=ghp_your_token_here
ack-scanner list

# Via flag
ack-scanner list --github-token ghp_your_token_here
```

Without a token, the scanner uses unauthenticated access (60 requests/hour limit).

## Report Categories

The gap analysis report classifies each field into one of four categories:

| Category | Description |
|----------|-------------|
| **Gap (Confirmed)** | Unannotated in ACK and confirmed as JSON in Terraform docs |
| **Gap (Unconfirmed)** | Unannotated in ACK but not found in Terraform docs |
| **Annotated** | Already has `is_document` or `is_iam_policy` in generator.yaml |
| **Terraform Only** | Found in Terraform docs but not matched to an ACK field |

## Development

```bash
# Run tests
go test ./...

# Build
go build ./...

# Run with verbose output
go run . scan --verbose --services iam
```

## License

This project is licensed under the Apache License 2.0 — see the [LICENSE](LICENSE) file for details.
