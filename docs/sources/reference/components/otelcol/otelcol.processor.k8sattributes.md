---
canonical: https://grafana.com/docs/alloy/latest/reference/components/otelcol/otelcol.processor.k8sattributes/
aliases:
  - ../otelcol.processor.k8sattributes/ # /docs/alloy/latest/reference/otelcol.processor.k8sattributes/
description: Learn about otelcol.processor.k8sattributes
labels:
  stage: general-availability
  products:
    - oss
title: otelcol.processor.k8sattributes
---

# `otelcol.processor.k8sattributes`

`otelcol.processor.k8sattributes` accepts telemetry data from other `otelcol` components and adds Kubernetes metadata to the resource attributes of spans, logs, or metrics.

{{< admonition type="note" >}}
`otelcol.processor.k8sattributes` is a wrapper over the upstream OpenTelemetry Collector [`k8sattributes`][] processor.
If necessary, bug reports or feature requests will be redirected to the upstream repository.

[`k8sattributes`]: https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/{{< param "OTEL_VERSION" >}}/processor/k8sattributesprocessor
{{< /admonition >}}

You can specify multiple `otelcol.processor.k8sattributes` components by giving them different labels.

## Usage

```alloy
otelcol.processor.k8sattributes "<LABEL>" {
  output {
    metrics = [...]
    logs    = [...]
    traces  = [...]
  }
}
```

## Arguments

You can use the following arguments with `otelcol.processor.k8sattributes`:

| Name                        | Type       | Description                                                                    | Default            | Required |
|-----------------------------|------------|--------------------------------------------------------------------------------|--------------------|----------|
| `auth_type`                 | `string`   | Authentication method when connecting to the Kubernetes API.                   | `"serviceAccount"` | no       |
| `passthrough`               | `bool`     | Pass through signals as-is, only adding a `k8s.pod.ip` resource attribute.     | `false`            | no       |
| `wait_for_metadata_timeout` | `duration` | How long to wait for Kubernetes metadata to arrive.                            | `"10s"`            | no       |
| `wait_for_metadata`         | `bool`     | Whether to wait for Kubernetes metadata to arrive before processing telemetry. | `false`            | no       |

The supported values for `auth_type` are:

* `none`: No authentication is required.
* `serviceAccount`: Use the built-in service account that Kubernetes automatically provisions for each Pod.
* `kubeConfig`: Use local credentials like those used by `kubectl`.
* `tls`: Use client TLS authentication.

Setting `passthrough` to `true` enables the "passthrough mode" of `otelcol.processor.k8sattributes`:

* Only a `k8s.pod.ip` resource attribute will be added.
* No other metadata will be added.
* The Kubernetes API won't be accessed.
* To correctly detect the Pod IPs, {{< param "PRODUCT_NAME" >}} must receive spans directly from services.
* The `passthrough` setting is useful when configuring {{< param "PRODUCT_NAME" >}} as a Kubernetes Deployment.

A {{< param "PRODUCT_NAME" >}} running as a Deployment can't detect the IP addresses of pods generating telemetry data without any of the well-known IP attributes.
If the Deployment {{< param "PRODUCT_NAME" >}} receives telemetry from {{< param "PRODUCT_NAME" >}}s deployed as DaemonSet, then some of those attributes might be missing.
As a workaround, you can configure the DaemonSet {{< param "PRODUCT_NAME" >}}s with `passthrough` set to `true`.

By default, `otelcol.processor.k8sattributes` is ready as soon as it starts, even if no metadata has been fetched yet.
If telemetry is sent to this processor before the metadata is synced, there will be no metadata to enrich the telemetry with.

To wait for the metadata to be synced before `otelcol.processor.k8sattributes` is ready, set the `wait_for_metadata` option to `true`.
Then, the processor won't be ready until the metadata is fully synced. As a result, the start-up of {{< param "PRODUCT_NAME" >}} will be blocked.
If the metadata can't be synced by the time the `wait_for_metadata_timeout` duration is reached,
`otelcol.processor.k8sattributes` will become unhealthy and fail to start.

If `otelcol.processor.k8sattributes` is unhealthy, other {{< param "PRODUCT_NAME" >}} components will still be able to start.
However, they may be unable to send telemetry to `otelcol.processor.k8sattributes`.

## Blocks

You can use the following blocks with `otelcol.processor.k8sattributes`:

| Block                                  | Description                                                                | Required |
|----------------------------------------|----------------------------------------------------------------------------|----------|
| [`output`][output]                     | Configures where to send received telemetry data.                          | yes      |
| [`debug_metrics`][debug_metrics]       | Configures the metrics that this component generates to monitor its state. | no       |
| [`exclude`][exclude]                   | Exclude pods from being processed.                                         | no       |
| `exclude` > [`pod`][pod]               | Pod information.                                                           | no       |
| [`extract`][extract]                   | Rules for extracting data from Kubernetes.                                 | no       |
| `extract` > [`annotation`][annotation] | Creating resource attributes from Kubernetes annotations.                  | no       |
| `extract` > [`label`][extract_label]   | Creating resource attributes from Kubernetes labels.                       | no       |
| [`filter`][filter]                     | Filters the data loaded from Kubernetes.                                   | no       |
| `filter` > [`field`][field]            | Filter pods by generic Kubernetes fields.                                  | no       |
| `filter` > [`label`][filter_label]     | Filter pods by Kubernetes labels.                                          | no       |
| [`pod_association`][pod_association]   | Rules to associate Pod metadata with telemetry signals.                    | no       |
| `pod_association` > [`source`][source] | Source information to identify a Pod.                                      | no       |

The > symbol indicates deeper levels of nesting.
For example, `extract` > `annotation` refers to an `annotation` block defined inside an `extract` block.

[output]: #output
[extract]: #extract
[annotation]: #annotation
[extract_label]: #label-extract
[filter]: #filter
[field]: #field
[filter_label]: #label-filter
[pod_association]: #pod_association
[source]: #source
[exclude]: #exclude
[pod]: #pod
[debug_metrics]: #debug_metrics

### `output`

{{< badge text="Required" >}}

{{< docs/shared lookup="reference/components/output-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `debug_metrics`

{{< docs/shared lookup="reference/components/otelcol-debug-metrics-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `exclude`

The `exclude` block configures which pods to exclude from the processor.

{{< admonition type="note" >}}
Pods with the name `jaeger-agent` or `jaeger-collector` are excluded by default.
{{< /admonition >}}

### `pod`

The `pod` block configures a Pod to be excluded from the processor.

The following attributes are supported:

| Name   | Type     | Description         | Default | Required |
|--------|----------|---------------------|---------|----------|
| `name` | `string` | The name of the Pod |         | yes      |

### `extract`

The `extract` block configures which metadata, annotations, and labels to extract from the Pod.

The following attributes are supported:

| Name               | Type           | Description                                                                 | Default     | Required |
|--------------------|----------------|-----------------------------------------------------------------------------|-------------|----------|
| `metadata`         | `list(string)` | Pre-configured metadata keys to add.                                        | _See below_ | no       |
| `otel_annotations` | `bool`         | Whether to set the [recommended resource attributes][semantic conventions]. | `false`     | no       |

The supported `metadata` keys are:

* `container.id`
* `container.image.name`
* `container.image.tag`
* `k8s.container.name`
* `k8s.cronjob.name`
* `k8s.daemonset.name`
* `k8s.daemonset.uid`
* `k8s.deployment.name`
* `k8s.job.name`
* `k8s.job.uid`
* `k8s.namespace.name`
* `k8s.node.name`
* `k8s.pod.name`
* `k8s.pod.start_time`
* `k8s.pod.uid`
* `k8s.replicaset.name`
* `k8s.replicaset.uid`
* `k8s.statefulset.name`
* `k8s.statefulset.uid`
* `service.instance.id`
* `service.name`
* `service.namespace`
* `service.version`

The `service.*` metadata are calculated following the OpenTelemetry [semantic conventions][].

By default, if `metadata` isn't specified, the following fields are extracted and added to spans, metrics, and logs as resource attributes:

* `container.image.name` (requires one of the following additional attributes to be set: `container.id` or `k8s.container.name`)
* `container.image.tag` (requires one of the following additional attributes to be set: `container.id` or `k8s.container.name`)
* `k8s.container.name` (requires an additional attribute to be set: `container.id`)
* `k8s.deployment.name` (if the Pod is controlled by a deployment)
* `k8s.namespace.name`
* `k8s.node.name`
* `k8s.pod.name`
* `k8s.pod.start_time`
* `k8s.pod.uid`

When `otel_annotations` is set to `true`, annotations such as `resource.opentelemetry.io/exampleResource` will be translated to the `exampleResource` resource attribute, etc.

[semantic conventions]: https://opentelemetry.io/docs/specs/semconv/non-normative/k8s-attributes

### `annotation`

The `annotation` block configures how to extract Kubernetes annotations.

{{< docs/shared lookup="reference/components/extract-field-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

{{< admonition type="caution" >}}
The `regex` argument has been removed.
Use the [ExtractPatterns][extract-patterns] function from `otelcol.processor.transform` instead.

[extract-patterns]: https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/{{< param "OTEL_VERSION" >}}/pkg/ottl/ottlfuncs/README.md#extractpatterns
{{< /admonition >}}

### `label` extract

The `label` block configures how to extract Kubernetes labels.

{{< docs/shared lookup="reference/components/extract-field-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

{{< admonition type="caution" >}}

The `regex` argument has been removed.
Use the [ExtractPatterns][extract-patterns] function from `otelcol.processor.transform` instead.

[extract-patterns]: https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/{{< param "OTEL_VERSION" >}}/pkg/ottl/ottlfuncs/README.md#extractpatterns

{{< /admonition >}}

### `filter`

The `filter` block configures which nodes to get data from and which fields and labels to fetch.

The following attributes are supported:

| Name        | Type     | Description                                                             | Default | Required |
|-------------|----------|-------------------------------------------------------------------------|---------|----------|
| `node`      | `string` | Configures a Kubernetes node name or host name.                         | `""`    | no       |
| `namespace` | `string` | Filters all pods by the provided namespace. All other pods are ignored. | `""`    | no       |

If `node` is specified, then any pods not running on the specified node will be ignored by `otelcol.processor.k8sattributes`.

### `field`

The `field` block allows you to filter pods by generic Kubernetes fields.

{{< docs/shared lookup="reference/components/field-filter-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `label` filter

The `label` block allows you to filter pods by generic Kubernetes labels.

{{< docs/shared lookup="reference/components/field-filter-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `pod_association`

The `pod_association` block configures rules on how to associate logs/traces/metrics to pods.

The `pod_association` block doesn't support any arguments and is configured fully through child blocks.

The `pod_association` block can be repeated multiple times, to configure additional rules.

#### Example

```alloy
pod_association {
    source {
        from = "resource_attribute"
        name = "k8s.pod.ip"
    }
}

pod_association {
    source {
        from = "resource_attribute"
        name = "k8s.pod.uid"
    }
    source {
        from = "connection"
    }
}
```

### `source`

The `source` block configures a Pod association rule.
This is used by the `k8sattributes` processor to determine the Pod associated with a telemetry signal.

When multiple `source` blocks are specified inside a `pod_association` block, both `source` blocks has to match for the Pod to be associated with the telemetry signal.

The following attributes are supported:

| Name   | Type     | Description                                                                      | Default | Required |
|--------|----------|----------------------------------------------------------------------------------|---------|----------|
| `from` | `string` | The association method. Currently supports `resource_attribute` and `connection` |         | yes      |
| `name` | `string` | Name represents extracted key name. For example, `ip`, `pod_uid`, `k8s.pod.ip`   |         | no       |

## Exported fields

The following fields are exported and can be referenced by other components:

| Name    | Type               | Description                                                      |
|---------|--------------------|------------------------------------------------------------------|
| `input` | `otelcol.Consumer` | A value that other components can use to send telemetry data to. |

`input` accepts `otelcol.Consumer` data for any telemetry signal (metrics, logs, or traces).

## Component health

`otelcol.processor.k8sattributes` is only reported as unhealthy if given an invalid configuration.

## Debug information

`otelcol.processor.k8sattributes` doesn't expose any component-specific debug
information.

## Examples

### Basic usage

In most cases, this is enough to get started. It'll add these resource attributes to all logs, metrics, and traces:

* `k8s.deployment.name`
* `k8s.namespace.name`
* `k8s.node.name`
* `k8s.pod.name`
* `k8s.pod.start_time`
* `k8s.pod.uid`

Example:

```alloy
otelcol.receiver.otlp "default" {
  http {}
  grpc {}

  output {
    metrics = [otelcol.processor.k8sattributes.default.input]
    logs    = [otelcol.processor.k8sattributes.default.input]
    traces  = [otelcol.processor.k8sattributes.default.input]
  }
}

otelcol.processor.k8sattributes "default" {
  output {
    metrics = [otelcol.exporter.otlp.default.input]
    logs    = [otelcol.exporter.otlp.default.input]
    traces  = [otelcol.exporter.otlp.default.input]
  }
}

otelcol.exporter.otlp "default" {
  client {
    endpoint = sys.env("<OTLP_ENDPOINT>")
  }
}
```

### Add additional metadata and labels

```alloy
otelcol.receiver.otlp "default" {
  http {}
  grpc {}

  output {
    metrics = [otelcol.processor.k8sattributes.default.input]
    logs    = [otelcol.processor.k8sattributes.default.input]
    traces  = [otelcol.processor.k8sattributes.default.input]
  }
}

otelcol.processor.k8sattributes "default" {
  extract {
    label {
      from      = "pod"
      key_regex = "(.*)/(.*)"
      tag_name  = "$1.$2"
    }

    metadata = [
      "k8s.namespace.name",
      "k8s.deployment.name",
      "k8s.statefulset.name",
      "k8s.daemonset.name",
      "k8s.cronjob.name",
      "k8s.job.name",
      "k8s.node.name",
      "k8s.pod.name",
      "k8s.pod.uid",
      "k8s.pod.start_time",
    ]

    otel_annotations = true
  }

  output {
    metrics = [otelcol.exporter.otlp.default.input]
    logs    = [otelcol.exporter.otlp.default.input]
    traces  = [otelcol.exporter.otlp.default.input]
  }
}

otelcol.exporter.otlp "default" {
  client {
    endpoint = sys.env("<OTLP_ENDPOINT>")
  }
}
```

### Add Kubernetes metadata to Prometheus metrics

`otelcol.processor.k8sattributes` adds metadata to metrics signals in the form of resource attributes.
To display the metadata as labels of Prometheus metrics, the OTLP attributes must be converted from resource attributes to datapoint attributes.
One way to do this is by using an `otelcol.processor.transform` component.

```alloy
otelcol.receiver.otlp "default" {
  http {}
  grpc {}

  output {
    metrics = [otelcol.processor.k8sattributes.default.input]
  }
}

otelcol.processor.k8sattributes "default" {
  extract {
    label {
      from = "pod"
    }

    metadata = [
      "k8s.namespace.name",
      "k8s.pod.name",
    ]
  }

  output {
    metrics = [otelcol.processor.transform.add_kube_attrs.input]
  }
}

otelcol.processor.transform "add_kube_attrs" {
  error_mode = "ignore"

  metric_statements {
    context = "datapoint"
    statements = [
      "set(attributes[\"k8s.pod.name\"], resource.attributes[\"k8s.pod.name\"])",
      "set(attributes[\"k8s.namespace.name\"], resource.attributes[\"k8s.namespace.name\"])",
    ]
  }

  output {
    metrics = [otelcol.exporter.prometheus.default.input]
  }
}

otelcol.exporter.prometheus "default" {
  forward_to = [prometheus.remote_write.mimir.receiver]
}

prometheus.remote_write "mimir" {
  endpoint {
    url = "http://mimir:9009/api/v1/push"
  }
}
```

<!-- START GENERATED COMPATIBLE COMPONENTS -->

## Compatible components

`otelcol.processor.k8sattributes` can accept arguments from the following components:

- Components that export [OpenTelemetry `otelcol.Consumer`](../../../compatibility/#opentelemetry-otelcolconsumer-exporters)

`otelcol.processor.k8sattributes` has exports that can be consumed by the following components:

- Components that consume [OpenTelemetry `otelcol.Consumer`](../../../compatibility/#opentelemetry-otelcolconsumer-consumers)

{{< admonition type="note" >}}
Connecting some components may not be sensible or components may require further configuration to make the connection work correctly.
Refer to the linked documentation for more details.
{{< /admonition >}}

<!-- END GENERATED COMPATIBLE COMPONENTS -->
