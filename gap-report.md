# ACK Scanner Gap Analysis Report

## Summary

| Metric | Count |
| --- | --- |
| Total Fields | 45 |
| Already Annotated | 7 |
| Gaps (Terraform Confirmed) | 11 |
| Gaps (Unconfirmed) | 0 |
| Terraform Only | 27 |

## Priority Services

| Service | Confirmed Gaps |
| --- | --- |
| sns | 6 |
| sqs | 2 |
| dynamodb | 1 |
| iam | 1 |
| kms | 1 |

## dynamodb

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Table | ⚠️ ResourcePolicy | none | Yes | Gap (Confirmed) |
| table_item | item | none | N/A | Terraform Only |
| resource_policy | resource_arn | none | N/A | Terraform Only |

## eks

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Addon | ConfigurationValues | is_document | N/A | Annotated |
| PodIdentityAssociation | Policy | is_iam_policy | N/A | Annotated |
| addon | assume_role_policy | none | N/A | Terraform Only |
| cluster | assume_role_policy | none | N/A | Terraform Only |
| fargate_profile | assume_role_policy | none | N/A | Terraform Only |
| node_group | assume_role_policy | none | N/A | Terraform Only |
| pod_identity_association | assume_role_policy | none | N/A | Terraform Only |

## iam

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Policy | ⚠️ PolicyDocument | none | Yes | Gap (Confirmed) |
| Role | AssumeRolePolicyDocument | is_iam_policy | N/A | Annotated |

## kms

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Key | ⚠️ Policy | none | Yes | Gap (Confirmed) |
| grant | assume_role_policy | none | N/A | Terraform Only |

## lambda

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| function | assume_role_policy | none | N/A | Terraform Only |
| permission | assume_role_policy | none | N/A | Terraform Only |
| event_source_mapping | event_record_format | none | N/A | Terraform Only |
| invocation | input | none | N/A | Terraform Only |
| invocation | lifecycle_scope | none | N/A | Terraform Only |
| function | log_format | none | N/A | Terraform Only |
| event_source_mapping | pattern | none | N/A | Terraform Only |
| function | policy | none | N/A | Terraform Only |
| invocation | terraform_key | none | N/A | Terraform Only |

## rds

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| export_task | assume_role_policy | none | N/A | Terraform Only |
| custom_db_engine_version | manifest | none | N/A | Terraform Only |
| custom_db_engine_version | manifest_hash | none | N/A | Terraform Only |
| export_task | policy | none | N/A | Terraform Only |
| integration | policy | none | N/A | Terraform Only |

## s3

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Bucket | Policy | is_iam_policy | N/A | Annotated |
| bucket_notification | assume_role_policy | none | N/A | Terraform Only |
| bucket_replication_configuration | assume_role_policy | none | N/A | Terraform Only |
| object_copy | kms_encryption_context | none | N/A | Terraform Only |
| bucket | routing_rules | none | N/A | Terraform Only |
| object_copy | source | none | N/A | Terraform Only |

## sns

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Subscription | ⚠️ DeliveryPolicy | none | Yes | Gap (Confirmed) |
| Topic | ⚠️ DeliveryPolicy | none | Yes | Gap (Confirmed) |
| Subscription | ⚠️ FilterPolicyScope | none | Yes | Gap (Confirmed) |
| Subscription | ⚠️ Protocol | none | Yes | Gap (Confirmed) |
| Subscription | ⚠️ RawMessageDelivery | none | Yes | Gap (Confirmed) |
| Subscription | ⚠️ RedrivePolicy | none | Yes | Gap (Confirmed) |
| Subscription | FilterPolicy | is_document | N/A | Annotated |
| Topic | Policy | is_iam_policy | N/A | Annotated |

## sqs

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Queue | ⚠️ RedriveAllowPolicy | none | Yes | Gap (Confirmed) |
| Queue | ⚠️ RedrivePolicy | none | Yes | Gap (Confirmed) |
| Queue | Policy | is_iam_policy | N/A | Annotated |

