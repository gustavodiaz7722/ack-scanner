# ACK Scanner Gap Analysis Report

## Summary

| Metric | Count |
| --- | --- |
| Total Fields | 187 |
| Already Annotated | 20 |
| Gaps (Terraform Confirmed) | 31 |
| Gaps (Unconfirmed) | 0 |
| Terraform Only | 136 |

## Priority Services

| Service | Confirmed Gaps |
| --- | --- |
| sagemaker | 10 |
| sns | 4 |
| dms | 3 |
| pipes | 3 |
| cloudwatch | 2 |
| lambda | 2 |
| apigatewayv2 | 1 |
| ecs | 1 |
| efs | 1 |
| kms | 1 |
| opensearchserverless | 1 |
| ram | 1 |
| ssm | 1 |

## acm

**Resources:** Certificate

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## acmpca

**Resources:** Certificate, CertificateAuthority, CertificateAuthorityActivation

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| policy | policy | none | N/A | Terraform Only |

## apigateway

**Resources:** APIIntegrationResponse, APIKey, APIMethodResponse, Authorizer, Deployment, Integration, Method, Resource, RestAPI, Stage, VPCLink

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## apigatewayv2

**Resources:** API, APIMapping, Authorizer, Deployment, DomainName, Integration, Route, Stage, VPCLink

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Authorizer | ⚠️ Issuer | none | Yes | Gap (Confirmed) |
| model | content_type | none | N/A | Terraform Only |
| model | schema | none | N/A | Terraform Only |

## applicationautoscaling

**Resources:** ScalableTarget, ScalingPolicy

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## athena

**Resources:** PreparedStatement, WorkGroup

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## autoscaling

**Resources:** AutoScalingGroup

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| group | notification_metadata | none | N/A | Terraform Only |
| lifecycle_hook | notification_metadata | none | N/A | Terraform Only |

## backup

**Resources:** BackupPlan, BackupSelection, BackupVault

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| report_plan | formats | none | N/A | Terraform Only |
| vault_policy | policy | none | N/A | Terraform Only |

## bedrock

**Resources:** InferenceProfile

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## bedrockagent

**Resources:** Agent

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| agent_action_group | api_schema | none | N/A | Terraform Only |
| flow | json | none | N/A | Terraform Only |
| prompt | json | none | N/A | Terraform Only |
| agent_action_group | payload | none | N/A | Terraform Only |

## bedrockagentcorecontrol

**Resources:** APIKeyCredentialProvider, AgentRuntime, AgentRuntimeEndpoint, Browser, BrowserProfile, CodeInterpreter, Gateway, GatewayTarget, Memory, WorkloadIdentity

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## cloudfront

**Resources:** CachePolicy, Distribution, Function, OriginAccessControl, OriginRequestPolicy, ResponseHeadersPolicy, VPCOrigin

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## cloudtrail

**Resources:** EventDataStore, Trail

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## cloudwatch

**Resources:** Dashboard, MetricAlarm, MetricStream

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

## cloudwatchlogs

**Resources:** LogGroup

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## codeartifact

**Resources:** Domain, PackageGroup

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| domain_permissions_policy | policy_document | none | N/A | Terraform Only |
| repository_permissions_policy | policy_document | none | N/A | Terraform Only |

## cognitoidentityprovider

**Resources:** UserPool

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## dms

**Resources:** Certificate, Endpoint, ReplicationSubnetGroup

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

## documentdb

**Resources:** DBCluster, DBInstance, DBSubnetGroup

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## dsql

**Resources:** Cluster

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Cluster | Policy | is_iam_policy | N/A | Annotated |

## dynamodb

**Resources:** Backup, GlobalTable, Table

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| table_item | item | none | N/A | Terraform Only |
| resource_policy | policy | none | N/A | Terraform Only |
| resource_policy | resource_arn | none | N/A | Terraform Only |

## ec2

**Resources:** CapacityReservation, DHCPOptions, EgressOnlyInternetGateway, ElasticIPAddress, FlowLog, Instance, InternetGateway, LaunchTemplate, ManagedPrefixList, NATGateway, NetworkACL, RouteTable, SecurityGroup, Subnet, TransitGateway, TransitGatewayVPCAttachment, VPC, VPCEndpoint, VPCEndpointServiceConfiguration, VPCPeeringConnection

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## ecr

**Resources:** PullThroughCacheRule, Repository, RepositoryCreationTemplate

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Repository | LifecyclePolicy | is_document | N/A | Annotated |
| RepositoryCreationTemplate | LifecyclePolicy | is_document | N/A | Annotated |
| Repository | Policy | is_iam_policy | N/A | Annotated |
| RepositoryCreationTemplate | RepositoryPolicy | is_iam_policy | N/A | Annotated |

## ecrpublic

**Resources:** Repository

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| repository_policy | policy | none | N/A | Terraform Only |

## ecs

**Resources:** CapacityProvider, Cluster, Service, TaskDefinition

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Service | ⚠️ LogDriver | none | Yes | Gap (Confirmed) |
| task_definition | container_definitions | none | N/A | Terraform Only |
| service | format | none | N/A | Terraform Only |

## efs

**Resources:** AccessPoint, FileSystem, MountTarget

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| FileSystem | ⚠️ Policy | none | Yes | Gap (Confirmed) |

## eks

**Resources:** AccessEntry, Addon, Capability, Cluster, FargateProfile, IdentityProviderConfig, Nodegroup, PodIdentityAssociation

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Addon | ConfigurationValues | is_document | N/A | Annotated |
| PodIdentityAssociation | Policy | is_iam_policy | N/A | Annotated |

## elasticache

**Resources:** CacheCluster, CacheParameterGroup, CacheSubnetGroup, ReplicationGroup, ServerlessCache, ServerlessCacheSnapshot, Snapshot, User, UserGroup

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## elbv2

**Resources:** Listener, LoadBalancer, Rule, TargetGroup

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## emrcontainers

**Resources:** JobRun, VirtualCluster

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## emrserverless

**Resources:** Application

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## eventbridge

**Resources:** Archive, Endpoint, EventBus, Rule

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## firehose

**Resources:** DeliveryStream

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## glue

**Resources:** Job

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

**Resources:** Group, InstanceProfile, OpenIDConnectProvider, Policy, Role, ServiceLinkedRole, User

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Role | AssumeRolePolicyDocument | is_iam_policy | N/A | Annotated |
| Policy | PolicyDocument | is_iam_policy | N/A | Annotated |

## kafka

**Resources:** Cluster, Configuration, ServerlessCluster

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## keyspaces

**Resources:** Keyspace, Table

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## kinesis

**Resources:** Stream

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

**Resources:** Alias, Grant, Key

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Key | ⚠️ Policy | none | Yes | Gap (Confirmed) |

## lambda

**Resources:** Alias, CodeSigningConfig, EventSourceMapping, Function, FunctionURLConfig, LayerVersion, Version

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| EventSourceMapping | ⚠️ EventRecordFormat | none | Yes | Gap (Confirmed) |
| Function | ⚠️ LogFormat | none | Yes | Gap (Confirmed) |
| invocation | input | none | N/A | Terraform Only |
| invocation | lifecycle_scope | none | N/A | Terraform Only |
| event_source_mapping | pattern | none | N/A | Terraform Only |
| invocation | terraform_key | none | N/A | Terraform Only |

## memorydb

**Resources:** ACL, Cluster, ParameterGroup, Snapshot, SubnetGroup, User

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## mq

**Resources:** Broker

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## mwaa

**Resources:** Environment

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## networkfirewall

**Resources:** Firewall, FirewallPolicy, RuleGroup

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| resource_policy | policy | none | N/A | Terraform Only |

## networkmanager

**Resources:** GlobalNetwork

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| core_network | create_base_policy | none | N/A | Terraform Only |
| core_network_policy_attachment | policy_document | none | N/A | Terraform Only |

## opensearchserverless

**Resources:** Collection, SecurityPolicy

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| SecurityPolicy | ⚠️ Policy | none | Yes | Gap (Confirmed) |

## opensearchservice

**Resources:** Domain

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Domain | AccessPolicies | is_iam_policy | N/A | Annotated |

## organizations

**Resources:** Account, OrganizationalUnit

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| policy | content | none | N/A | Terraform Only |
| resource_policy | content | none | N/A | Terraform Only |

## pipes

**Resources:** Pipe

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Pipe | ⚠️ InputTemplate | none | Yes | Gap (Confirmed) |
| Pipe | ⚠️ Time | none | Yes | Gap (Confirmed) |
| Pipe | ⚠️ Timestamp | none | Yes | Gap (Confirmed) |
| pipe | output_format | none | N/A | Terraform Only |
| pipe | pattern | none | N/A | Terraform Only |

## prometheusservice

**Resources:** AlertManagerDefinition, LoggingConfiguration, RuleGroupsNamespace, Workspace

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## quicksight

**Resources:** Analysis, Dashboard, DataSet, DataSource

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| data_set | format | none | N/A | Terraform Only |

## ram

**Resources:** Permission, ResourceShare

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Permission | ⚠️ PolicyTemplate | none | Yes | Gap (Confirmed) |

## rds

**Resources:** DBCluster, DBClusterEndpoint, DBClusterParameterGroup, DBClusterSnapshot, DBInstance, DBParameterGroup, DBProxy, DBSnapshot, DBSubnetGroup, GlobalCluster

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| custom_db_engine_version | manifest | none | N/A | Terraform Only |
| custom_db_engine_version | manifest_hash | none | N/A | Terraform Only |

## recyclebin

**Resources:** Rule

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## route53

**Resources:** HealthCheck, HostedZone, RecordSet

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| traffic_policy | document | none | N/A | Terraform Only |

## route53resolver

**Resources:** ResolverEndpoint, ResolverQueryLogConfig, ResolverRule, ResolverRuleAssociation

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## s3

**Resources:** Bucket

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Bucket | Policy | is_iam_policy | N/A | Annotated |
| object_copy | kms_encryption_context | none | N/A | Terraform Only |
| bucket | routing_rules | none | N/A | Terraform Only |
| object_copy | source | none | N/A | Terraform Only |

## s3control

**Resources:** AccessPoint

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| object_lambda_access_point | function_payload | none | N/A | Terraform Only |
| access_grants_instance_resource_policy | policy | none | N/A | Terraform Only |
| access_point_policy | policy | none | N/A | Terraform Only |
| bucket_policy | policy | none | N/A | Terraform Only |
| multi_region_access_point_policy | policy | none | N/A | Terraform Only |
| object_lambda_access_point_policy | policy | none | N/A | Terraform Only |

## s3files

**Resources:** AccessPoint, FileSystem, MountTarget

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| FileSystem | Policy | is_iam_policy | N/A | Annotated |

## s3tables

**Resources:** Namespace, TableBucket

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| table_bucket_policy | resource_policy | none | N/A | Terraform Only |
| table_policy | resource_policy | none | N/A | Terraform Only |

## sagemaker

**Resources:** App, DataQualityJobDefinition, Domain, Endpoint, EndpointConfig, FeatureGroup, HyperParameterTuningJob, InferenceComponent, LabelingJob, Model, ModelBiasJobDefinition, ModelExplainabilityJobDefinition, ModelPackage, ModelPackageGroup, ModelQualityJobDefinition, MonitoringSchedule, NotebookInstance, NotebookInstanceLifecycleConfig, Pipeline, PipelineExecution, ProcessingJob, Project, Space, TrainingJob, TransformJob, UserProfile

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

**Resources:** Secret

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| secret | policy | none | N/A | Terraform Only |
| secret_policy | policy | none | N/A | Terraform Only |
| secret_version | secret_string | none | N/A | Terraform Only |

## ses

**Resources:** ConfigurationSet

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| identity_policy | policy | none | N/A | Terraform Only |

## sfn

**Resources:** Activity, StateMachine, StateMachineAlias

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |

_No findings for this controller._

## sns

**Resources:** PlatformApplication, PlatformEndpoint, Subscription, Topic

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Topic | ⚠️ DeliveryPolicy | none | Yes | Gap (Confirmed) |
| Subscription | ⚠️ FilterPolicyScope | none | Yes | Gap (Confirmed) |
| Subscription | ⚠️ Protocol | none | Yes | Gap (Confirmed) |
| Subscription | ⚠️ RawMessageDelivery | none | Yes | Gap (Confirmed) |
| Subscription | DeliveryPolicy | is_document | N/A | Annotated |
| Subscription | FilterPolicy | is_document | N/A | Annotated |
| Topic | Policy | is_iam_policy | N/A | Annotated |
| Subscription | RedrivePolicy | is_document | N/A | Annotated |
| topic_subscription | replay_policy | none | N/A | Terraform Only |

## sqs

**Resources:** Queue

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Queue | Policy | is_iam_policy | N/A | Annotated |
| Queue | RedriveAllowPolicy | is_document | N/A | Annotated |
| Queue | RedrivePolicy | is_document | N/A | Annotated |

## ssm

**Resources:** Document, Parameter, PatchBaseline, ResourceDataSync

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| ResourceDataSync | ⚠️ SyncFormat | none | Yes | Gap (Confirmed) |
| Document | Content | is_document | N/A | Annotated |
| maintenance_window_task | payload | none | N/A | Terraform Only |

## wafv2

**Resources:** IPSet, RuleGroup, WebACL

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

