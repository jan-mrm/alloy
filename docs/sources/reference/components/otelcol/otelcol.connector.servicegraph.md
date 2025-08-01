---
canonical: https://grafana.com/docs/alloy/latest/reference/components/otelcol/otelcol.connector.servicegraph/
aliases:
  - ../otelcol.connector.servicegraph/ # /docs/alloy/latest/reference/components/otelcol.connector.servicegraph/
description: Learn about otelcol.connector.servicegraph
labels:
  stage: general-availability
  products:
    - oss
title: otelcol.connector.servicegraph
---

# `otelcol.connector.servicegraph`

`otelcol.connector.servicegraph` accepts span data from other `otelcol` components and outputs metrics representing the relationship between various services in a system.
A metric represents an edge in the service graph.
Those metrics can then be used by a data visualization application, for example, [Grafana][], to draw the service graph.

[Grafana]: https://grafana.com/docs/grafana/latest/explore/trace-integration/#service-graph

{{< admonition type="note" >}}
`otelcol.connector.servicegraph` is a wrapper over the upstream OpenTelemetry Collector [`servicegraph`][] connector.
Bug reports or feature requests will be redirected to the upstream repository, if necessary.

[`servicegraph`]: https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/{{< param "OTEL_VERSION" >}}/connector/servicegraphconnector
{{< /admonition >}}

You can specify multiple `otelcol.connector.servicegraph` components by giving them different labels.

This component is based on the [Grafana Tempo service graph processor](https://github.com/grafana/tempo/tree/main/modules/generator/processor/servicegraphs).

Service graphs are useful for a number of use-cases:

* Infer the topology of a distributed system. As distributed systems grow, they become more complex.
  Service graphs can help you understand the structure of the system.
* Provide a high level overview of the health of your system.
  Service graphs show error rates, latencies, and other relevant data.
* Provide a historic view of a system's topology.
  Distributed systems change very frequently, and service graphs offer a way of seeing how these systems have evolved over time.

Since `otelcol.connector.servicegraph` has to process both sides of an edge, it needs to process all spans of a trace to function properly.
If spans of a trace are spread out over multiple {{< param "PRODUCT_NAME" >}} instances, spans can't be paired reliably.
A solution to this problem is using [otelcol.exporter.loadbalancing][] in front of {{< param "PRODUCT_NAME" >}} instances running `otelcol.connector.servicegraph`.

[otelcol.exporter.loadbalancing]: ../otelcol.exporter.loadbalancing/

## Usage

```alloy
otelcol.connector.servicegraph "<LABEL>" {
  output {
    metrics = [...]
  }
}
```

## Arguments

You can use the following arguments with `otelcol.connector.servicegraph`:

| Name                        | Type             | Description                                                                              | Default                                                                                                                      | Required |
|-----------------------------|------------------|------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------|----------|
| `cache_loop`                | `duration`       | Configures how often to delete series which have not been updated.                       | `"1m"`                                                                                                                       | no       |
| `database_name_attribute`   | `string`         | (Deprecated) The attribute name used to identify the database name from span attributes. | `"db.name"`                                                                                                                  | no       |
| `database_name_attributes`  | `list(string)`   | The list of attribute names used to identify the database name from span attributes.     | `["db.name"]`                                                                                                                | no       |
| `dimensions`                | `list(string)`   | A list of dimensions to add with the default dimensions.                                 | `[]`                                                                                                                         | no       |
| `latency_histogram_buckets` | `list(duration)` | Buckets for latency histogram metrics.                                                   | `["2ms", "4ms", "6ms", "8ms", "10ms", "50ms", "100ms", "200ms", "400ms", "800ms", "1s", "1400ms", "2s", "5s", "10s", "15s"]` | no       |
| `metrics_flush_interval`    | `duration`       | The interval at which metrics are flushed to downstream components.                      | `"60s"`                                                                                                                      | no       |
| `store_expiration_loop`     | `duration`       | The time to expire old entries from the store periodically.                              | `"2s"`                                                                                                                       | no       |

Service graphs work by inspecting traces and looking for spans with parent-children relationship that represent a request.
`otelcol.connector.servicegraph` uses OpenTelemetry semantic conventions to detect a myriad of requests.
The following requests are supported:

* A direct request between two services, where the outgoing and the incoming span must have a [Span Kind][] value of `client` and `server` respectively.
* A request across a messaging system, where the outgoing and the incoming span must have a [Span Kind][] value of `producer` and `consumer` respectively.
* A database request, where spans have a [Span Kind][] with a value of `client`, as well as an attribute with a key of `db.name`.

Every span which can be paired up to form a request is kept in an in-memory store:

* If the TTL of the span expires before it can be paired, it's deleted from the store.
  TTL is configured in the [store][] block.
* If the span is paired prior to its expiration, a metric is recorded and the span is deleted from the store.

The following metrics are emitted by the processor:

| Metric                                      | Type      | Labels                                | Description                                                               |
|---------------------------------------------|-----------|---------------------------------------|---------------------------------------------------------------------------|
| `traces_service_graph_dropped_spans_total`  | Counter   | `client`, `server`, `connection_type` | Total count of dropped spans                                              |
| `traces_service_graph_request_client`       | Histogram | `client`, `server`, `connection_type` | Number of seconds for a request between two nodes as seen from the client |
| `traces_service_graph_request_failed_total` | Counter   | `client`, `server`, `connection_type` | Total count of failed requests between two nodes                          |
| `traces_service_graph_request_server`       | Histogram | `client`, `server`, `connection_type` | Number of seconds for a request between two nodes as seen from the server |
| `traces_service_graph_request_total`        | Counter   | `client`, `server`, `connection_type` | Total count of requests between two nodes                                 |
| `traces_service_graph_unpaired_spans_total` | Counter   | `client`, `server`, `connection_type` | Total count of unpaired spans                                             |

Duration is measured both from the client and the server sides.

The `latency_histogram_buckets` argument controls the buckets for `traces_service_graph_request_server` and `traces_service_graph_request_client`.

Each emitted metrics series have a `client` and a `server` label corresponding with the service doing the request and the service receiving the request.
The value of the label is derived from the `service.name` resource attribute of the two spans.

The `connection_type` label may not be set. If it's set, its value will be either `messaging_system` or `database`.

Additional labels can be included using the `dimensions` configuration option:

* Those labels will have a prefix to mark where they originate (client or server span kinds).
  The `client_` prefix relates to the dimensions coming from spans with a [Span Kind][] of `client`.
  The `server_` prefix relates to the dimensions coming from spans with a [Span Kind][] of `server`.
* Firstly the resource attributes will be searched. If the attribute isn't found, the span attributes will be searched.

When `metrics_flush_interval` is set to `0s`, metrics will be flushed on every received batch of traces.

The attributes in `database_name_attributes` are tried in order, selecting the first match.

[Span Kind]: https://opentelemetry.io/docs/concepts/signals/traces/#span-kind

## Blocks

You can use the following blocks with `otelcol.connector.servicegraph`:

| Block                            | Description                                                                | Required |
|----------------------------------|----------------------------------------------------------------------------|----------|
| [`output`][output]               | Configures where to send telemetry data.                                   | yes      |
| [`debug_metrics`][debug_metrics] | Configures the metrics that this component generates to monitor its state. | no       |
| [`store`][store]                 | Configures the in-memory store for spans.                                  | no       |

[store]: #store
[output]: #output
[debug_metrics]: #debug_metrics

### `output`

{{< badge text="Required" >}}

{{< docs/shared lookup="reference/components/output-block-metrics.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `debug_metrics`

{{< docs/shared lookup="reference/components/otelcol-debug-metrics-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `store`

The `store` block configures the in-memory store for spans.

| Name        | Type       | Description                                   | Default | Required |
|-------------|------------|-----------------------------------------------|---------|----------|
| `max_items` | `number`   | Maximum number of items to keep in the store. | `1000`  | no       |
| `ttl`       | `duration` | The time to live for spans in the store.      | `"2s"`  | no       |

## Exported fields

The following fields are exported and can be referenced by other components:

| Name    | Type               | Description                                                      |
|---------|--------------------|------------------------------------------------------------------|
| `input` | `otelcol.Consumer` | A value that other components can use to send telemetry data to. |

`input` accepts `otelcol.Consumer` traces telemetry data.
It doesn't accept metrics and logs.

## Component health

`otelcol.connector.servicegraph` is only reported as unhealthy if given an invalid configuration.

## Debug information

`otelcol.connector.servicegraph` doesn't expose any component-specific debug information.

## Example

The example below accepts traces, creates service graph metrics from them, and writes the metrics to Mimir.
The traces are written to Tempo.

`otelcol.connector.servicegraph` also adds a label to each metric with the value of the "http.method" span/resource attribute.

```alloy
otelcol.receiver.otlp "default" {
  grpc {
    endpoint = "0.0.0.0:4320"
  }

  output {
    traces  = [otelcol.connector.servicegraph.default.input,otelcol.exporter.otlp.grafana_cloud_traces.input]
  }
}

otelcol.connector.servicegraph "default" {
  dimensions = ["http.method"]
  output {
    metrics = [otelcol.exporter.prometheus.default.input]
  }
}

otelcol.exporter.prometheus "default" {
  forward_to = [prometheus.remote_write.mimir.receiver]
}

prometheus.remote_write "mimir" {
  endpoint {
    url = "https://prometheus-xxx.grafana.net/api/prom/push"

    basic_auth {
      username = sys.env("<PROMETHEUS_USERNAME>")
      password = sys.env("<GRAFANA_CLOUD_API_KEY>")
    }
  }
}

otelcol.exporter.otlp "grafana_cloud_traces" {
  client {
    endpoint = "https://tempo-xxx.grafana.net/tempo"
    auth     = otelcol.auth.basic.grafana_cloud_traces.handler
  }
}

otelcol.auth.basic "grafana_cloud_traces" {
  username = sys.env("<TEMPO_USERNAME>")
  password = sys.env("<GRAFANA_CLOUD_API_KEY>")
}
```

Some of the metrics in Mimir may look like this:

```text
traces_service_graph_request_total{client="shop-backend",failed="false",server="article-service",client_http_method="DELETE",server_http_method="DELETE"}
traces_service_graph_request_failed_total{client="shop-backend",client_http_method="POST",failed="false",server="auth-service",server_http_method="POST"}
```

<!-- START GENERATED COMPATIBLE COMPONENTS -->

## Compatible components

`otelcol.connector.servicegraph` can accept arguments from the following components:

- Components that export [OpenTelemetry `otelcol.Consumer`](../../../compatibility/#opentelemetry-otelcolconsumer-exporters)

`otelcol.connector.servicegraph` has exports that can be consumed by the following components:

- Components that consume [OpenTelemetry `otelcol.Consumer`](../../../compatibility/#opentelemetry-otelcolconsumer-consumers)

{{< admonition type="note" >}}
Connecting some components may not be sensible or components may require further configuration to make the connection work correctly.
Refer to the linked documentation for more details.
{{< /admonition >}}

<!-- END GENERATED COMPATIBLE COMPONENTS -->
