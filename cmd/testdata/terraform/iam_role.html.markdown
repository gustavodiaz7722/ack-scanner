---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "AWS: aws_iam_role"
description: |-
  Provides an IAM role.
---

# Resource: aws_iam_role

Provides an IAM role.

## Example Usage

```terraform
resource "aws_iam_role" "example" {
  name = "example"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
    ]
  })

  inline_policy {
    name   = "my_inline_policy"
    policy = data.aws_iam_policy_document.inline_policy.json
  }

  tags = {
    tag-key = "tag-value"
  }
}
```

## Argument Reference

This resource supports the following arguments:

* `assume_role_policy` - (Required) Policy document associated with the role as a JSON string.
* `description` - (Optional) Description of the role.
* `force_detach_policies` - (Optional) Whether to force detaching any policies the role has.
* `inline_policy` - (Optional) Configuration block defining an exclusive set of IAM inline policies.
* `managed_policy_arns` - (Optional) Set of exclusive IAM managed policy ARNs to attach to the IAM role.
* `max_session_duration` - (Optional) Maximum session duration (in seconds).
* `name` - (Optional) Friendly name of the role.
* `permissions_boundary` - (Optional) ARN of the policy that is used to set the permissions boundary.
* `tags` - (Optional) Key-value mapping of tags for the IAM role.

## Attribute Reference

This resource exports the following attributes:

* `arn` - Amazon Resource Name (ARN) specifying the role.
* `id` - Name of the role.
* `unique_id` - Stable and unique string identifying the role.
