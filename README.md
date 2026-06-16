# ACK Scanner

ACK Scanner is a CLI tool that proactively identifies missing `is_document: true` and `is_iam_policy: true` annotations across [AWS Controllers for Kubernetes (ACK)](https://github.com/aws-controllers-k8s) controller repositories.

## The Problem

ACK controllers compare desired state against observed state to detect drift. For fields containing JSON or YAML documents (e.g., IAM policies, SNS filter policies, SQS redrive policies), a simple string comparison produces false positives when semantic content is identical but formatting differs. This causes **infinite reconciliation loops** where controllers continuously update resources that are already in the desired state.

The fix is marking these fields with `is_document: true` or `is_iam_policy: true` in `generator.yaml`, but this is done manually and reactively — only after customers report issues.

## What ACK Scanner Does

1. **Discovers** all `*-controller` repos in the `aws-controllers-k8s` GitHub organization
2. **Parses CRD YAML files** (`config/crd/bases/*.yaml`) to extract all string fields defined under `spec` — this captures every field regardless of which Go source file defines it
3. **Checks** `generator.yaml` files to determine which fields already have annotations
4. **Cross-references with Terraform** AWS provider documentation to identify which CRD string fields accept JSON — Terraform is the source of truth for detection
5. **Produces** a prioritized gap analysis report showing fields that need annotations

## How Detection Works

The scanner uses **Terraform as the source of truth**. It identifies JSON-accepting fields in the Terraform AWS provider documentation using two signals:

1. **Description phrases** in the `## Argument Reference` section (e.g., "JSON String", "JSON formatted", "policy document")
2. **`jsonencode()` usage** in `## Example Usage` code blocks (e.g., `assume_role_policy = jsonencode({...})`)

These Terraform JSON fields are then cross-referenced against every string field in ACK controller CRDs. A field appears as a "Gap (Confirmed)" if Terraform identifies it as JSON but the ACK controller has no `is_document` or `is_iam_policy` annotation.

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
# Full report across all controllers (default: markdown format)
ack-scanner report

# Report for specific services
ack-scanner report --services iam,s3,eks,sns,sqs

# Output as JSON
ack-scanner report --output json

# Save to file
ack-scanner report > gap-report.md
```

## Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--github-token` | GitHub personal access token | `$GITHUB_TOKEN` env var |
| `--cache-dir` | Directory for caching cloned repositories | `$HOME/.ack-scanner/cache` |
| `--output` | Output format: `json`, `table`, or `markdown` | `table` (`report` defaults to `markdown`) |
| `--verbose` | Enable detailed progress logging to stderr | `false` |

## Authentication

For higher GitHub API rate limits, provide a token:

```bash
# Via environment variable (recommended)
export GITHUB_TOKEN=ghp_your_token_here
ack-scanner report

# Via flag
ack-scanner report --github-token ghp_your_token_here
```

Without a token, the scanner uses unauthenticated access (60 requests/hour limit).

## Report Categories

The gap analysis report classifies fields into three categories:

| Category | Description |
|----------|-------------|
| **Gap (Confirmed)** | Unannotated in ACK CRD and confirmed as JSON by Terraform documentation |
| **Annotated** | Already has `is_document` or `is_iam_policy` in generator.yaml |
| **Terraform Only** | Found in Terraform docs but no matching string field exists in the ACK CRD |

Fields that are regular strings (not identified as JSON by Terraform and not already annotated) are excluded from the report entirely.

## Example Output

```
## Priority Services

| Service | Confirmed Gaps |
| --- | --- |
| sns | 6 |
| ecr | 4 |
| wafv2 | 3 |
| sqs | 2 |
| kms | 1 |
| iam | 1 |

## sns

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Subscription | ⚠️ DeliveryPolicy | none | Yes | Gap (Confirmed) |
| Subscription | ⚠️ RedrivePolicy | none | Yes | Gap (Confirmed) |
| Queue | ⚠️ RedriveAllowPolicy | none | Yes | Gap (Confirmed) |
| Subscription | FilterPolicy | is_document | N/A | Annotated |
| Topic | Policy | is_iam_policy | N/A | Annotated |
```

## Architecture

```
ack-scanner/
├── cmd/                        # CLI commands (cobra)
│   ├── root.go                 # Root command, global flags
│   ├── list.go                 # list command
│   ├── scan.go                 # scan command
│   ├── terraform_scan.go       # terraform-scan command
│   └── report.go              # report command
├── pkg/
│   ├── discovery/             # GitHub repository discovery
│   ├── parser/                # CRD, generator.yaml, and Terraform parsers
│   ├── matcher/               # Field matching and cross-referencing
│   ├── cache/                 # Git clone/fetch caching
│   ├── reporter/              # Output formatting (JSON, table, markdown)
│   └── types/                 # Shared data types
└── main.go
```

## Development

```bash
# Run tests
go test ./...

# Build
go build ./...

# Run full report with verbose output
go run . report --verbose

# Run for specific services
go run . report --services iam,sns,sqs --verbose
```

## License

This project is licensed under the Apache License 2.0 — see the [LICENSE](LICENSE) file for details.
