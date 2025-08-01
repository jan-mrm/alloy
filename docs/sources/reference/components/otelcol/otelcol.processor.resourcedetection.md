---
canonical: https://grafana.com/docs/alloy/latest/reference/components/otelcol/otelcol.processor.resourcedetection/
aliases:
  - ../otelcol.processor.resourcedetection/ # /docs/alloy/latest/reference/otelcol.processor.resourcedetection/
title: otelcol.processor.resourcedetection
labels:
  stage: general-availability
  products:
    - oss
description: Learn about otelcol.processor.resourcedetection
---

# `otelcol.processor.resourcedetection`

`otelcol.processor.resourcedetection` detects resource information from the host in a format that conforms to the [OpenTelemetry resource semantic conventions][], and appends or overrides the resource values in the telemetry data with this information.

[OpenTelemetry resource semantic conventions]: https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/resource/semantic_conventions/

{{< admonition type="note" >}}
`otelcol.processor.resourcedetection` is a wrapper over the upstream OpenTelemetry Collector Contrib [`resourcedetection`][] processor.
If necessary, bug reports or feature requests are redirected to the upstream repository.

[`resourcedetection`]: https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/{{< param "OTEL_VERSION" >}}/processor/resourcedetectionprocessor
{{< /admonition >}}

You can specify multiple `otelcol.processor.resourcedetection` components by giving them different labels.

## Usage

```alloy
otelcol.processor.resourcedetection "<LABEL>" {
  output {
    logs    = [...]
    metrics = [...]
    traces  = [...]
  }
}
```

## Arguments

You can use the following arguments with `otelcol.processor.resourcedetection`:

| Name        | Type           | Description                                                                        | Default   | Required |
|-------------|----------------|------------------------------------------------------------------------------------|-----------|----------|
| `detectors` | `list(string)` | An ordered list of named detectors used to detect resource information.            | `["env"]` | no       |
| `override`  | `bool`         | Configures whether existing resource attributes should be overridden or preserved. | `true`    | no       |
| `timeout`   | `duration`     | Timeout by which all specified detectors must complete.                            | `"5s"`    | no       |

`detectors` could contain the following values:

{{< column-list >}}

* `aks`
* `azure`
* `consul`
* `docker`
* `dynatrace`
* `ec2`
* `ecs`
* `eks`
* `elasticbeanstalk`
* `env`
* `gcp`
* `heroku`
* `kubeadm`
* `kubernetes_node`
* `lambda`
* `openshift`
* `system`

{{< /column-list >}}

`env` is the only detector that's not configured through a block.
The `env` detector reads resource information from the `OTEL_RESOURCE_ATTRIBUTES` environment variable.
This variable must be in the format `<key1>=<value1>,<key2>=<value2>,...`, the details of which are currently pending confirmation in the OpenTelemetry specification.

If a detector other than `env` is needed, you can customize it with the relevant block.
For example, you can customize the `ec2` detector with the [ec2][] block.
If you omit the [ec2][] block, the defaults specified in the [ec2][] block documentation are used.

If multiple detectors are inserting the same attribute name, the first detector to insert wins.
For example, if you had `detectors = ["eks", "ec2"]` then `cloud.platform` will be `aws_eks` instead of `ec2`.

The following order is recommended for AWS:

  1. [`lambda`][lambda]
  1. [`elasticbeanstalk`][elasticbeanstalk]
  1. [`eks`][eks]
  1. [`ecs`][ecs]
  1. [`ec2`][ec2]

## Blocks

You can use the following blocks with `otelcol.processor.resourcedetection`:

| Block                                  | Description                                                                                             | Required |
|----------------------------------------|---------------------------------------------------------------------------------------------------------|----------|
| [`output`][output]                     | Configures where to send received telemetry data.                                                       | yes      |
| [`aks`][aks]                           | Adds resource attributes related to Azure AKS.                                                          | no       |
| [`azure`][azure]                       | Queries the Azure Instance Metadata Service to retrieve various resource attributes.                    | no       |
| [`consul`][consul]                     | Queries a Consul agent and reads its configuration endpoint to retrieve values for resource attributes. | no       |
| [`debug_metrics`][debug_metrics]       | Configures the metrics that this component generates to monitor its state.                              | no       |
| [`docker`][docker]                     | Queries the Docker daemon to retrieve various resource attributes from the host machine.                | no       |
| [`dynatrace`][dynatrace]               | Loads resource information from the `dt_host_metadata.properties` file.                                 | no       |
| [`ec2`][ec2]                           | Reads resource information from the EC2 instance metadata API.                                          | no       |
| [`ecs`][ecs]                           | Queries the Task Metadata Endpoint to record information about the current ECS Task.                    | no       |
| [`eks`][eks]                           | Adds resource attributes for Amazon EKS.                                                                | no       |
| [`elasticbeanstalk`][elasticbeanstalk] | Reads the AWS X-Ray configuration file available on all Beanstalk instances with X-Ray Enabled.         | no       |
| [`gcp`][gcp]                           | Detects resource attributes using the Google Cloud Client Libraries for Go.                             | no       |
| [`heroku`][heroku]                     | Adds resource attributes derived from Heroku dyno metadata.                                             | no       |
| [`kubeadm`][kubeadm]                   | Queries the Kubernetes API server to retrieve kubeadm resource attributes.                              | no       |
| [`kubernetes_node`][kubernetes_node]   | Queries the Kubernetes API server to retrieve various node resource attributes.                         | no       |
| [`lambda`][lambda]                     | Uses the AWS Lambda runtime environment variables to retrieve various resource attributes.              | no       |
| [`openshift`][openshift]               | Queries the OpenShift and Kubernetes APIs to retrieve various resource attributes.                      | no       |
| [`system`][system]                     | Queries the host machine to retrieve various resource attributes.                                       | no       |

[output]: #output
[debug_metrics]: #debug_metrics
[ec2]: #ec2
[ecs]: #ecs
[eks]: #eks
[elasticbeanstalk]: #elasticbeanstalk
[lambda]: #lambda
[azure]: #azure
[aks]: #aks
[consul]: #consul
[docker]: #docker
[gcp]: #gcp
[heroku]: #heroku
[system]: #system
[openshift]: #openshift
[kubernetes_node]: #kubernetes_node
[kubeadm]: #kubeadm
[res-attr-cfg]: #resource-attribute-configuration
[dynatrace]: #dynatrace

### `output`

{{< badge text="Required" >}}

{{< docs/shared lookup="reference/components/output-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `aks`

The `aks` block adds resource attributes related to Azure AKS.

The `aks` block supports the following blocks:

| Block                                              | Description                                  | Required |
|----------------------------------------------------|----------------------------------------------|----------|
| [`resource_attributes`](#aks--resource_attributes) | Configures which resource attributes to add. | no       |

#### `aks` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                              | Description                                                                              | Required |
|------------------------------------|------------------------------------------------------------------------------------------|----------|
| [`cloud.platform`][res-attr-cfg]   | Toggles the `cloud.platform` resource attribute. Sets `enabled` to `true` by default.    | no       |
| [`cloud.provider`][res-attr-cfg]   | Toggles the `cloud.provider` resource attribute. Sets `enabled` to `true` by default.    | no       |
| [`k8s.cluster.name`][res-attr-cfg] | Toggles the `k8s.cluster.name` resource attribute. Sets `enabled` to `false` by default. | no       |

Example values:

* `cloud.platform`: `"azure_vm"`
* `cloud.provider`: `"azure"`

Azure AKS cluster name is derived from the Azure Instance Metadata Service's (IMDS) infrastructure resource group field.
This field contains the resource group and name of the cluster, separated by underscores. For example: `MC_<resource group>_<cluster name>_<location>`.

Example:

* Resource group: `my-resource-group`
* Cluster name: `my-cluster`
* Location: `eastus`
* Generated name: `MC_my-resource-group_my-cluster_eastus`

The cluster name is detected if it doesn't contain underscores and if a custom infrastructure resource group name wasn't used.

If accurate parsing can't be performed, the infrastructure resource group value is returned.
This value can be used to uniquely identify the cluster, because Azure won't allow users to create multiple clusters with the same infrastructure resource group name.

### `azure`

The `azure` block queries the [Azure Instance Metadata Service][] to retrieve various resource attributes.

[Azure Instance Metadata Service]: https://aka.ms/azureimds

The `azure` block supports the following blocks:

| Block                                                | Description                                  | Required |
|------------------------------------------------------|----------------------------------------------|----------|
| [`resource_attributes`](#azure--resource_attributes) | Configures which resource attributes to add. | no       |

The `azure` block supports the following attributes:

| Attribute | Type           | Description                                                                                     | Default | Required |
|-----------|----------------|-------------------------------------------------------------------------------------------------|---------|----------|
| `tags`    | `list(string)` | A list of regular expressions to match tag keys to add as resource attributes can be specified. | `[]`    | no       |

#### `azure` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                                      | Description                                                                                     | Required |
|--------------------------------------------|-------------------------------------------------------------------------------------------------|----------|
| [`azure.resourcegroup.name`][res-attr-cfg] | Toggles the `azure.resourcegroup.name` resource attribute. Sets `enabled` to `true` by default. | no       |
| [`azure.vm.name`][res-attr-cfg]            | Toggles the `azure.vm.name` resource attribute. Sets `enabled` to `true` by default.            | no       |
| [`azure.vm.scaleset.name`][res-attr-cfg]   | Toggles the `azure.vm.scaleset.name` resource attribute. Sets `enabled` to `true` by default.   | no       |
| [`azure.vm.size`][res-attr-cfg]            | Toggles the `azure.vm.size` resource attribute. Sets `enabled` to `true` by default.            | no       |
| [`cloud.account.id`][res-attr-cfg]         | Toggles the `cloud.account.id` resource attribute. Sets `enabled` to `true` by default.         | no       |
| [`cloud.platform`][res-attr-cfg]           | Toggles the `cloud.platform` resource attribute. Sets `enabled` to `true` by default.           | no       |
| [`cloud.provider`][res-attr-cfg]           | Toggles the `cloud.provider` resource attribute. Sets `enabled` to `true` by default.           | no       |
| [`cloud.region`][res-attr-cfg]             | Toggles the `cloud.region` resource attribute. Sets `enabled` to `true` by default.             | no       |
| [`host.id`][res-attr-cfg]                  | Toggles the `host.id` resource attribute. Sets `enabled` to `true` by default.                  | no       |
| [`host.name`][res-attr-cfg]                | Toggles the `host.name` resource attribute. Sets `enabled` to `true` by default.                | no       |

Example values:

* `cloud.platform`: `"azure_vm"`
* `cloud.provider`: `"azure"`

### `consul`

The `consul` block queries a Consul agent and reads its configuration endpoint to retrieve values for resource attributes.

The `consul` block supports the following attributes:

| Attribute    | Type           | Description                                                                       | Default | Required |
|--------------|----------------|-----------------------------------------------------------------------------------|---------|----------|
| `address`    | `string`       | The address of the Consul server                                                  | `""`    | no       |
| `datacenter` | `string`       | Data center to use. If not provided, the default agent data center is used.       | `""`    | no       |
| `meta`       | `list(string)` | Allowlist of [Consul Metadata][] keys to use as resource attributes.              | `[]`    | no       |
| `namespace`  | `string`       | The name of the namespace to send along for the request.                          | `""`    | no       |
| `token`      | `secret`       | A per-request ACL token which overrides the Consul agent's default (empty) token. | `""`    | no       |

`token` is only required if the [Consul ACL System][] is enabled.

[Consul Metadata]: https://www.consul.io/docs/agent/options#node_meta
[Consul ACL System]: https://www.consul.io/docs/security/acl/acl-system

The `consul` block supports the following blocks:

| Block                                                 | Description                                  | Required |
|-------------------------------------------------------|----------------------------------------------|----------|
| [`resource_attributes`](#consul--resource_attributes) | Configures which resource attributes to add. | no       |

#### `consul` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                          | Description                                                                         | Required |
|--------------------------------|-------------------------------------------------------------------------------------|----------|
| [`cloud.region`][res-attr-cfg] | Toggles the `cloud.region` resource attribute. Sets `enabled` to `true` by default. | no       |
| [`host.id`][res-attr-cfg]      | Toggles the `host.id` resource attribute. Sets `enabled` to `true` by default.      | no       |
| [`host.name`][res-attr-cfg]    | Toggles the `host.name` resource attribute. Sets `enabled` to `true` by default.    | no       |

### `debug_metrics`

{{< docs/shared lookup="reference/components/otelcol-debug-metrics-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `docker`

The `docker` block queries the Docker daemon to retrieve various resource attributes from the host machine.

You need to mount the Docker socket (`/var/run/docker.sock` on Linux) to contact the Docker daemon.
Docker detection doesn't work on MacOS.

The `docker` block supports the following blocks:

| Block                                                 | Description                                  | Required |
|-------------------------------------------------------|----------------------------------------------|----------|
| [`resource_attributes`](#docker--resource_attributes) | Configures which resource attributes to add. | no       |

#### `docker` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                       | Description                                                                      | Required |
|-----------------------------|----------------------------------------------------------------------------------|----------|
| [`host.name`][res-attr-cfg] | Toggles the `host.name` resource attribute. Sets `enabled` to `true` by default. | no       |
| [`os.type`][res-attr-cfg]   | Toggles the `os.type` resource attribute. Sets `enabled` to `true` by default.   | no       |

### `dynatrace`

The `dynatrace` block loads resource information from the `dt_host_metadata.properties` file which is located in the `/var/lib/dynatrace/enrichment` (on Unix systems) or `%ProgramData%\dynatrace\enrichment` (on Windows) directories.

The `dynatrace` block supports the following blocks:

| Block                                                    | Description                                  | Required |
|----------------------------------------------------------|----------------------------------------------|----------|
| [`resource_attributes`](#dynatrace--resource_attributes) | Configures which resource attributes to add. | no       |

#### `dynatrace` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                            | Description                                                                           | Required |
|----------------------------------|---------------------------------------------------------------------------------------|----------|
| [`host.name`][res-attr-cfg]      | Toggles the `host.name` resource attribute. Sets `enabled` to `true` by default.      | no       |
| [`dt.entity.host`][res-attr-cfg] | Toggles the `dt.entity.host` resource attribute. Sets `enabled` to `true` by default. | no       |

### `ec2`

The `ec2` block reads resource information from the [EC2 instance metadata API][] using the [AWS SDK for Go][].

The `ec2` block supports the following attributes:

| Attribute      | Type           | Description                                                                 | Default | Required |
|----------------|----------------|-----------------------------------------------------------------------------|---------|----------|
| `max_attempts` | `int`          | The maximum number of attempts to retrieve metadata.                        | `3`     | no       |
| `max_backoff`  | `duration`     | The maximum backoff time between retries.                                   | `"20s"` | no       |
| `tags`         | `list(string)` | A list of regular expressions to match against tag keys of an EC2 instance. | `[]`    | no       |

<!-- The following commented behavior is implemented in https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/37453 but it does not appear 
that the code actually will cause the collector to fail to start, waiting for a response from the author 
<!-- | `fail_on_missing_metadata` | `bool`         | Whether to fail if metadata is missing.                                     | `false` | no       |

By default the ec2 detector will log errors if the metadata endpoint is unavailable, but if `fail_on_missing_metadata` is `true` it will propagate that error instead which will cause {{< param "PRODUCT_NAME" >}} to fail to start. -->

If you are using a proxy server on your EC2 instance, it's important that you exempt requests for instance metadata as described in the [AWS cli user guide][].
Failing to do so can result in proxied or missing instance data.

If the instance is part of AWS ParallelCluster and the detector is failing to connect to the metadata server,
check the iptable and make sure the chain `PARALLELCLUSTER_IMDS` contains a rule that allows the {{< param "PRODUCT_NAME" >}} user to access `169.254.169.254/32`.

[AWS SDK for Go]: https://docs.aws.amazon.com/sdk-for-go/api/aws/ec2metadata/
[EC2 instance metadata API]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html
[AWS cli user guide]: https://github.com/awsdocs/aws-cli-user-guide/blob/a2393582590b64bd2a1d9978af15b350e1f9eb8e/doc_source/cli-configure-proxy.md#using-a-proxy-on-amazon-ec2-instances

`tags` can be used to gather tags for the EC2 instance which {{< param "PRODUCT_NAME" >}} is running on.
To fetch EC2 tags, the IAM role assigned to the EC2 instance must have a policy that includes the `ec2:DescribeTags` permission.

The `ec2` block supports the following blocks:

| Block                           | Description                                  | Required |
|---------------------------------|----------------------------------------------|----------|
| [``](#ec2--resource_attributes) | Configures which resource attributes to add. | no       |

#### `ec2` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                                     | Description                                                                                    | Required |
|-------------------------------------------|------------------------------------------------------------------------------------------------|----------|
| [`cloud.account.id`][res-attr-cfg]        | Toggles the `cloud.account.id` resource attribute. Sets `enabled` to `true` by default.        | no       |
| [`cloud.availability_zone`][res-attr-cfg] | Toggles the `cloud.availability_zone` resource attribute. Sets `enabled` to `true` by default. | no       |
| [`cloud.platform`][res-attr-cfg]          | Toggles the `cloud.platform` resource attribute. Sets `enabled` to `true` by default.          | no       |
| [`cloud.provider`][res-attr-cfg]          | Toggles the `cloud.provider` resource attribute. Sets `enabled` to `true` by default.          | no       |
| [`cloud.region`][res-attr-cfg]            | Toggles the `cloud.region` resource attribute. Sets `enabled` to `true` by default.            | no       |
| [`host.id`][res-attr-cfg]                 | Toggles the `host.id` resource attribute. Sets `enabled` to `true` by default.                 | no       |
| [`host.image.id`][res-attr-cfg]           | Toggles the `host.image.id` resource attribute. Sets `enabled` to `true` by default.           | no       |
| [`host.name`][res-attr-cfg]               | Toggles the `host.name` resource attribute. Sets `enabled` to `true` by default.               | no       |
| [`host.type`][res-attr-cfg]               | Toggles the `host.type` resource attribute. Sets `enabled` to `true` by default.               | no       |

### `ecs`

The `ecs` block queries the [Task Metadata Endpoint][] (TMDE) to record information about the current ECS Task.
Only TMDE V4 and V3 are supported.

[Task Metadata Endpoint]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint.html

The `ecs` block supports the following blocks:

| Block                                             | Description                                  | Required |
| ------------------------------------------------- | -------------------------------------------- | -------- |
| [resource_attr`ibutes](#ecs--resource_attributes) | Configures which resource attributes to add. | no       |

#### `ecs` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                                     | Description                                                                                    | Required |
|-------------------------------------------|------------------------------------------------------------------------------------------------|----------|
| [`aws.ecs.cluster.arn`][res-attr-cfg]     | Toggles the `aws.ecs.cluster.arn` resource attribute. Sets `enabled` to `true` by default.     | no       |
| [`aws.ecs.launchtype`][res-attr-cfg]      | Toggles the `aws.ecs.launchtype` resource attribute. Sets `enabled` to `true` by default.      | no       |
| [`aws.ecs.task.arn`][res-attr-cfg]        | Toggles the `aws.ecs.task.arn` resource attribute. Sets `enabled` to `true` by default.        | no       |
| [`aws.ecs.task.family`][res-attr-cfg]     | Toggles the `aws.ecs.task.family` resource attribute. Sets `enabled` to `true` by default.     | no       |
| [`aws.ecs.task.id`][res-attr-cfg]         | Toggles the `aws.ecs.task.id` resource attribute. Sets `enabled` to `true` by default.         | no       |
| [`aws.ecs.task.revision`][res-attr-cfg]   | Toggles the `aws.ecs.task.revision` resource attribute. Sets `enabled` to `true` by default.   | no       |
| [`aws.log.group.arns`][res-attr-cfg]      | Toggles the `aws.log.group.arns` resource attribute. Sets `enabled` to `true` by default.      | no       |
| [`aws.log.group.names`][res-attr-cfg]     | Toggles the `aws.log.group.names` resource attribute. Sets `enabled` to `true` by default.     | no       |
| [`aws.log.stream.arns`][res-attr-cfg]     | Toggles the `aws.log.stream.arns` resource attribute. Sets `enabled` to `true` by default.     | no       |
| [`aws.log.stream.names`][res-attr-cfg]    | Toggles the `aws.log.stream.names` resource attribute. Sets `enabled` to `true` by default.    | no       |
| [`cloud.account.id`][res-attr-cfg]        | Toggles the `cloud.account.id` resource attribute. Sets `enabled` to `true` by default.        | no       |
| [`cloud.availability_zone`][res-attr-cfg] | Toggles the `cloud.availability_zone` resource attribute. Sets `enabled` to `true` by default. | no       |
| [`cloud.platform`][res-attr-cfg]          | Toggles the `cloud.platform` resource attribute. Sets `enabled` to `true` by default.          | no       |
| [`cloud.provider`][res-attr-cfg]          | Toggles the `cloud.provider` resource attribute. Sets `enabled` to `true` by default.          | no       |
| [`cloud.region`][res-attr-cfg]            | Toggles the `cloud.region` resource attribute. Sets `enabled` to `true` by default.            | no       |

### `eks`

The `eks` block adds resource attributes for Amazon EKS.

The `eks` block supports the following blocks:

| Block                                              | Description                                  | Required |
|----------------------------------------------------|----------------------------------------------|----------|
| [`resource_attributes`](#eks--resource_attributes) | Configures which resource attributes to add. | no       |

#### `eks` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                              | Description                                                                              | Required |
|------------------------------------|------------------------------------------------------------------------------------------|----------|
| [`cloud.account.id`][res-attr-cfg] | Toggles the `cloud.account.id` resource attribute. Sets `enabled` to `false` by default. | no       |
| [`cloud.platform`][res-attr-cfg]   | Toggles the `cloud.platform` resource attribute. Sets `enabled` to `true` by default.    | no       |
| [`cloud.provider`][res-attr-cfg]   | Toggles the `cloud.provider` resource attribute. Sets `enabled` to `true` by default.    | no       |
| [`k8s.cluster.name`][res-attr-cfg] | Toggles the `k8s.cluster.name` resource attribute. Sets `enabled` to `false` by default. | no       |

Example values:

* `cloud.provider`: `"aws"`
* `cloud.platform`: `"aws_eks"`

### `elasticbeanstalk`

The `elasticbeanstalk` block reads the AWS X-Ray configuration file available on all Beanstalk instances with [X-Ray Enabled][].

[X-Ray Enabled]: https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/environment-configuration-debugging.html

The `elasticbeanstalk` block supports the following blocks:

| Block                                                           | Description                                  | Required |
|-----------------------------------------------------------------|----------------------------------------------|----------|
| [`resource_attributes`](#elasticbeanstalk--resource_attributes) | Configures which resource attributes to add. | no       |

#### `elasticbeanstalk` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                              | Description                                                                             | Required |
|------------------------------------|-----------------------------------------------------------------------------------------|----------|
| [`cloud.platform`][res-attr-cfg]   | Toggles the `cloud.platform` resource attribute. Sets `enabled` to `true` by default.   | no       |
| [`cloud.provider`][res-attr-cfg]   | Toggles the `cloud.provider` resource attribute. Sets `enabled` to `true` by default.   | no       |
| [`deployment.envir`][res-attr-cfg] | Toggles the `deployment.envir` resource attribute. Sets `enabled` to `true` by default. | no       |
| [`service.instance`][res-attr-cfg] | Toggles the `service.instance` resource attribute. Sets `enabled` to `true` by default. | no       |
| [`service.version`][res-attr-cfg]  | Toggles the `service.version` resource attribute. Sets `enabled` to `true` by default.  | no       |

Example values:

* `cloud.provider`: `"aws"`
* `cloud.platform`: `"aws_elastic_beanstalk"`

### `gcp`

The `gcp` block detects resource attributes using the [Google Cloud Client Libraries for Go][], which reads resource information from the [GCP metadata server][].
The detector also uses environment variables to identify which GCP platform the application is running on, and assigns appropriate resource attributes for that platform.

Use the `gcp` detector regardless of the GCP platform {{< param "PRODUCT_NAME" >}} is running on.

[Google Cloud Client Libraries for Go]: https://github.com/googleapis/google-cloud-go
[GCP metadata server]: https://cloud.google.com/compute/docs/storing-retrieving-metadata

The `gcp` block supports the following blocks:

| Block                                              | Description                                  | Required |
|----------------------------------------------------|----------------------------------------------|----------|
| [`resource_attributes`](#gcp--resource_attributes) | Configures which resource attributes to add. | no       |

#### `gcp` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                                                   | Description                                                                                                  | Required |
|---------------------------------------------------------|--------------------------------------------------------------------------------------------------------------|----------|
| [`cloud.account.id`][res-attr-cfg]                      | Toggles the `cloud.account.id` resource attribute. Sets `enabled` to `true` by default.                      | no       |
| [`cloud.availability_zone`][res-attr-cfg]               | Toggles the `cloud.availability_zone` resource attribute. Sets `enabled` to `true` by default.               | no       |
| [`cloud.platform`][res-attr-cfg]                        | Toggles the `cloud.platform` resource attribute. Sets `enabled` to `true` by default.                        | no       |
| [`cloud.provider`][res-attr-cfg]                        | Toggles the `cloud.provider` resource attribute. Sets `enabled` to `true` by default.                        | no       |
| [`cloud.region`][res-attr-cfg]                          | Toggles the `cloud.region` resource attribute. Sets `enabled` to `true` by default.                          | no       |
| [`faas.id`][res-attr-cfg]                               | Toggles the `faas.id` resource attribute. Sets `enabled` to `true` by default.                               | no       |
| [`faas.instance`][res-attr-cfg]                         | Toggles the `faas.instance` resource attribute. Sets `enabled` to `true` by default.                         | no       |
| [`faas.name`][res-attr-cfg]                             | Toggles the `faas.name` resource attribute. Sets `enabled` to `true` by default.                             | no       |
| [`faas.version`][res-attr-cfg]                          | Toggles the `faas.version` resource attribute. Sets `enabled` to `true` by default.                          | no       |
| [`gcp.cloud_run.job.execution`][res-attr-cfg]           | Toggles the `gcp.cloud_run.job.execution` resource attribute. Sets `enabled` to `true` by default.           | no       |
| [`gcp.cloud_run.job.task_index`][res-attr-cfg]          | Toggles the `gcp.cloud_run.job.task_index` resource attribute. Sets `enabled` to `true` by default.          | no       |
| [`gcp.gce.instance.group_manager.name`][res-attr-cfg]   | Toggles the `gcp.gce.instance.group_manager.name` resource attribute. Sets `enabled` to `true` by default.   | no       |
| [`gcp.gce.instance.group_manager.region`][res-attr-cfg] | Toggles the `gcp.gce.instance.group_manager.region` resource attribute. Sets `enabled` to `true` by default. | no       |
| [`gcp.gce.instance.group_manager.zone`][res-attr-cfg]   | Toggles the `gcp.gce.instance.group_manager.zone` resource attribute. Sets `enabled` to `true` by default.   | no       |
| [`gcp.gce.instance.hostname`][res-attr-cfg]             | Toggles the `gcp.gce.instance.hostname` resource attribute. Sets `enabled` to `false` by default.            | no       |
| [`gcp.gce.instance.name`][res-attr-cfg]                 | Toggles the `gcp.gce.instance.name` resource attribute. Sets `enabled` to `false` by default.                | no       |
| [`host.id`][res-attr-cfg]                               | Toggles the `host.id` resource attribute. Sets `enabled` to `true` by default.                               | no       |
| [`host.name`][res-attr-cfg]                             | Toggles the `host.name` resource attribute. Sets `enabled` to `true` by default.                             | no       |
| [`host.type`][res-attr-cfg]                             | Toggles the `host.type` resource attribute. Sets `enabled` to `true` by default.                             | no       |
| [`k8s.cluster.name`][res-attr-cfg]                      | Toggles the `k8s.cluster.name` resource attribute. Sets `enabled` to `true` by default.                      | no       |

#### Google Compute Engine (GCE) metadata

* `cloud.provider`: `"gcp"`
* `cloud.platform`: `"gcp_compute_engine"`
* `cloud.account.id`: Project ID
* `cloud.region`: For example, `"us-central1"`
* `cloud.availability_zone`: For example, `"us-central1-c"`
* `host.id`: Instance ID
* `host.name`: Instance name
* `host.type`: Machine type
* (optional) `gcp.gce.instance.hostname`
* (optional) `gcp.gce.instance.name`
* `gcp.gce.instance.group_manager.name`:  Managed instance group name
* `gcp.gce.instance.group_manager.region`:  Managed instance group region
* `gcp.gce.instance.group_manager.zone`:  Managed instance group zone

#### Google Kubernetes Engine (GKE) metadata

* `cloud.provider`: `"gcp"`
* `cloud.platform`: `"gcp_kubernetes_engine"`
* `cloud.account.id`: Project ID
* `cloud.region`: Only for regional GKE clusters, for example `"us-central1"`
* `cloud.availability_zone`: only for zonal GKE clusters, for example, `"us-central1-c"`
* `k8s.cluster.name`
* `host.id`: Instance ID
* `host.name`: Instance name, only when workload identity is disabled

One known issue happens when GKE workload identity is enabled.
The GCE metadata endpoints won't be available, and the GKE resource detector won't be able to determine `host.name`.
If this happens, you can set `host.name` from one of the following resources:

* Get the `node.name` through the [downward API][] with the `env` detector.
* Get the Kubernetes node name from the Kubernetes API (with `k8s.io/client-go`).

[downward API]: https://kubernetes.io/docs/concepts/workloads/pods/downward-api/

#### Google Cloud Run Services metadata

* `cloud.provider`: `"gcp"`
* `cloud.platform`: `"gcp_cloud_run"`
* `cloud.account.id`: Project ID
* `cloud.region`: For example, `"us-central1"`
* `faas.id`: Instance ID
* `faas.name`: Service name
* `faas.version`: Service revision

#### Cloud Run Jobs metadata

* `cloud.provider`: `"gcp"`
* `cloud.platform`: `"gcp_cloud_run"`
* `cloud.account.id`: Project ID
* `cloud.region`: For example, `"us-central1"`
* `faas.id`: Instance ID
* `faas.name`: Service name
* `gcp.cloud_run.job.execution`: For example, `"my-service-ajg89"`
* `gcp.cloud_run.job.task_index`: For example, `"0"`

#### Google Cloud Functions metadata

* `cloud.provider`: `"gcp"`
* `cloud.platform`: `"gcp_cloud_functions"`
* `cloud.account.id`: Project ID
* `cloud.region`: For example, `"us-central1"`
* `faas.id`: Instance ID
* `faas.name`: Function name
* `faas.version`: Function version

#### Google App Engine metadata

* `cloud.provider`: `"gcp"`
* `cloud.platform`: `"gcp_app_engine"`
* `cloud.account.id`: Project ID
* `cloud.region`: For example, `"us-central1"`
* `cloud.availability_zone`: For example, `"us-central1-c"`
* `faas.id`: Instance ID
* `faas.name`: Service name
* `faas.version`: Service version

### `heroku`

The `heroku` block adds resource attributes derived from [Heroku dyno metadata][].

The `heroku` block supports the following blocks:

| Block                                                 | Description                                  | Required |
|-------------------------------------------------------|----------------------------------------------|----------|
| [`resource_attributes`](#heroku--resource_attributes) | Configures which resource attributes to add. | no       |

#### `heroku` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                                               | Description                                                                                              | Required |
|-----------------------------------------------------|----------------------------------------------------------------------------------------------------------|----------|
| [`cloud.provider`][res-attr-cfg]                    | Toggles the `cloud.provider` resource attribute. Sets `enabled` to `true` by default.                    | no       |
| [`heroku.app.id`][res-attr-cfg]                     | Toggles the `heroku.app.id` resource attribute. Sets `enabled` to `true` by default.                     | no       |
| [`heroku.dyno.id`][res-attr-cfg]                    | Toggles the `heroku.dyno.id` resource attribute. Sets `enabled` to `true` by default.                    | no       |
| [`heroku.release.commit`][res-attr-cfg]             | Toggles the `heroku.release.commit` resource attribute. Sets `enabled` to `true` by default.             | no       |
| [`heroku.release.creation_timestamp`][res-attr-cfg] | Toggles the `heroku.release.creation_timestamp` resource attribute. Sets `enabled` to `true` by default. | no       |
| [`service.instance.id`][res-attr-cfg]               | Toggles the `service.instance.id` resource attribute. Sets `enabled` to `true` by default.               | no       |
| [`service.name`][res-attr-cfg]                      | Toggles the `service.name` resource attribute. Sets `enabled` to `true` by default.                      | no       |
| [`service.version`][res-attr-cfg]                   | Toggles the `service.version` resource attribute. Sets `enabled` to `true` by default.                   | no       |

When [Heroku dyno metadata][] is active, Heroku applications publish information through environment variables.
These environment variables map to resource attributes as follows:

| Dyno metadata environment variable | Resource attribute                  |
|------------------------------------|-------------------------------------|
| `HEROKU_APP_ID`                    | `heroku.app.id`                     |
| `HEROKU_APP_NAME`                  | `service.name`                      |
| `HEROKU_DYNO_ID`                   | `service.instance.id`               |
| `HEROKU_RELEASE_CREATED_AT`        | `heroku.release.creation_timestamp` |
| `HEROKU_RELEASE_VERSION`           | `service.version`                   |
| `HEROKU_SLUG_COMMIT`               | `heroku.release.commit`             |

For more information, refer to the [Heroku cloud provider documentation][] under the [OpenTelemetry specification semantic conventions][].

[Heroku dyno metadata]: https://devcenter.heroku.com/articles/dyno-metadata
[Heroku cloud provider documentation]: https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/resource/semantic_conventions/cloud_provider/heroku.md
[OpenTelemetry specification semantic conventions]: https://github.com/open-telemetry/opentelemetry-specification

### `kubeadm`

The `kubeadm` block queries the Kubernetes API server to retrieve kubeadm resource attributes.

The `kubeadm` block supports the following attributes:

| Attribute   | Type     | Description                                                             | Default  | Required |
|-------------|----------|-------------------------------------------------------------------------|----------|----------|
| `auth_type` | `string` | Configures how to authenticate to the Kubernetes API server.            | `"none"` | no       |
| `context`   | `string` | Override the current context when `auth_type` is set to `"kubeConfig"`. | `""`     | no       |

The following permissions are required:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: otel-collector
  namespace: kube-system
rules:
  - apiGroups: [""]
    resources: ["configmaps"]
    resourceNames: ["kubeadm-config"]
    verbs: ["get"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: otel-collector-rolebinding
  namespace: kube-system
subjects:
- kind: ServiceAccount
  name: default
  namespace: default
roleRef:
  kind: Role
  name: otel-collector
  apiGroup: rbac.authorization.k8s.io
```

You can set `auth_type` to one of the following:

* `none`: No authentication.
* `serviceAccount`: Use the standard service account token provided to the {{< param "PRODUCT_NAME" >}} Pod.
* `kubeConfig`: Use credentials from `~/.kube/config`.

The `kubeadm` block supports the following blocks:

| Block                                                  | Description                                  | Required |
|--------------------------------------------------------|----------------------------------------------|----------|
| [`resource_attributes`](#kubeadm--resource_attributes) | Configures which resource attributes to add. | no       |

#### `kubeadm` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                              | Description                                                                             | Required |
|------------------------------------|-----------------------------------------------------------------------------------------|----------|
| [`k8s.cluster.name`][res-attr-cfg] | Toggles the `k8s.cluster.name` resource attribute. Sets `enabled` to `true` by default. | no       |
| [`k8s.cluster.uid`][res-attr-cfg]  | Toggles the `k8s.cluster.uid` resource attribute. Sets `enabled` to `true` by default.  | no       |

### `kubernetes_node`

The `kubernetes_node` block queries the Kubernetes API server to retrieve various node resource attributes.

The `kubernetes_node` block supports the following attributes:

| Attribute           | Type     | Description                                                               | Default           | Required |
|---------------------|----------|---------------------------------------------------------------------------|-------------------|----------|
| `auth_type`         | `string` | Configures how to authenticate to the K8s API server.                     | `"none"`          | no       |
| `context`           | `string` | Override the current context when `auth_type` is set to `"kubeConfig"`.   | `""`              | no       |
| `node_from_env_var` | `string` | The name of an environment variable from which to retrieve the node name. | `"K8S_NODE_NAME"` | no       |

The "get" and "list" permissions are required:

```yaml
kind: ClusterRole
metadata:
  name: alloy
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list"]
```

`auth_type` can be set to one of the following:

* `none`: No authentication.
* `serviceAccount`: Use the standard service account token provided to the {{< param "PRODUCT_NAME" >}} Pod.
* `kubeConfig`: Use credentials from `~/.kube/config`.

The `kubernetes_node` block supports the following blocks:

| Block                                                          | Description                                  | Required |
|----------------------------------------------------------------|----------------------------------------------|----------|
| [`resource_attributes`](#kubernetes_node--resource_attributes) | Configures which resource attributes to add. | no       |

#### `kubernetes_node` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                           | Description                                                                          | Required |
|---------------------------------|--------------------------------------------------------------------------------------|----------|
| [`k8s.node.name`][res-attr-cfg] | Toggles the `k8s.node.name` resource attribute. Sets `enabled` to `true` by default. | no       |
| [`k8s.node.uid`][res-attr-cfg]  | Toggles the `k8s.node.uid` resource attribute. Sets `enabled` to `true` by default.  | no       |

### `lambda`

The `lambda` block uses the AWS Lambda [runtime environment variables][lambda-env-vars] to retrieve various resource attributes.

[lambda-env-vars]: https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html#configuration-envvars-runtime

The `lambda` block supports the following blocks:

| Block                                                | Description                                  | Required |
| ---------------------------------------------------- | -------------------------------------------- | -------- |
| [resource_attri`butes](#lambda--resource_attributes) | Configures which resource attributes to add. | no       |

#### `lambda` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                                  | Description                                                                                 | Required |
|----------------------------------------|---------------------------------------------------------------------------------------------|----------|
| [`aws.log.group.names`][res-attr-cfg]  | Toggles the `aws.log.group.names` resource attribute. Sets `enabled` to `true` by default.  | no       |
| [`aws.log.stream.names`][res-attr-cfg] | Toggles the `aws.log.stream.names` resource attribute. Sets `enabled` to `true` by default. | no       |
| [`cloud.platform`][res-attr-cfg]       | Toggles the `cloud.platform` resource attribute. Sets `enabled` to `true` by default.       | no       |
| [`cloud.provider`][res-attr-cfg]       | Toggles the `cloud.provider` resource attribute. Sets `enabled` to `true` by default.       | no       |
| [`cloud.region`][res-attr-cfg]         | Toggles the `cloud.region` resource attribute. Sets `enabled` to `true` by default.         | no       |
| [`faas.instance`][res-attr-cfg]        | Toggles the `faas.instance` resource attribute. Sets `enabled` to `true` by default.        | no       |
| [`faas.max_memory`][res-attr-cfg]      | Toggles the `faas.max_memory` resource attribute. Sets `enabled` to `true` by default.      | no       |
| [`faas.name`][res-attr-cfg]            | Toggles the `faas.name` resource attribute. Sets `enabled` to `true` by default.            | no       |
| [`faas.version`][res-attr-cfg]         | Toggles the `faas.version` resource attribute. Sets `enabled` to `true` by default.         | no       |

[Cloud semantic conventions][]:

* `cloud.provider`: `"aws"`
* `cloud.platform`: `"aws_lambda"`
* `cloud.region`: `$AWS_REGION`

[Function as a Service semantic conventions][] and [AWS Lambda semantic conventions][]:

* `faas.name`: `$AWS_LAMBDA_FUNCTION_NAME`
* `faas.version`: `$AWS_LAMBDA_FUNCTION_VERSION`
* `faas.instance`: `$AWS_LAMBDA_LOG_STREAM_NAME`
* `faas.max_memory`: `$AWS_LAMBDA_FUNCTION_MEMORY_SIZE`

[AWS Logs semantic conventions][]:

* `aws.log.group.names`: `$AWS_LAMBDA_LOG_GROUP_NAME`
* `aws.log.stream.names`: `$AWS_LAMBDA_LOG_STREAM_NAME`

[Cloud semantic conventions]: https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/resource/semantic_conventions/cloud.md
[Function as a Service semantic conventions]: https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/resource/semantic_conventions/faas.md
[AWS Lambda semantic conventions]: https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/instrumentation/aws-lambda.md#resource-detector
[AWS Logs semantic conventions]: https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/resource/semantic_conventions/cloud_provider/aws/logs.md

### `openshift`

The `openshift` block queries the OpenShift and Kubernetes APIs to retrieve various resource attributes.

The `openshift` block supports the following attributes:

| Attribute | Type     | Description                                              | Default     | Required |
|-----------|----------|----------------------------------------------------------|-------------|----------|
| `address` | `string` | Address of the OpenShift API server.                     | _See below_ | no       |
| `token`   | `string` | Token used to identify against the OpenShift API server. | ""          | no       |

The "get", "watch", and "list" permissions are required:

```yaml
kind: ClusterRole
metadata:
  name: alloy
rules:
- apiGroups: ["config.openshift.io"]
  resources: ["infrastructures", "infrastructures/status"]
  verbs: ["get", "watch", "list"]
```

By default, the API address is determined from the environment variables `KUBERNETES_SERVICE_HOST`, `KUBERNETES_SERVICE_PORT` and the service token is read from `/var/run/secrets/kubernetes.io/serviceaccount/token`.
If TLS isn't explicitly disabled and no `ca_file` is configured, `/var/run/secrets/kubernetes.io/serviceaccount/ca.crt` is used.
The determination of the API address, `ca_file`, and the service token is skipped if they are set in the configuration.

The `openshift` block supports the following blocks:

| Block                                                    | Description                                             | Required |
|----------------------------------------------------------|---------------------------------------------------------|----------|
| [`tls`](#openshift--tls)                                 | TLS settings for the connection with the OpenShift API. | yes      |
| [`resource_attributes`](#openshift--resource_attributes) | Configures which resource attributes to add.            | no       |

#### `openshift` > `tls`

The `tls` block configures TLS settings used for the connection to the gRPC server.

{{< docs/shared lookup="reference/components/otelcol-tls-client-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `tpm`

The `tpm` block configures retrieving the TLS `key_file` from a trusted device.

{{< docs/shared lookup="reference/components/otelcol-tls-tpm-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

#### `openshift` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                              | Description                                                                             | Required |
|------------------------------------|-----------------------------------------------------------------------------------------|----------|
| [`cloud.platform`][res-attr-cfg]   | Toggles the `cloud.platform` resource attribute. Sets `enabled` to `true` by default.   | no       |
| [`cloud.provider`][res-attr-cfg]   | Toggles the `cloud.provider` resource attribute. Sets `enabled` to `true` by default.   | no       |
| [`cloud.region`][res-attr-cfg]     | Toggles the `cloud.region` resource attribute. Sets `enabled` to `true` by default.     | no       |
| [`k8s.cluster.name`][res-attr-cfg] | Toggles the `k8s.cluster.name` resource attribute. Sets `enabled` to `true` by default. | no       |

### `system`

The `system` block queries the host machine to retrieve various resource attributes.

{{< admonition type="note" >}}
Use the [Docker](#docker) detector if running {{< param "PRODUCT_NAME" >}} as a Docker container.
{{< /admonition >}}

The `system` block supports the following attributes:

| Attribute          | Type           | Description                                                         | Default         | Required |
|--------------------|----------------|---------------------------------------------------------------------|-----------------|----------|
| `hostname_sources` | `list(string)` | A priority list of sources from which the hostname will be fetched. | `["dns", "os"]` | no       |

The valid options for `hostname_sources` are:

* `"dns"`: Uses multiple sources to get the fully qualified domain name.
  First, it looks up the host name in the local machine's `hosts` file.
  If that fails, it looks up the CNAME.
  If the CNAME lookup fails, it does a reverse DNS query.
  This hostname source may produce unreliable results on Windows.
  To produce a FQDN, Windows hosts might have better results using the "lookup" hostname source.
* `"os"`: Provides the hostname provided by the local machine's kernel.
* `"cname"`: Provides the canonical name, as provided by `net.LookupCNAME` in the Go standard library.
  This hostname source may produce unreliable results on Windows.
* `"lookup"`: Does a reverse DNS lookup of the current host's IP address.

If there is an error fetching a hostname from a source, the next source from the list of `hostname_sources` will be considered.

The `system` block supports the following blocks:

| Block                                                 | Description                                  | Required |
|-------------------------------------------------------|----------------------------------------------|----------|
| [`resource_attributes`](#system--resource_attributes) | Configures which resource attributes to add. | no       |

#### `system` > `resource_attributes`

The `resource_attributes` block supports the following blocks:

| Block                                    | Description                                                                                    | Required |
|------------------------------------------|------------------------------------------------------------------------------------------------|----------|
| [`host.arch`][res-attr-cfg]              | Toggles the `host.arch` resource attribute. Sets `enabled` to `false` by default.              | no       |
| [`host.cpu.cache.l2.size`][res-attr-cfg] | Toggles the `host.cpu.cache.l2.size` resource attribute. Sets `enabled` to `false` by default. | no       |
| [`host.cpu.family`][res-attr-cfg]        | Toggles the `host.cpu.family` resource attribute. Sets `enabled` to `false` by default.        | no       |
| [`host.cpu.model.id`][res-attr-cfg]      | Toggles the `host.cpu.model.id` resource attribute. Sets `enabled` to `false` by default.      | no       |
| [`host.cpu.model.name`][res-attr-cfg]    | Toggles the `host.cpu.model.name` resource attribute. Sets `enabled` to `false` by default.    | no       |
| [`host.cpu.stepping`][res-attr-cfg]      | Toggles the `host.cpu.stepping` resource attribute. Sets `enabled` to `false` by default.      | no       |
| [`host.cpu.vendor.id`][res-attr-cfg]     | Toggles the `host.cpu.vendor.id` resource attribute. Sets `enabled` to `false` by default.     | no       |
| [`host.id`][res-attr-cfg]                | Toggles the `host.id` resource attribute. Sets `enabled` to `false` by default.                | no       |
| [`host.interface`][res-attr-cfg]         | Toggles the `host.interface` resource attribute. Sets `enabled` to `false` by default.         | no       |
| [`host.ip`][res-attr-cfg]                | Toggles the `host.ip` resource attribute. Sets `enabled` to `false` by default.                | no       |
| [`host.mac`][res-attr-cfg]               | Toggles the `host.mac` resource attribute. Sets `enabled` to `false` by default.               | no       |
| [`host.name`][res-attr-cfg]              | Toggles the `host.name` resource attribute. Sets `enabled` to `true` by default.               | no       |
| [`os.build.id`][res-attr-cfg]            | Toggles the `os.build.id` resource attribute. Sets `enabled` to `false` by default.            | no       |
| [`os.description`][res-attr-cfg]         | Toggles the `os.description` resource attribute. Sets `enabled` to `false` by default.         | no       |
| [`os.name`][res-attr-cfg]                | Toggles the `os.name` resource attribute. Sets `enabled` to `false` by default.                | no       |
| [`os.type`][res-attr-cfg]                | Toggles the `os.type` resource attribute. Sets `enabled` to `true` by default.                 | no       |
| [`os.version`][res-attr-cfg]             | Toggles the `os.version` resource attribute. Sets `enabled` to `false` by default.             | no       |

## Common configuration

### Resource attribute configuration

This block describes how to configure resource attributes such as `k8s.node.name` and `azure.vm.name`.
Every block is configured using the same set of attributes.
Only the default values for those attributes might differ across resource attributes.
For example, some resource attributes have `enabled` set to `true` by default, whereas others don't.

The following attributes are supported:

| Attribute | Type   | Description                                                                         | Default     | Required |
|-----------|--------|-------------------------------------------------------------------------------------|-------------|----------|
| `enabled` | `bool` | Toggles whether to add the resource attribute to the span, log, or metric resource. | _See below_ | no       |

To see the default value for `enabled`, refer to the tables in the sections above which list the resource attributes blocks.
The "Description" column will state either:

> Sets `enabled` to `true` by default.

or:

> Sets `enabled` to `false` by default.

## Exported fields

The following fields are exported and can be referenced by other components:

| Name    | Type               | Description                                                      |
|---------|--------------------|------------------------------------------------------------------|
| `input` | `otelcol.Consumer` | A value that other components can use to send telemetry data to. |

`input` accepts `otelcol.Consumer` OTLP-formatted data for any telemetry signal of these types:

* logs
* metrics
* traces

## Component health

`otelcol.processor.resourcedetection` is only reported as unhealthy if given an invalid configuration.

## Debug information

`otelcol.processor.resourcedetection` doesn't expose any component-specific debug information.

## Examples

### `env` detector

If you set up a `OTEL_RESOURCE_ATTRIBUTES` environment variable with value of `TestKey=TestValue`,
then all logs, metrics, and traces have a resource attribute with a key `TestKey` and value of `TestValue`.

```alloy
otelcol.processor.resourcedetection "default" {
  detectors = ["env"]

  output {
    logs    = [otelcol.exporter.otlp.default.input]
    metrics = [otelcol.exporter.otlp.default.input]
    traces  = [otelcol.exporter.otlp.default.input]
  }
}
```

### `env` and `ec2`

There is no need to put in an `ec2 {}` block.
The `ec2` defaults are applied automatically, as specified in [`ec2`][ec2].

```alloy
otelcol.processor.resourcedetection "default" {
  detectors = ["env", "ec2"]

  output {
    logs    = [otelcol.exporter.otlp.default.input]
    metrics = [otelcol.exporter.otlp.default.input]
    traces  = [otelcol.exporter.otlp.default.input]
  }
}
```

### `ec2` with default resource attributes

There is no need to put in a `ec2 {}` block.
The `ec2` defaults are applied automatically, as specified in [`ec2`][ec2].

```alloy
otelcol.processor.resourcedetection "default" {
  detectors = ["ec2"]

  output {
    logs    = [otelcol.exporter.otlp.default.input]
    metrics = [otelcol.exporter.otlp.default.input]
    traces  = [otelcol.exporter.otlp.default.input]
  }
}
```

### `ec2` with explicit resource attributes

```alloy
otelcol.processor.resourcedetection "default" {
  detectors = ["ec2"]
  ec2 {
    tags = ["^tag1$", "^tag2$", "^label.*$"]
    resource_attributes {
      cloud.account.id  { enabled = true }
      cloud.availability_zone  { enabled = true }
      cloud.platform  { enabled = true }
      cloud.provider  { enabled = true }
      cloud.region  { enabled = true }
      host.id  { enabled = true }
      host.image.id  { enabled = false }
      host.name  { enabled = false }
      host.type  { enabled = false }
    }
  }

  output {
    logs    = [otelcol.exporter.otlp.default.input]
    metrics = [otelcol.exporter.otlp.default.input]
    traces  = [otelcol.exporter.otlp.default.input]
  }
}
```

### `kubernetes_node` with a default node

This example uses the default `node_from_env_var` option of `K8S_NODE_NAME`.

There is no need to put in a `kubernetes_node {}` block.
The `kubernetes_node` defaults are applied automatically, as specified in [`kubernetes_node`][kubernetes_node].

```alloy
otelcol.processor.resourcedetection "default" {
  detectors = ["kubernetes_node"]

  output {
    logs    = [otelcol.exporter.otlp.default.input]
    metrics = [otelcol.exporter.otlp.default.input]
    traces  = [otelcol.exporter.otlp.default.input]
  }
}
```

You need to add this to your workload:

```yaml
        env:
          - name: K8S_NODE_NAME
            valueFrom:
              fieldRef:
                fieldPath: spec.nodeName
```

### `kubernetes_node` with a custom environment variable

This example uses a custom `node_from_env_var` set to `my_custom_var`.

```alloy
otelcol.processor.resourcedetection "default" {
  detectors = ["kubernetes_node"]
  kubernetes_node {
    node_from_env_var = "my_custom_var"
  }

  output {
    logs    = [otelcol.exporter.otlp.default.input]
    metrics = [otelcol.exporter.otlp.default.input]
    traces  = [otelcol.exporter.otlp.default.input]
  }
}
```

You need to add this to your workload:

```yaml
        env:
          - name: my_custom_var
            valueFrom:
              fieldRef:
                fieldPath: spec.nodeName
```
<!-- START GENERATED COMPATIBLE COMPONENTS -->

## Compatible components

`otelcol.processor.resourcedetection` can accept arguments from the following components:

- Components that export [OpenTelemetry `otelcol.Consumer`](../../../compatibility/#opentelemetry-otelcolconsumer-exporters)

`otelcol.processor.resourcedetection` has exports that can be consumed by the following components:

- Components that consume [OpenTelemetry `otelcol.Consumer`](../../../compatibility/#opentelemetry-otelcolconsumer-consumers)

{{< admonition type="note" >}}
Connecting some components may not be sensible or components may require further configuration to make the connection work correctly.
Refer to the linked documentation for more details.
{{< /admonition >}}

<!-- END GENERATED COMPATIBLE COMPONENTS -->
