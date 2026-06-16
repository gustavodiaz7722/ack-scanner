# ACK Scanner Gap Analysis Report

## Summary

| Metric | Count |
| --- | --- |
| Total Fields | 248 |
| Already Annotated | 10 |
| Gaps (Terraform Confirmed) | 29 |
| Gaps (Unconfirmed) | 0 |
| Terraform Only | 209 |

## Priority Services

| Service | Confirmed Gaps |
| --- | --- |
| sns | 6 |
| ecr | 4 |
| wafv2 | 3 |
| cloudwatch | 2 |
| glue | 2 |
| sqs | 2 |
| dms | 1 |
| dynamodb | 1 |
| efs | 1 |
| iam | 1 |
| kinesis | 1 |
| kms | 1 |
| opensearchserverless | 1 |
| ram | 1 |
| sagemaker | 1 |
| ssm | 1 |

## acmpca

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| certificate_authority | policy | none | N/A | Terraform Only |
| policy | policy | none | N/A | Terraform Only |

## apigatewayv2

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| model | content_type | none | N/A | Terraform Only |
| authorizer | issuer | none | N/A | Terraform Only |
| model | schema | none | N/A | Terraform Only |

## autoscaling

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| group | notification_metadata | none | N/A | Terraform Only |
| lifecycle_hook | notification_metadata | none | N/A | Terraform Only |

## backup

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| selection | assume_role_policy | none | N/A | Terraform Only |
| report_plan | formats | none | N/A | Terraform Only |
| vault_notifications | policy | none | N/A | Terraform Only |
| vault_policy | policy | none | N/A | Terraform Only |

## bedrockagent

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| agent_action_group | api_schema | none | N/A | Terraform Only |
| agent | assume_role_policy | none | N/A | Terraform Only |
| agent_alias | assume_role_policy | none | N/A | Terraform Only |
| agent_collaborator | assume_role_policy | none | N/A | Terraform Only |
| flow | json | none | N/A | Terraform Only |
| prompt | json | none | N/A | Terraform Only |
| agent_action_group | payload | none | N/A | Terraform Only |
| agent | policy | none | N/A | Terraform Only |
| agent_alias | policy | none | N/A | Terraform Only |
| agent_collaborator | policy | none | N/A | Terraform Only |

## cloudfront

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| realtime_log_config | assume_role_policy | none | N/A | Terraform Only |
| distribution | policy | none | N/A | Terraform Only |
| origin_access_identity | policy | none | N/A | Terraform Only |
| realtime_log_config | policy | none | N/A | Terraform Only |

## cloudtrail

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
|  | policy | none | N/A | Terraform Only |

## cloudwatch

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Dashboard | ⚠️ DashboardBody | none | Yes | Gap (Confirmed) |
| MetricStream | ⚠️ OutputFormat | none | Yes | Gap (Confirmed) |
| log_destination_policy | access_policy | none | N/A | Terraform Only |
| event_target | assume_role_policy | none | N/A | Terraform Only |
| metric_stream | assume_role_policy | none | N/A | Terraform Only |
| event_target | content | none | N/A | Terraform Only |
| log_delivery_destination_policy | delivery_destination_policy | none | N/A | Terraform Only |
| event_archive | event_pattern | none | N/A | Terraform Only |
| event_rule | event_pattern | none | N/A | Terraform Only |
| event_target | event_pattern | none | N/A | Terraform Only |
| event_target | input | none | N/A | Terraform Only |
| event_target | input_path | none | N/A | Terraform Only |
| event_target | input_paths | none | N/A | Terraform Only |
| event_target | input_template | none | N/A | Terraform Only |
| log_transformer | parse_cloudfront | none | N/A | Terraform Only |
| log_transformer | parse_json | none | N/A | Terraform Only |
| log_transformer | parse_postgres | none | N/A | Terraform Only |
| log_transformer | parse_route53 | none | N/A | Terraform Only |
| log_transformer | parse_vpc | none | N/A | Terraform Only |
| log_transformer | parse_waf | none | N/A | Terraform Only |
| event_target | partition_key_path | none | N/A | Terraform Only |
| event_archive | policy | none | N/A | Terraform Only |
| event_bus | policy | none | N/A | Terraform Only |
| event_bus_policy | policy | none | N/A | Terraform Only |
| event_connection | policy | none | N/A | Terraform Only |
| event_rule | policy | none | N/A | Terraform Only |
| event_target | policy | none | N/A | Terraform Only |
| metric_stream | policy | none | N/A | Terraform Only |
| event_bus | policy_document | none | N/A | Terraform Only |
| event_target | policy_document | none | N/A | Terraform Only |
| log_account_policy | policy_document | none | N/A | Terraform Only |
| log_data_protection_policy | policy_document | none | N/A | Terraform Only |
| log_index_policy | policy_document | none | N/A | Terraform Only |
| log_resource_policy | policy_document | none | N/A | Terraform Only |
| contributor_insight_rule | rule_definition | none | N/A | Terraform Only |
| metric_stream | statistics_configuration | none | N/A | Terraform Only |
| log_transformer | target | none | N/A | Terraform Only |

## codeartifact

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| domain_permissions_policy | policy_document | none | N/A | Terraform Only |
| repository_permissions_policy | policy_document | none | N/A | Terraform Only |

## dms

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Endpoint | ⚠️ ExternalTableDefinition | none | Yes | Gap (Confirmed) |
| replication_instance | assume_role_policy | none | N/A | Terraform Only |
| replication_subnet_group | assume_role_policy | none | N/A | Terraform Only |
| endpoint | message_format | none | N/A | Terraform Only |
| replication_config | replication_settings | none | N/A | Terraform Only |
| replication_task | replication_task_settings | none | N/A | Terraform Only |
| endpoint | spatial_data_option_to_geo_json_function_name | none | N/A | Terraform Only |
| replication_config | supplemental_settings | none | N/A | Terraform Only |
| replication_config | table_mappings | none | N/A | Terraform Only |
| replication_task | table_mappings | none | N/A | Terraform Only |

## dsql

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Cluster | Policy | is_iam_policy | N/A | Annotated |

## dynamodb

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Table | ⚠️ ResourcePolicy | none | Yes | Gap (Confirmed) |
| table_item | item | none | N/A | Terraform Only |
| resource_policy | resource_arn | none | N/A | Terraform Only |

## ecr

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Repository | ⚠️ LifecyclePolicy | none | Yes | Gap (Confirmed) |
| RepositoryCreationTemplate | ⚠️ LifecyclePolicy | none | Yes | Gap (Confirmed) |
| Repository | ⚠️ Policy | none | Yes | Gap (Confirmed) |
| RepositoryCreationTemplate | ⚠️ RepositoryPolicy | none | Yes | Gap (Confirmed) |
| pull_time_update_exclusion | assume_role_policy | none | N/A | Terraform Only |

## ecrpublic

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| repository_policy | policy | none | N/A | Terraform Only |

## ecs

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| daemon_task_definition | assume_role_policy | none | N/A | Terraform Only |
| task_definition | container_definitions | none | N/A | Terraform Only |
| service | format | none | N/A | Terraform Only |
| daemon_task_definition | log_driver | none | N/A | Terraform Only |
| cluster | policy | none | N/A | Terraform Only |
| task_definition | secret_string | none | N/A | Terraform Only |

## efs

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| FileSystem | ⚠️ Policy | none | Yes | Gap (Confirmed) |

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

## glue

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Job | ⚠️ SecurityConfiguration | none | Yes | Gap (Confirmed) |
| Job | ⚠️ WorkerType | none | Yes | Gap (Confirmed) |
| connection | SparkProperties | none | N/A | Terraform Only |
| dev_endpoint | assume_role_policy | none | N/A | Terraform Only |
| job | assume_role_policy | none | N/A | Terraform Only |
| classifier | classification | none | N/A | Terraform Only |
| connection | connection_type | none | N/A | Terraform Only |
| schema | data_format | none | N/A | Terraform Only |
| catalog_table | initial_default | none | N/A | Terraform Only |
| classifier | json_classifier | none | N/A | Terraform Only |
| classifier | json_path | none | N/A | Terraform Only |
| connection | secret_string | none | N/A | Terraform Only |
| catalog_table | write_default | none | N/A | Terraform Only |

## iam

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Policy | ⚠️ PolicyDocument | none | Yes | Gap (Confirmed) |
| Role | AssumeRolePolicyDocument | is_iam_policy | N/A | Annotated |

## kinesis

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Stream | ⚠️ Name | none | Yes | Gap (Confirmed) |
| firehose_delivery_stream | assume_role_policy | none | N/A | Terraform Only |
| firehose_delivery_stream | case_insensitive | none | N/A | Terraform Only |
| firehose_delivery_stream | column_to_json_key_mappings | none | N/A | Terraform Only |
| firehose_delivery_stream | convert_dots_in_json_keys_to_underscores | none | N/A | Terraform Only |
| firehose_delivery_stream | data_format_conversion_configuration | none | N/A | Terraform Only |
| firehose_delivery_stream | deserializer | none | N/A | Terraform Only |
| firehose_delivery_stream | file_extension | none | N/A | Terraform Only |
| firehose_delivery_stream | hive_json_ser_de | none | N/A | Terraform Only |
| firehose_delivery_stream | input_format_configuration | none | N/A | Terraform Only |
| analytics_application | json | none | N/A | Terraform Only |
| firehose_delivery_stream | open_x_json_ser_de | none | N/A | Terraform Only |
| resource_policy | policy | none | N/A | Terraform Only |
| analytics_application | record_format_type | none | N/A | Terraform Only |
| firehose_delivery_stream | timestamp_formats | none | N/A | Terraform Only |

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

## networkfirewall

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| resource_policy | policy | none | N/A | Terraform Only |

## networkmanager

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| core_network | create_base_policy | none | N/A | Terraform Only |
| core_network_policy_attachment | policy_document | none | N/A | Terraform Only |
| site_to_site_vpn_attachment | policy_document | none | N/A | Terraform Only |

## opensearchserverless

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| SecurityPolicy | ⚠️ Policy | none | Yes | Gap (Confirmed) |

## opensearchservice

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Domain | AccessPolicies | is_iam_policy | N/A | Annotated |

## organizations

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| policy | content | none | N/A | Terraform Only |
| resource_policy | content | none | N/A | Terraform Only |

## pipes

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| pipe | assume_role_policy | none | N/A | Terraform Only |
| pipe | input_template | none | N/A | Terraform Only |
| pipe | output_format | none | N/A | Terraform Only |
| pipe | pattern | none | N/A | Terraform Only |
| pipe | policy | none | N/A | Terraform Only |
| pipe | time | none | N/A | Terraform Only |
| pipe | timestamp | none | N/A | Terraform Only |

## quicksight

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| data_source | assume_role_policy | none | N/A | Terraform Only |
| vpc_connection | assume_role_policy | none | N/A | Terraform Only |
| data_source | content | none | N/A | Terraform Only |
| data_set | format | none | N/A | Terraform Only |
| data_source | policy | none | N/A | Terraform Only |
| vpc_connection | policy | none | N/A | Terraform Only |

## ram

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Permission | ⚠️ PolicyTemplate | none | Yes | Gap (Confirmed) |

## rds

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| export_task | assume_role_policy | none | N/A | Terraform Only |
| custom_db_engine_version | manifest | none | N/A | Terraform Only |
| custom_db_engine_version | manifest_hash | none | N/A | Terraform Only |
| export_task | policy | none | N/A | Terraform Only |
| integration | policy | none | N/A | Terraform Only |

## route53

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| traffic_policy | document | none | N/A | Terraform Only |
| hosted_zone_dnssec | policy | none | N/A | Terraform Only |
| key_signing_key | policy | none | N/A | Terraform Only |

## s3

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Bucket | Policy | is_iam_policy | N/A | Annotated |
| bucket_notification | assume_role_policy | none | N/A | Terraform Only |
| bucket_replication_configuration | assume_role_policy | none | N/A | Terraform Only |
| object_copy | kms_encryption_context | none | N/A | Terraform Only |
| bucket | routing_rules | none | N/A | Terraform Only |
| object_copy | source | none | N/A | Terraform Only |

## s3control

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| object_lambda_access_point | function_payload | none | N/A | Terraform Only |
| access_grants_instance_resource_policy | policy | none | N/A | Terraform Only |
| access_point_policy | policy | none | N/A | Terraform Only |
| bucket_policy | policy | none | N/A | Terraform Only |
| multi_region_access_point_policy | policy | none | N/A | Terraform Only |
| object_lambda_access_point_policy | policy | none | N/A | Terraform Only |

## s3files

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| FileSystem | Policy | is_iam_policy | N/A | Annotated |

## s3tables

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| table_bucket_policy | resource_policy | none | N/A | Terraform Only |
| table_policy | resource_policy | none | N/A | Terraform Only |

## sagemaker

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Pipeline | ⚠️ PipelineDefinition | none | Yes | Gap (Confirmed) |
| algorithm | assume_role_policy | none | N/A | Terraform Only |
| domain | assume_role_policy | none | N/A | Terraform Only |
| model | assume_role_policy | none | N/A | Terraform Only |
| algorithm | attribute_names | none | N/A | Terraform Only |
| model_card | content | none | N/A | Terraform Only |
| endpoint_configuration | csv_content_types | none | N/A | Terraform Only |
| flow_definition | human_loop_activation_conditions | none | N/A | Terraform Only |
| data_quality_job_definition | json | none | N/A | Terraform Only |
| monitoring_schedule | json | none | N/A | Terraform Only |
| endpoint_configuration | json_content_types | none | N/A | Terraform Only |
| workforce | jwks_uri | none | N/A | Terraform Only |
| data_quality_job_definition | line | none | N/A | Terraform Only |
| monitoring_schedule | line | none | N/A | Terraform Only |
| algorithm | policy | none | N/A | Terraform Only |
| data_quality_job_definition | record_preprocessor_source_uri | none | N/A | Terraform Only |
| model_package_group_policy | resource_policy | none | N/A | Terraform Only |
| code_repository | secret_string | none | N/A | Terraform Only |

## secretsmanager

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| secret | policy | none | N/A | Terraform Only |
| secret_policy | policy | none | N/A | Terraform Only |
| secret_version | secret_string | none | N/A | Terraform Only |

## ses

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| identity_policy | policy | none | N/A | Terraform Only |

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

## ssm

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Document | ⚠️ Content | none | Yes | Gap (Confirmed) |
| activation | assume_role_policy | none | N/A | Terraform Only |
| maintenance_window_task | payload | none | N/A | Terraform Only |
| resource_data_sync | policy | none | N/A | Terraform Only |
| resource_data_sync | sync_format | none | N/A | Terraform Only |

## wafv2

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| IPSet | ⚠️ Scope | none | Yes | Gap (Confirmed) |
| RuleGroup | ⚠️ Scope | none | Yes | Gap (Confirmed) |
| WebACL | ⚠️ Scope | none | Yes | Gap (Confirmed) |
| web_acl_rule | all | none | N/A | Terraform Only |
| web_acl_association | body | none | N/A | Terraform Only |
| rule_group | content_type | none | N/A | Terraform Only |
| web_acl | content_type | none | N/A | Terraform Only |
| web_acl | default_size_inspection_limit | none | N/A | Terraform Only |
| rule_group | fallback_behavior | none | N/A | Terraform Only |
| web_acl | fallback_behavior | none | N/A | Terraform Only |
| web_acl_rule_group_association | identifier | none | N/A | Terraform Only |
| web_acl_rule | included_paths | none | N/A | Terraform Only |
| rule_group | invalid_fallback_behavior | none | N/A | Terraform Only |
| web_acl | invalid_fallback_behavior | none | N/A | Terraform Only |
| web_acl_rule | invalid_fallback_behavior | none | N/A | Terraform Only |
| web_acl_rule_group_association | json | none | N/A | Terraform Only |
| rule_group | json_body | none | N/A | Terraform Only |
| web_acl | json_body | none | N/A | Terraform Only |
| web_acl_rule | json_body | none | N/A | Terraform Only |
| rule_group | match_pattern | none | N/A | Terraform Only |
| web_acl | match_pattern | none | N/A | Terraform Only |
| web_acl_rule | match_pattern | none | N/A | Terraform Only |
| web_acl_rule_group_association | payload_type | none | N/A | Terraform Only |
| web_acl_logging_configuration | policy_document | none | N/A | Terraform Only |
| rule_group | rules_json | none | N/A | Terraform Only |

