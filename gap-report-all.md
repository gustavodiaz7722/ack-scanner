# ACK Scanner Gap Analysis Report

## Summary

| Metric | Count |
| --- | --- |
| Total Fields | 523 |
| Already Annotated | 10 |
| Gaps (Terraform Confirmed) | 29 |
| Gaps (Unconfirmed) | 0 |
| Terraform Only | 484 |

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

## amplify

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| branch | event_pattern | none | N/A | Terraform Only |
| branch | policy | none | N/A | Terraform Only |

## api

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| gateway_account | assume_role_policy | none | N/A | Terraform Only |
| gateway_authorizer | assume_role_policy | none | N/A | Terraform Only |
| gateway_integration | assume_role_policy | none | N/A | Terraform Only |
| gateway_deployment | body | none | N/A | Terraform Only |
| gateway_method_settings | body | none | N/A | Terraform Only |
| gateway_rest_api | body | none | N/A | Terraform Only |
| gateway_rest_api_put | body | none | N/A | Terraform Only |
| gateway_stage | body | none | N/A | Terraform Only |
| gateway_usage_plan | body | none | N/A | Terraform Only |
| gateway_account | policy | none | N/A | Terraform Only |
| gateway_authorizer | policy | none | N/A | Terraform Only |
| gateway_domain_name | policy | none | N/A | Terraform Only |
| gateway_rest_api | policy | none | N/A | Terraform Only |
| gateway_rest_api_policy | policy | none | N/A | Terraform Only |
| gateway_documentation_part | properties | none | N/A | Terraform Only |
| gateway_method_response | schema | none | N/A | Terraform Only |
| gateway_model | schema | none | N/A | Terraform Only |

## apigatewayv2

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| model | content_type | none | N/A | Terraform Only |
| authorizer | issuer | none | N/A | Terraform Only |
| model | schema | none | N/A | Terraform Only |

## appconfig

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| extension | assume_role_policy | none | N/A | Terraform Only |
| extension_association | assume_role_policy | none | N/A | Terraform Only |
| hosted_configuration_version | content | none | N/A | Terraform Only |
| configuration_profile | type | none | N/A | Terraform Only |

## appfabric

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| ingestion_destination | format | none | N/A | Terraform Only |

## appflow

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| connector_profile | assume_role_policy | none | N/A | Terraform Only |
| flow | file_type | none | N/A | Terraform Only |
| flow | policy | none | N/A | Terraform Only |
| flow | s3_input_file_type | none | N/A | Terraform Only |

## applicationinsights

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| application | query | none | N/A | Terraform Only |

## appmesh

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| virtual_gateway | json | none | N/A | Terraform Only |
| virtual_node | json | none | N/A | Terraform Only |
| virtual_gateway | key | none | N/A | Terraform Only |
| virtual_node | key | none | N/A | Terraform Only |
| virtual_gateway | value | none | N/A | Terraform Only |
| virtual_node | value | none | N/A | Terraform Only |

## appsync

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| datasource | assume_role_policy | none | N/A | Terraform Only |
| graphql_api | assume_role_policy | none | N/A | Terraform Only |
| type | format | none | N/A | Terraform Only |
| datasource | policy | none | N/A | Terraform Only |

## arcregionswitch

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| plan | assume_role_policy | none | N/A | Terraform Only |

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

## batch

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| compute_environment | assume_role_policy | none | N/A | Terraform Only |
| job_definition | assume_role_policy | none | N/A | Terraform Only |
| job_definition | container_properties | none | N/A | Terraform Only |
| job_definition | ecs_properties | none | N/A | Terraform Only |
| job_definition | node_properties | none | N/A | Terraform Only |

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

## bedrockagentcore

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| agent_runtime | assume_role_policy | none | N/A | Terraform Only |
| browser | assume_role_policy | none | N/A | Terraform Only |
| code_interpreter | assume_role_policy | none | N/A | Terraform Only |
| gateway | assume_role_policy | none | N/A | Terraform Only |
| gateway_target | assume_role_policy | none | N/A | Terraform Only |
| harness | assume_role_policy | none | N/A | Terraform Only |
| memory | assume_role_policy | none | N/A | Terraform Only |
| online_evaluation_config | assume_role_policy | none | N/A | Terraform Only |
| harness | input_schema | none | N/A | Terraform Only |
| gateway_target | items_json | none | N/A | Terraform Only |
| agent_runtime | policy | none | N/A | Terraform Only |
| harness | policy | none | N/A | Terraform Only |
| resource_policy | policy | none | N/A | Terraform Only |
| gateway_target | properties_json | none | N/A | Terraform Only |

## budgets

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| budget_action | assume_role_policy | none | N/A | Terraform Only |
| budget_action | policy | none | N/A | Terraform Only |

## ce

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| anomaly_monitor | monitor_specification | none | N/A | Terraform Only |
| anomaly_subscription | policy | none | N/A | Terraform Only |

## chime

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| voice_connector_streaming | assume_role_policy | none | N/A | Terraform Only |

## chimesdkmediapipelines

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| media_insights_pipeline_configuration | assume_role_policy | none | N/A | Terraform Only |

## cloudcontrolapi

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| resource | desired_state | none | N/A | Terraform Only |
| resource | schema | none | N/A | Terraform Only |

## cloudformation

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| stack_instances | assume_role_policy | none | N/A | Terraform Only |
| stack_set | assume_role_policy | none | N/A | Terraform Only |
| stack_set_instance | assume_role_policy | none | N/A | Terraform Only |
| stack_instances | policy | none | N/A | Terraform Only |
| stack_set | policy | none | N/A | Terraform Only |
| stack_set_instance | policy | none | N/A | Terraform Only |
| stack | template_body | none | N/A | Terraform Only |
| stack_set | template_body | none | N/A | Terraform Only |

## cloudfront

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| realtime_log_config | assume_role_policy | none | N/A | Terraform Only |
| distribution | policy | none | N/A | Terraform Only |
| origin_access_identity | policy | none | N/A | Terraform Only |
| realtime_log_config | policy | none | N/A | Terraform Only |

## cloudsearch

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| domain_service_access_policy | access_policy | none | N/A | Terraform Only |

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

## codebuild

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| project | assume_role_policy | none | N/A | Terraform Only |
| project | policy | none | N/A | Terraform Only |
| report_group | policy | none | N/A | Terraform Only |
| resource_policy | policy | none | N/A | Terraform Only |

## codecommit

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| approval_rule_template | content | none | N/A | Terraform Only |

## codedeploy

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| deployment_group | assume_role_policy | none | N/A | Terraform Only |

## codepipeline

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
|  | assume_role_policy | none | N/A | Terraform Only |
| webhook | json_path | none | N/A | Terraform Only |
|  | policy | none | N/A | Terraform Only |

## codestarnotifications

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| notification_rule | policy | none | N/A | Terraform Only |

## cognito

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| identity_pool_roles_attachment | assume_role_policy | none | N/A | Terraform Only |
| log_delivery_configuration | assume_role_policy | none | N/A | Terraform Only |
| managed_user_pool_client | assume_role_policy | none | N/A | Terraform Only |
| user_group | assume_role_policy | none | N/A | Terraform Only |
| user_pool_client | assume_role_policy | none | N/A | Terraform Only |
| identity_pool_roles_attachment | policy | none | N/A | Terraform Only |
| log_delivery_configuration | policy | none | N/A | Terraform Only |
| user_pool_client | policy | none | N/A | Terraform Only |
| managed_login_branding | settings | none | N/A | Terraform Only |

## comprehend

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| document_classifier | attribute_names | none | N/A | Terraform Only |
| entity_recognizer | attribute_names | none | N/A | Terraform Only |

## config

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| config_rule | assume_role_policy | none | N/A | Terraform Only |
| configuration_aggregator | assume_role_policy | none | N/A | Terraform Only |
| configuration_recorder | assume_role_policy | none | N/A | Terraform Only |
| configuration_recorder_status | assume_role_policy | none | N/A | Terraform Only |
| delivery_channel | assume_role_policy | none | N/A | Terraform Only |
| config_rule | input_parameters | none | N/A | Terraform Only |
| organization_custom_policy_rule | input_parameters | none | N/A | Terraform Only |
| organization_custom_rule | input_parameters | none | N/A | Terraform Only |
| organization_managed_rule | input_parameters | none | N/A | Terraform Only |
| config_rule | policy | none | N/A | Terraform Only |
| configuration_recorder_status | policy | none | N/A | Terraform Only |
| delivery_channel | policy | none | N/A | Terraform Only |

## connect

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| contact_flow | content | none | N/A | Terraform Only |
| contact_flow_module | content | none | N/A | Terraform Only |
| contact_flow | content_hash | none | N/A | Terraform Only |
| contact_flow_module | content_hash | none | N/A | Terraform Only |

## controltower

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| landing_zone | manifest_json | none | N/A | Terraform Only |
| control | value | none | N/A | Terraform Only |

## customerprofiles

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| domain | policy | none | N/A | Terraform Only |

## datazone

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| domain | assume_role_policy | none | N/A | Terraform Only |
| environment_profile | assume_role_policy | none | N/A | Terraform Only |
| form_type | assume_role_policy | none | N/A | Terraform Only |
| glossary | assume_role_policy | none | N/A | Terraform Only |
| glossary_term | assume_role_policy | none | N/A | Terraform Only |
| domain | policy | none | N/A | Terraform Only |
| environment_profile | policy | none | N/A | Terraform Only |
| form_type | policy | none | N/A | Terraform Only |
| glossary | policy | none | N/A | Terraform Only |
| glossary_term | policy | none | N/A | Terraform Only |

## dlm

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| lifecycle_policy | assume_role_policy | none | N/A | Terraform Only |
| lifecycle_policy | policy | none | N/A | Terraform Only |

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

## elasticsearch

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| domain | access_policies | none | N/A | Terraform Only |
| domain_policy | access_policies | none | N/A | Terraform Only |
| domain | policy_document | none | N/A | Terraform Only |

## emr

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| cluster | additional_info | none | N/A | Terraform Only |
| cluster | assume_role_policy | none | N/A | Terraform Only |
| cluster | autoscaling_policy | none | N/A | Terraform Only |
| instance_group | autoscaling_policy | none | N/A | Terraform Only |
| security_configuration | configuration | none | N/A | Terraform Only |
| cluster | configurations_json | none | N/A | Terraform Only |
| instance_group | configurations_json | none | N/A | Terraform Only |
| cluster | policy | none | N/A | Terraform Only |

## finspace

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| kx_user | assume_role_policy | none | N/A | Terraform Only |

## fis

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| experiment_template | assume_role_policy | none | N/A | Terraform Only |
| experiment_template | policy | none | N/A | Terraform Only |

## flow

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| log | assume_role_policy | none | N/A | Terraform Only |
| log | policy | none | N/A | Terraform Only |

## fms

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| policy | managed_service_data | none | N/A | Terraform Only |

## gamelift

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| game_server_group | assume_role_policy | none | N/A | Terraform Only |

## glacier

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| vault | access_policy | none | N/A | Terraform Only |
| vault_lock | policy | none | N/A | Terraform Only |

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

## grafana

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| license_association | assume_role_policy | none | N/A | Terraform Only |
| role_association | assume_role_policy | none | N/A | Terraform Only |
| workspace | assume_role_policy | none | N/A | Terraform Only |
| workspace_saml_configuration | assume_role_policy | none | N/A | Terraform Only |
| workspace | configuration | none | N/A | Terraform Only |

## guardduty

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| publishing_destination | policy | none | N/A | Terraform Only |

## iam

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Policy | ⚠️ PolicyDocument | none | Yes | Gap (Confirmed) |
| Role | AssumeRolePolicyDocument | is_iam_policy | N/A | Annotated |

## imagebuilder

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| lifecycle_policy | assume_role_policy | none | N/A | Terraform Only |

## iot

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| provisioning_template | assume_role_policy | none | N/A | Terraform Only |
| role_alias | assume_role_policy | none | N/A | Terraform Only |
| topic_rule | assume_role_policy | none | N/A | Terraform Only |
| topic_rule | batch_mode | none | N/A | Terraform Only |
| topic_rule | message_format | none | N/A | Terraform Only |
| policy | policy | none | N/A | Terraform Only |
| policy_attachment | policy | none | N/A | Terraform Only |
| provisioning_template | policy | none | N/A | Terraform Only |
| topic_rule | policy | none | N/A | Terraform Only |
| provisioning_template | template_body | none | N/A | Terraform Only |

## ivschat

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| logging_configuration | assume_role_policy | none | N/A | Terraform Only |

## kendra

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| index | json_token_type_configuration | none | N/A | Terraform Only |
| data_source | template | none | N/A | Terraform Only |

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

## kinesisanalyticsv2

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| application | json_mapping_parameters | none | N/A | Terraform Only |
| application | mapping_parameters | none | N/A | Terraform Only |
| application | record_format_type | none | N/A | Terraform Only |

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

## lb

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| listener | content_type | none | N/A | Terraform Only |
| listener_rule | content_type | none | N/A | Terraform Only |
| listener | jwks_endpoint | none | N/A | Terraform Only |
| listener_rule | jwks_endpoint | none | N/A | Terraform Only |
| listener | routing_http_response_content_security_policy_header_value | none | N/A | Terraform Only |

## lexv2models

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| bot | assume_role_policy | none | N/A | Terraform Only |
| intent | assume_role_policy | none | N/A | Terraform Only |

## lightsail

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| container_service | policy | none | N/A | Terraform Only |

## m2

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| application | content | none | N/A | Terraform Only |
| application | definition | none | N/A | Terraform Only |

## media

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| store_container_policy | policy | none | N/A | Terraform Only |

## msk

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| cluster | assume_role_policy | none | N/A | Terraform Only |
| topic | configs | none | N/A | Terraform Only |
| cluster_policy | policy | none | N/A | Terraform Only |
| scram_secret_association | policy | none | N/A | Terraform Only |
| scram_secret_association | secret_string | none | N/A | Terraform Only |

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

## notifications

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| event_rule | event_pattern | none | N/A | Terraform Only |

## oam

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| sink_policy | policy | none | N/A | Terraform Only |

## observabilityadmin

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| s3_table_integration | assume_role_policy | none | N/A | Terraform Only |
| telemetry_pipeline | assume_role_policy | none | N/A | Terraform Only |
| telemetry_rule | output_format | none | N/A | Terraform Only |
| telemetry_rule_for_organization | output_format | none | N/A | Terraform Only |
| s3_table_integration | policy | none | N/A | Terraform Only |
| telemetry_pipeline | policy | none | N/A | Terraform Only |

## opensearch

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| domain | access_policies | none | N/A | Terraform Only |
| domain_policy | access_policies | none | N/A | Terraform Only |
| application | assume_role_policy | none | N/A | Terraform Only |
| domain | jwks_url | none | N/A | Terraform Only |
| application | policy | none | N/A | Terraform Only |
| domain | policy_document | none | N/A | Terraform Only |

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

## osis

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| pipeline | assume_role_policy | none | N/A | Terraform Only |

## pinpoint

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| gcm_channel | api_key | none | N/A | Terraform Only |
| email_channel | assume_role_policy | none | N/A | Terraform Only |
| event_stream | assume_role_policy | none | N/A | Terraform Only |
| email_template | default_substitutions | none | N/A | Terraform Only |
| email_channel | policy | none | N/A | Terraform Only |
| event_stream | policy | none | N/A | Terraform Only |
| gcm_channel | service_json | none | N/A | Terraform Only |

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

## prometheus

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| resource_policy | policy_document | none | N/A | Terraform Only |

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

## redshift

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| scheduled_action | assume_role_policy | none | N/A | Terraform Only |
| authentication_profile | authentication_profile_content | none | N/A | Terraform Only |
| cluster_snapshot | cluster_snapshot_content | none | N/A | Terraform Only |
| integration | policy | none | N/A | Terraform Only |
| resource_policy | policy | none | N/A | Terraform Only |
| scheduled_action | policy | none | N/A | Terraform Only |

## redshiftserverless

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| resource_policy | policy | none | N/A | Terraform Only |

## rekognition

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| stream_processor | assume_role_policy | none | N/A | Terraform Only |
| stream_processor | policy | none | N/A | Terraform Only |

## resourcegroups

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| group | query | none | N/A | Terraform Only |

## rolesanywhere

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| profile | assume_role_policy | none | N/A | Terraform Only |

## route53

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| traffic_policy | document | none | N/A | Terraform Only |
| hosted_zone_dnssec | policy | none | N/A | Terraform Only |
| key_signing_key | policy | none | N/A | Terraform Only |

## route53domains

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| delegation_signer_record | policy | none | N/A | Terraform Only |

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

## s3vectors

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| vector_bucket_policy | policy | none | N/A | Terraform Only |

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

## scheduler

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| schedule | input | none | N/A | Terraform Only |

## schemas

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| schema | content | none | N/A | Terraform Only |
| registry_policy | policy | none | N/A | Terraform Only |
| schema | type | none | N/A | Terraform Only |

## secretsmanager

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| secret | policy | none | N/A | Terraform Only |
| secret_policy | policy | none | N/A | Terraform Only |
| secret_version | secret_string | none | N/A | Terraform Only |

## securityhub

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| automation_rule_v2 | ocsf_finding_criteria_json | none | N/A | Terraform Only |

## servicecatalog

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| constraint | parameters | none | N/A | Terraform Only |
| service_action | parameters | none | N/A | Terraform Only |

## servicecatalogappregistry

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| attribute_group | attributes | none | N/A | Terraform Only |
| attribute_group_association | attributes | none | N/A | Terraform Only |

## ses

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| identity_policy | policy | none | N/A | Terraform Only |

## sesv2

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| email_identity_policy | policy | none | N/A | Terraform Only |

## shield

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| drt_access_role_arn_association | assume_role_policy | none | N/A | Terraform Only |
| proactive_engagement | assume_role_policy | none | N/A | Terraform Only |

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

## ssmquicksetup

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| configuration_manager | selected_patch_baselines | none | N/A | Terraform Only |

## ssoadmin

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| permission_set_inline_policy | inline_policy | none | N/A | Terraform Only |
| trusted_token_issuer | jwks_retrieval_option | none | N/A | Terraform Only |
| trusted_token_issuer | oidc_jwt_configuration | none | N/A | Terraform Only |
| customer_managed_policy_attachment | policy | none | N/A | Terraform Only |
| customer_managed_policy_attachments_exclusive | policy | none | N/A | Terraform Only |
| permissions_boundary_attachment | policy | none | N/A | Terraform Only |

## timestreaminfluxdb

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| db_cluster | policy | none | N/A | Terraform Only |
| db_instance | policy | none | N/A | Terraform Only |

## timestreamquery

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| scheduled_query | assume_role_policy | none | N/A | Terraform Only |
| scheduled_query | policy | none | N/A | Terraform Only |

## transcribe

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| language_model | assume_role_policy | none | N/A | Terraform Only |
| language_model | policy | none | N/A | Terraform Only |

## transfer

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| server | assume_role_policy | none | N/A | Terraform Only |
| ssh_key | assume_role_policy | none | N/A | Terraform Only |
| user | assume_role_policy | none | N/A | Terraform Only |
| web_app | assume_role_policy | none | N/A | Terraform Only |
| access | policy | none | N/A | Terraform Only |
| ssh_key | policy | none | N/A | Terraform Only |
| user | policy | none | N/A | Terraform Only |
| web_app | policy | none | N/A | Terraform Only |

## verifiedaccess

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| endpoint | policy_document | none | N/A | Terraform Only |
| group | policy_document | none | N/A | Terraform Only |

## verifiedpermissions

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| schema | value | none | N/A | Terraform Only |

## vpc

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| endpoint | policy | none | N/A | Terraform Only |
| endpoint_connection_notification | policy | none | N/A | Terraform Only |
| endpoint_policy | policy | none | N/A | Terraform Only |

## vpclattice

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| auth_policy | policy | none | N/A | Terraform Only |
| resource_policy | policy | none | N/A | Terraform Only |

## vpn

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| connection | bgp_log_output_format | none | N/A | Terraform Only |
| connection | log_output_format | none | N/A | Terraform Only |

## wafregional

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| web_acl_association | body | none | N/A | Terraform Only |

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

## workspaces

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| directory | assume_role_policy | none | N/A | Terraform Only |

## workspacesweb

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| browser_settings | browser_policy | none | N/A | Terraform Only |
| browser_settings_association | browser_policy | none | N/A | Terraform Only |
| session_logger | log_file_format | none | N/A | Terraform Only |
| session_logger | policy | none | N/A | Terraform Only |
| session_logger_association | policy | none | N/A | Terraform Only |

## xray

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| encryption_config | policy | none | N/A | Terraform Only |
| resource_policy | policy_document | none | N/A | Terraform Only |

