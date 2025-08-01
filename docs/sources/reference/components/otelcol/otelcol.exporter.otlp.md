---
canonical: https://grafana.com/docs/alloy/latest/reference/components/otelcol/otelcol.exporter.otlp/
aliases:
  - ../otelcol.exporter.otlp/ # /docs/alloy/latest/reference/components/otelcol.exporter.otlp/
description: Learn about otelcol.exporter.otlp
labels:
  stage: general-availability
  products:
    - oss
title: otelcol.exporter.otlp
---

# `otelcol.exporter.otlp`

`otelcol.exporter.otlp` accepts telemetry data from other `otelcol` components and writes them over the network using the OTLP gRPC protocol.

{{< admonition type="note" >}}
`otelcol.exporter.otlp` is a wrapper over the upstream OpenTelemetry Collector [`otlp`][] exporter.
Bug reports or feature requests will be redirected to the upstream repository, if necessary.

[`otlp`]: https://github.com/open-telemetry/opentelemetry-collector/tree/{{< param "OTEL_VERSION" >}}/exporter/otlpexporter
{{< /admonition >}}

You can specify multiple `otelcol.exporter.otlp` components by giving them different labels.

## Usage

```alloy
otelcol.exporter.otlp "<LABEL>" {
  client {
    endpoint = "<HOST>:<PORT>"
  }
}
```

## Arguments

You can use the following argument with `otelcol.exporter.otlp`:

| Name      | Type       | Description                                      | Default | Required |
|-----------|------------|--------------------------------------------------|---------|----------|
| `timeout` | `duration` | Time to wait before marking a request as failed. | `"5s"`  | no       |

## Blocks

You can use the following blocks with `otelcol.exporter.otlp`:

| Block                                  | Description                                                                | Required |
|----------------------------------------|----------------------------------------------------------------------------|----------|
| [`client`][client]                     | Configures the gRPC client to send telemetry data to.                      | yes      |
| `client` > [`keepalive`][keepalive]    | Configures keepalive settings for the gRPC client.                         | no       |
| `client` > [`tls`][tls]                | Configures TLS for the gRPC client.                                        | no       |
| `client` > `tls` > [`tpm`][tpm]        | Configures TPM settings for the TLS key_file.                              | no       |
| [`debug_metrics`][debug_metrics]       | Configures the metrics that this component generates to monitor its state. | no       |
| [`retry_on_failure`][retry_on_failure] | Configures retry mechanism for failed requests.                            | no       |
| [`sending_queue`][sending_queue]       | Configures batching of data before sending.                                | no       |

The > symbol indicates deeper levels of nesting.
For example, `client` > `tls` refers to a `tls` block defined inside a `client` block.

[client]: #client
[tls]: #tls
[tpm]: #tpm
[keepalive]: #keepalive
[sending_queue]: #sending_queue
[retry_on_failure]: #retry_on_failure
[debug_metrics]: #debug_metrics

### `client`

{{< badge text="Required" >}}

The `client` block configures the gRPC client used by the component.

The following arguments are supported:

| Name                | Type                       | Description                                                                      | Default         | Required |
|---------------------|----------------------------|----------------------------------------------------------------------------------|-----------------|----------|
| `endpoint`          | `string`                   | `host:port` to send telemetry data to.                                           |                 | yes      |
| `auth`              | `capsule(otelcol.Handler)` | Handler from an `otelcol.auth` component to use for authenticating requests.     |                 | no       |
| `authority`         | `string`                   | Overrides the default `:authority` header in gRPC requests from the gRPC client. |                 | no       |
| `balancer_name`     | `string`                   | Which gRPC client-side load balancer to use for requests.                        | `"round_robin"` | no       |
| `compression`       | `string`                   | Compression mechanism to use for requests.                                       | `"gzip"`        | no       |
| `headers`           | `map(string)`              | Additional headers to send with the request.                                     | `{}`            | no       |
| `read_buffer_size`  | `string`                   | Size of the read buffer the gRPC client to use for reading server responses.     |                 | no       |
| `wait_for_ready`    | `boolean`                  | Waits for gRPC connection to be in the `READY` state before sending data.        | `false`         | no       |
| `write_buffer_size` | `string`                   | Size of the write buffer the gRPC client to use for writing requests.            | `"512KiB"`      | no       |

{{< docs/shared lookup="reference/components/otelcol-compression-field.md" source="alloy" version="<ALLOY_VERSION>" >}}

{{< docs/shared lookup="reference/components/otelcol-grpc-balancer-name.md" source="alloy" version="<ALLOY_VERSION>" >}}

{{< docs/shared lookup="reference/components/otelcol-grpc-authority.md" source="alloy" version="<ALLOY_VERSION>" >}}

An HTTP proxy can be configured through the following environment variables:

* `HTTPS_PROXY`
* `NO_PROXY`

The `HTTPS_PROXY` environment variable specifies a URL to use for proxying requests.
Connections to the proxy are established via [the `HTTP CONNECT` method][HTTP CONNECT].

The `NO_PROXY` environment variable is an optional list of comma-separated hostnames for which the HTTPS proxy should _not_ be used.
Each hostname can be provided as an IP address (`1.2.3.4`), an IP address in CIDR notation (`1.2.3.4/8`), a domain name (`example.com`), or `*`.
A domain name matches that domain and all subdomains.
A domain name with a leading "." (`.example.com`) matches subdomains only.
`NO_PROXY` is only read when `HTTPS_PROXY` is set.

Because `otelcol.exporter.otlp` uses gRPC, the configured proxy server must be able to handle and proxy HTTP/2 traffic.

[HTTP CONNECT]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/CONNECT

### `keepalive`

The `keepalive` block configures keepalive settings for gRPC client connections.

The following arguments are supported:

| Name                    | Type       | Description                                                                               | Default | Required |
|-------------------------|------------|-------------------------------------------------------------------------------------------|---------|----------|
| `ping_wait`             | `duration` | How often to ping the server after no activity.                                           |         | no       |
| `ping_response_timeout` | `duration` | Time to wait before closing inactive connections if the server doesn't respond to a ping. |         | no       |
| `ping_without_stream`   | `boolean`  | Send pings even if there is no active stream request.                                     |         | no       |

### `tls`

The `tls` block configures TLS settings used for the connection to the gRPC
server.

{{< docs/shared lookup="reference/components/otelcol-tls-client-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `tpm`

The `tpm` block configures retrieving the TLS `key_file` from a trusted device.

{{< docs/shared lookup="reference/components/otelcol-tls-tpm-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

{{< admonition type="note" >}}
`otelcol.exporter.otlp` uses gRPC, which doesn't allow you to send sensitive credentials like `auth` over insecure channels.
Sending sensitive credentials over insecure non-TLS connections is supported by non-gRPC exporters such as [`otelcol.exporter.otlphttp`][otelcol.exporter.otlphttp].

[otelcol.exporter.otlphttp]: ../otelcol.exporter.otlphttp/
{{< /admonition >}}

### `debug_metrics`

{{< docs/shared lookup="reference/components/otelcol-debug-metrics-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `retry_on_failure`

The `retry_on_failure` block configures how failed requests to the gRPC server are retried.

{{< docs/shared lookup="reference/components/otelcol-retry-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `sending_queue`

The `sending_queue` block configures an in-memory buffer of batches before data is sent
to the gRPC server.

{{< docs/shared lookup="reference/components/otelcol-queue-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

## Exported fields

The following fields are exported and can be referenced by other components:

| Name    | Type               | Description                                                      |
|---------|--------------------|------------------------------------------------------------------|
| `input` | `otelcol.Consumer` | A value that other components can use to send telemetry data to. |

`input` accepts `otelcol.Consumer` data for any telemetry signal (metrics, logs, or traces).

## Component health

`otelcol.exporter.otlp` is only reported as unhealthy if given an invalid configuration.

## Debug information

`otelcol.exporter.otlp` doesn't expose any component-specific debug information.

## Debug metrics

* `otelcol_exporter_queue_capacity` (gauge): Fixed capacity of the retry queue (in batches)
* `otelcol_exporter_queue_size` (gauge): Current size of the retry queue (in batches)
* `otelcol_exporter_send_failed_spans_total` (counter): Number of spans in failed attempts to send to destination.
* `otelcol_exporter_sent_spans_total` (counter): Number of spans successfully sent to destination.
* `rpc_client_duration_milliseconds` (histogram): Measures the duration of inbound RPC.
* `rpc_client_request_size_bytes` (histogram): Measures size of RPC request messages (uncompressed).
* `rpc_client_requests_per_rpc` (histogram): Measures the number of messages received per RPC. Should be 1 for all non-streaming RPCs.
* `rpc_client_response_size_bytes` (histogram): Measures size of RPC response messages (uncompressed).
* `rpc_client_responses_per_rpc` (histogram): Measures the number of messages received per RPC. Should be 1 for all non-streaming RPCs.

## Examples

The following examples show you how to create an exporter to send data to different destinations.

### Send data to a local Tempo instance

You can create an exporter that sends your data to a local Grafana Tempo instance without TLS:

```alloy
otelcol.exporter.otlp "tempo" {
    client {
        endpoint = "tempo:4317"
        tls {
            insecure             = true
            insecure_skip_verify = true
        }
    }
}
```

### Send data to a managed service

You can create an `otlp` exporter that sends your data to a managed service, for example, Grafana Cloud.
The Tempo username and Grafana Cloud API Key are injected in this example through environment variables.

```alloy
otelcol.exporter.otlp "grafana_cloud_traces" {
    client {
        endpoint = "tempo-xxx.grafana.net/tempo:443"
        auth     = otelcol.auth.basic.grafana_cloud_traces.handler
    }
}
otelcol.auth.basic "grafana_cloud_traces" {
    username = sys.env("TEMPO_USERNAME")
    password = sys.env("GRAFANA_CLOUD_API_KEY")
}
```
<!-- START GENERATED COMPATIBLE COMPONENTS -->

## Compatible components

`otelcol.exporter.otlp` has exports that can be consumed by the following components:

- Components that consume [OpenTelemetry `otelcol.Consumer`](../../../compatibility/#opentelemetry-otelcolconsumer-consumers)

{{< admonition type="note" >}}
Connecting some components may not be sensible or components may require further configuration to make the connection work correctly.
Refer to the linked documentation for more details.
{{< /admonition >}}

<!-- END GENERATED COMPATIBLE COMPONENTS -->
