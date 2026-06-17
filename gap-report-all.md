# ACK Scanner Gap Analysis Report

## Summary

| Metric | Count |
| --- | --- |
| Total Fields | 187 |
| Already Annotated | 10 |
| Gaps (Terraform Confirmed) | 41 |
| Gaps (Unconfirmed) | 0 |
| Terraform Only | 136 |

## Priority Services

| Service | Confirmed Gaps |
| --- | --- |
| sagemaker | 10 |
| sns | 6 |
| ecr | 4 |
| dms | 3 |
| pipes | 3 |
| cloudwatch | 2 |
| lambda | 2 |
| sqs | 2 |
| ssm | 2 |
| apigatewayv2 | 1 |
| ecs | 1 |
| efs | 1 |
| iam | 1 |
| kms | 1 |
| opensearchserverless | 1 |
| ram | 1 |

## acmpca

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| policy | policy | none | N/A | Terraform Only |

## apigatewayv2

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Authorizer | ⚠️ Issuer | none | Yes | Gap (Confirmed) |
| model | content_type | none | N/A | Terraform Only |
| model | schema | none | N/A | Terraform Only |

## autoscaling

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| group | notification_metadata | none | N/A | Terraform Only |
| lifecycle_hook | notification_metadata | none | N/A | Terraform Only |

## backup

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| report_plan | formats | none | N/A | Terraform Only |
| vault_policy | policy | none | N/A | Terraform Only |

## bedrockagent

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| agent_action_group | api_schema | none | N/A | Terraform Only |
| flow | json | none | N/A | Terraform Only |
| prompt | json | none | N/A | Terraform Only |
| agent_action_group | payload | none | N/A | Terraform Only |

## cloudwatch

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Dashboard | ⚠️ DashboardBody | none | Yes | Gap (Confirmed) |
| MetricStream | ⚠️ OutputFormat | none | Yes | Gap (Confirmed) |
| log_destination_policy | access_policy | none | N/A | Terraform Only |
| log_delivery_destination_policy | delivery_destination_policy | none | N/A | Terraform Only |
| event_archive | event_pattern | none | N/A | Terraform Only |
| event_rule | event_pattern | none | N/A | Terraform Only |
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
| event_bus_policy | policy | none | N/A | Terraform Only |
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
| Endpoint | ⚠️ MessageFormat | none | Yes | Gap (Confirmed) |
| Endpoint | ⚠️ SpatialDataOptionToGeoJSONFunctionName | none | Yes | Gap (Confirmed) |
| replication_config | replication_settings | none | N/A | Terraform Only |
| replication_task | replication_task_settings | none | N/A | Terraform Only |
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
| table_item | item | none | N/A | Terraform Only |
| resource_policy | policy | none | N/A | Terraform Only |
| resource_policy | resource_arn | none | N/A | Terraform Only |

## ecr

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Repository | ⚠️ LifecyclePolicy | none | Yes | Gap (Confirmed) |
| RepositoryCreationTemplate | ⚠️ LifecyclePolicy | none | Yes | Gap (Confirmed) |
| Repository | ⚠️ Policy | none | Yes | Gap (Confirmed) |
| RepositoryCreationTemplate | ⚠️ RepositoryPolicy | none | Yes | Gap (Confirmed) |

## ecrpublic

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| repository_policy | policy | none | N/A | Terraform Only |

## ecs

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Service | ⚠️ LogDriver | none | Yes | Gap (Confirmed) |
| task_definition | container_definitions | none | N/A | Terraform Only |
| service | format | none | N/A | Terraform Only |

## efs

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| FileSystem | ⚠️ Policy | none | Yes | Gap (Confirmed) |

## eks

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Addon | ConfigurationValues | is_document | N/A | Annotated |
| PodIdentityAssociation | Policy | is_iam_policy | N/A | Annotated |

## glue

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| connection | SparkProperties | none | N/A | Terraform Only |
| classifier | classification | none | N/A | Terraform Only |
| crawler | configuration | none | N/A | Terraform Only |
| connection | connection_type | none | N/A | Terraform Only |
| schema | data_format | none | N/A | Terraform Only |
| catalog_table | initial_default | none | N/A | Terraform Only |
| classifier | json_classifier | none | N/A | Terraform Only |
| classifier | json_path | none | N/A | Terraform Only |
| catalog_table | type | none | N/A | Terraform Only |
| catalog_table | write_default | none | N/A | Terraform Only |

## iam

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Policy | ⚠️ PolicyDocument | none | Yes | Gap (Confirmed) |
| Role | AssumeRolePolicyDocument | is_iam_policy | N/A | Annotated |

## kinesis

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
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
| firehose_delivery_stream | parameter_name | none | N/A | Terraform Only |
| resource_policy | policy | none | N/A | Terraform Only |
| analytics_application | record_format_type | none | N/A | Terraform Only |
| firehose_delivery_stream | timestamp_formats | none | N/A | Terraform Only |

## kms

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Key | ⚠️ Policy | none | Yes | Gap (Confirmed) |

## lambda

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| EventSourceMapping | ⚠️ EventRecordFormat | none | Yes | Gap (Confirmed) |
| Function | ⚠️ LogFormat | none | Yes | Gap (Confirmed) |
| invocation | input | none | N/A | Terraform Only |
| invocation | lifecycle_scope | none | N/A | Terraform Only |
| event_source_mapping | pattern | none | N/A | Terraform Only |
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
| Pipe | ⚠️ InputTemplate | none | Yes | Gap (Confirmed) |
| Pipe | ⚠️ Time | none | Yes | Gap (Confirmed) |
| Pipe | ⚠️ Timestamp | none | Yes | Gap (Confirmed) |
| pipe | output_format | none | N/A | Terraform Only |
| pipe | pattern | none | N/A | Terraform Only |

## quicksight

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| data_set | format | none | N/A | Terraform Only |

## ram

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Permission | ⚠️ PolicyTemplate | none | Yes | Gap (Confirmed) |

## rds

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| custom_db_engine_version | manifest | none | N/A | Terraform Only |
| custom_db_engine_version | manifest_hash | none | N/A | Terraform Only |

## route53

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| traffic_policy | document | none | N/A | Terraform Only |

## s3

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Bucket | Policy | is_iam_policy | N/A | Annotated |
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
| ModelPackage | ⚠️ ContentDigest | none | Yes | Gap (Confirmed) |
| ModelPackage | ⚠️ ContentType | none | Yes | Gap (Confirmed) |
| TransformJob | ⚠️ ContentType | none | Yes | Gap (Confirmed) |
| Pipeline | ⚠️ PipelineDefinition | none | Yes | Gap (Confirmed) |
| DataQualityJobDefinition | ⚠️ PostAnalyticsProcessorSourceURI | none | Yes | Gap (Confirmed) |
| ModelQualityJobDefinition | ⚠️ PostAnalyticsProcessorSourceURI | none | Yes | Gap (Confirmed) |
| MonitoringSchedule | ⚠️ PostAnalyticsProcessorSourceURI | none | Yes | Gap (Confirmed) |
| DataQualityJobDefinition | ⚠️ RecordPreprocessorSourceURI | none | Yes | Gap (Confirmed) |
| ModelQualityJobDefinition | ⚠️ RecordPreprocessorSourceURI | none | Yes | Gap (Confirmed) |
| MonitoringSchedule | ⚠️ RecordPreprocessorSourceURI | none | Yes | Gap (Confirmed) |
| algorithm | attribute_names | none | N/A | Terraform Only |
| endpoint_configuration | csv_content_types | none | N/A | Terraform Only |
| flow_definition | human_loop_activation_conditions | none | N/A | Terraform Only |
| data_quality_job_definition | json | none | N/A | Terraform Only |
| monitoring_schedule | json | none | N/A | Terraform Only |
| endpoint_configuration | json_content_types | none | N/A | Terraform Only |
| workforce | jwks_uri | none | N/A | Terraform Only |
| data_quality_job_definition | line | none | N/A | Terraform Only |
| monitoring_schedule | line | none | N/A | Terraform Only |
| model_package_group_policy | resource_policy | none | N/A | Terraform Only |

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
| topic_subscription | replay_policy | none | N/A | Terraform Only |

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
| ResourceDataSync | ⚠️ SyncFormat | none | Yes | Gap (Confirmed) |
| maintenance_window_task | payload | none | N/A | Terraform Only |

## wafv2

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| web_acl_rule | all | none | N/A | Terraform Only |
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
| rule_group | match_scope | none | N/A | Terraform Only |
| web_acl | match_scope | none | N/A | Terraform Only |
| web_acl_rule | match_scope | none | N/A | Terraform Only |
| web_acl_rule_group_association | payload_type | none | N/A | Terraform Only |
| rule_group | rules_json | none | N/A | Terraform Only |

