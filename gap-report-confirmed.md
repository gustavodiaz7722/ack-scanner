# ACK Scanner Gap Analysis Report

## Summary

| Metric | Count |
| --- | --- |
| Total Fields | 41 |
| Already Annotated | 0 |
| Gaps (Terraform Confirmed) | 41 |
| Gaps (Unconfirmed) | 0 |
| Terraform Only | 0 |

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

## apigatewayv2

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Authorizer | ⚠️ Issuer | none | Yes | Gap (Confirmed) |

## cloudwatch

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Dashboard | ⚠️ DashboardBody | none | Yes | Gap (Confirmed) |
| MetricStream | ⚠️ OutputFormat | none | Yes | Gap (Confirmed) |

## dms

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Endpoint | ⚠️ ExternalTableDefinition | none | Yes | Gap (Confirmed) |
| Endpoint | ⚠️ MessageFormat | none | Yes | Gap (Confirmed) |
| Endpoint | ⚠️ SpatialDataOptionToGeoJSONFunctionName | none | Yes | Gap (Confirmed) |

## ecr

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Repository | ⚠️ LifecyclePolicy | none | Yes | Gap (Confirmed) |
| RepositoryCreationTemplate | ⚠️ LifecyclePolicy | none | Yes | Gap (Confirmed) |
| Repository | ⚠️ Policy | none | Yes | Gap (Confirmed) |
| RepositoryCreationTemplate | ⚠️ RepositoryPolicy | none | Yes | Gap (Confirmed) |

## ecs

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Service | ⚠️ LogDriver | none | Yes | Gap (Confirmed) |

## efs

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| FileSystem | ⚠️ Policy | none | Yes | Gap (Confirmed) |

## iam

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Policy | ⚠️ PolicyDocument | none | Yes | Gap (Confirmed) |

## kms

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Key | ⚠️ Policy | none | Yes | Gap (Confirmed) |

## lambda

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| EventSourceMapping | ⚠️ EventRecordFormat | none | Yes | Gap (Confirmed) |
| Function | ⚠️ LogFormat | none | Yes | Gap (Confirmed) |

## opensearchserverless

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| SecurityPolicy | ⚠️ Policy | none | Yes | Gap (Confirmed) |

## pipes

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Pipe | ⚠️ InputTemplate | none | Yes | Gap (Confirmed) |
| Pipe | ⚠️ Time | none | Yes | Gap (Confirmed) |
| Pipe | ⚠️ Timestamp | none | Yes | Gap (Confirmed) |

## ram

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Permission | ⚠️ PolicyTemplate | none | Yes | Gap (Confirmed) |

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

## sns

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Subscription | ⚠️ DeliveryPolicy | none | Yes | Gap (Confirmed) |
| Topic | ⚠️ DeliveryPolicy | none | Yes | Gap (Confirmed) |
| Subscription | ⚠️ FilterPolicyScope | none | Yes | Gap (Confirmed) |
| Subscription | ⚠️ Protocol | none | Yes | Gap (Confirmed) |
| Subscription | ⚠️ RawMessageDelivery | none | Yes | Gap (Confirmed) |
| Subscription | ⚠️ RedrivePolicy | none | Yes | Gap (Confirmed) |

## sqs

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Queue | ⚠️ RedriveAllowPolicy | none | Yes | Gap (Confirmed) |
| Queue | ⚠️ RedrivePolicy | none | Yes | Gap (Confirmed) |

## ssm

| Resource | Field | Annotation | TF Confirmed | Category |
| --- | --- | --- | --- | --- |
| Document | ⚠️ Content | none | Yes | Gap (Confirmed) |
| ResourceDataSync | ⚠️ SyncFormat | none | Yes | Gap (Confirmed) |

