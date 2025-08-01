---
canonical: https://grafana.com/docs/alloy/latest/reference/components/otelcol/otelcol.receiver.jaeger/
aliases:
  - ../otelcol.receiver.jaeger/ # /docs/alloy/latest/reference/otelcol.receiver.jaeger/
description: Learn about otelcol.receiver.jaeger
labels:
  stage: general-availability
  products:
    - oss
title: otelcol.receiver.jaeger
---

# `otelcol.receiver.jaeger`

`otelcol.receiver.jaeger` accepts Jaeger-formatted data over the network and forwards it to other `otelcol.*` components.

{{< admonition type="note" >}}
`otelcol.receiver.jaeger` is a wrapper over the upstream OpenTelemetry Collector [`jaeger`][] receiver.
Bug reports or feature requests will be redirected to the upstream repository, if necessary.

[`jaeger`]: https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/{{< param "OTEL_VERSION" >}}/receiver/jaegerreceiver
{{< /admonition >}}

You can specify multiple `otelcol.receiver.jaeger` components by giving them different labels.

## Usage

```alloy
otelcol.receiver.jaeger "<LABEL>" {
  protocols {
    grpc {}
    thrift_http {}
    thrift_binary {}
    thrift_compact {}
  }

  output {
    metrics = [...]
    logs    = [...]
    traces  = [...]
  }
}
```

## Arguments

The `otelcol.receiver.jaeger` component doesn't support any arguments. You can configure this component with blocks.

## Blocks

You can use the following blocks with `otelcol.receiver.jaeger`:

| Block                                                                           | Description                                                                | Required |
|---------------------------------------------------------------------------------|----------------------------------------------------------------------------|----------|
| [`output`][output]                                                              | Configures where to send received telemetry data.                          | yes      |
| [`protocols`][protocols]                                                        | Configures the protocols the component can accept traffic over.            | yes      |
| `protocols` > [`grpc`][grpc]                                                    | Configures a Jaeger gRPC server to receive traces.                         | no       |
| `protocols` > `grpc` > [`tls`][tls]                                             | Configures TLS for the gRPC server.                                        | no       |
| `protocols` > `grpc` > `tls` > [`tpm`][tpm]                                     | Configures TPM settings for the TLS key_file.                              | no       |
| `protocols` > `grpc` > [`keepalive`][keepalive]                                 | Configures keepalive settings for the configured server.                   | no       |
| `protocols` > `grpc` > `keepalive` > [`server_parameters`][server_parameters]   | Server parameters used to configure keepalive settings.                    | no       |
| `protocols` > `grpc` > `keepalive` > [`enforcement_policy`][enforcement_policy] | Enforcement policy for keepalive settings.                                 | no       |
| `protocols` > [`thrift_http`][thrift_http]                                      | Configures a Thrift HTTP server to receive traces.                         | no       |
| `protocols` > `thrift_http` > [`cors`][cors]                                    | Configures CORS for the Thrift HTTP server.                                | no       |
| `protocols` > `thrift_http` > [`tls`][tls]                                      | Configures TLS for the Thrift HTTP server.                                 | no       |
| `protocols` > `thrift_http` > `tls` > [`tpm`][tpm]                              | Configures TPM settings for the TLS key_file.                              | no       |
| `protocols` > [`thrift_binary`][thrift_binary]                                  | Configures a Thrift binary UDP server to receive traces.                   | no       |
| `protocols` > [`thrift_compact`][thrift_compact]                                | Configures a Thrift compact UDP server to receive traces.                  | no       |
| [`debug_metrics`][debug_metrics]                                                | Configures the metrics that this component generates to monitor its state. | no       |

The > symbol indicates deeper levels of nesting.
For example, `protocols` > `grpc` refers to a `grpc` block defined inside a `protocols` block.

[protocols]: #protocols
[grpc]: #grpc
[tls]: #tls
[tpm]: #tpm
[keepalive]: #keepalive
[server_parameters]: #server_parameters
[enforcement_policy]: #enforcement_policy
[thrift_http]: #thrift_http
[cors]: #cors
[thrift_binary]: #thrift_binary
[thrift_compact]: #thrift_compact
[debug_metrics]: #debug_metrics
[output]: #output

### `output`

{{< badge text="Required" >}}

{{< docs/shared lookup="reference/components/output-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `protocols`

{{< badge text="Required" >}}

The `protocols` block defines a set of protocols used to accept traces over the network.

`protocols` doesn't support any arguments and is configured fully through inner blocks.

`otelcol.receiver.jeager` requires at least one protocol block (`grpc`, `thrift_http`, `thrift_binary`, or `thrift_compact`).

### `grpc`

The `grpc` block configures a gRPC server which can accept Jaeger traces.
If the `grpc` block isn't provided, a gRPC server isn't started.

The following arguments are supported:

| Name                     | Type                       | Description                                                                  | Default           | Required |
|--------------------------|----------------------------|------------------------------------------------------------------------------|-------------------|----------|
| `auth`                   | `capsule(otelcol.Handler)` | Handler from an `otelcol.auth` component to use for authenticating requests. |                   | no       |
| `endpoint`               | `string`                   | `host:port` to listen for traffic on.                                        | `"0.0.0.0:14250"` | no       |
| `include_metadata`       | `boolean`                  | Propagate incoming connection metadata to downstream consumers.              |                   | no       |
| `max_concurrent_streams` | `number`                   | Limit the number of concurrent streaming RPC calls.                          |                   | no       |
| `max_recv_msg_size`      | `string`                   | Maximum size of messages the server will accept.                             | `"4MiB`"          | no       |
| `read_buffer_size`       | `string`                   | Size of the read buffer the gRPC server will use for reading from clients.   | `"512KiB"`        | no       |
| `transport`              | `string`                   | Transport to use for the gRPC server.                                        | `"tcp"`           | no       |
| `write_buffer_size`      | `string`                   | Size of the write buffer the gRPC server will use for writing to clients.    |                   | no       |

### `tls`

The `tls` block configures TLS settings used for a server.
If the `tls` block isn't provided, TLS isn't used for connections to the server.

{{< docs/shared lookup="reference/components/otelcol-tls-server-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `tpm`

The `tpm` block configures retrieving the TLS `key_file` from a trusted device.

{{< docs/shared lookup="reference/components/otelcol-tls-tpm-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

### `keepalive`

The `keepalive` block configures keepalive settings for connections to a gRPC server.

`keepalive` doesn't support any arguments and is configured fully through inner blocks.

### `server_parameters`

The `server_parameters` block controls keepalive and maximum age settings for gRPC servers.

The following arguments are supported:

| Name                       | Type       | Description                                                                         | Default      | Required |
|----------------------------|------------|-------------------------------------------------------------------------------------|--------------|----------|
| `max_connection_age_grace` | `duration` | Time to wait before forcibly closing connections.                                   | `"infinity"` | no       |
| `max_connection_age`       | `duration` | Maximum age for non-idle connections.                                               | `"infinity"` | no       |
| `max_connection_idle`      | `duration` | Maximum age for idle connections.                                                   | `"infinity"` | no       |
| `time`                     | `duration` | How often to ping inactive clients to check for liveness.                           | `"2h"`       | no       |
| `timeout`                  | `duration` | Time to wait before closing inactive clients that don't respond to liveness checks. | `"20s"`      | no       |

### `enforcement_policy`

The `enforcement_policy` block configures the keepalive enforcement policy for gRPC servers.
The server closes connections from clients that violate the configured policy.

The following arguments are supported:

| Name                    | Type       | Description                                                             | Default | Required |
|-------------------------|------------|-------------------------------------------------------------------------|---------|----------|
| `min_time`              | `duration` | Minimum time clients should wait before sending a keepalive ping.       | `"5m"`  | no       |
| `permit_without_stream` | `boolean`  | Allow clients to send keepalive pings when there are no active streams. | `false` | no       |

### `thrift_http`

The `thrift_http` block configures an HTTP server which can accept Thrift-formatted traces.
If the `thrift_http` block isn't specified, an HTTP server isn't started.

The following arguments are supported:

| Name                     | Type                       | Description                                                                  | Default                                                    | Required |
|--------------------------|----------------------------|------------------------------------------------------------------------------|------------------------------------------------------------|----------|
| `auth`                   | `capsule(otelcol.Handler)` | Handler from an `otelcol.auth` component to use for authenticating requests. |                                                            | no       |
| `compression_algorithms` | `list(string)`             | A list of compression algorithms the server can accept.                      | `["", "gzip", "zstd", "zlib", "snappy", "deflate", "lz4"]` | no       |
| `endpoint`               | `string`                   | `host:port` to listen for traffic on.                                        | `"0.0.0.0:14268"`                                          | no       |
| `include_metadata`       | `boolean`                  | Propagate incoming connection metadata to downstream consumers.              |                                                            | no       |
| `max_request_body_size`  | `string`                   | Maximum request body size the server will allow.                             | `"20MiB"`                                                  | no       |

### `cors`

The `cors` block configures CORS settings for an HTTP server.

The following arguments are supported:

| Name              | Type           | Description                                              | Default                | Required |
|-------------------|----------------|----------------------------------------------------------|------------------------|----------|
| `allowed_headers` | `list(string)` | Accepted headers from CORS requests.                     | `["X-Requested-With"]` | no       |
| `allowed_origins` | `list(string)` | Allowed values for the `Origin` header.                  |                        | no       |
| `max_age`         | `number`       | Configures the `Access-Control-Max-Age` response header. |                        | no       |

The `allowed_headers` specifies which headers are acceptable from a CORS request.
The following headers are always implicitly allowed:

* `Accept`
* `Accept-Language`
* `Content-Type`
* `Content-Language`

If `allowed_headers` includes `"*"`, all headers are permitted.

### `thrift_binary`

The `thrift_binary` block configures a UDP server which can accept traces formatted to the Thrift binary protocol.
If the `thrift_binary` block isn't provided, a UDP server isn't started.

The following arguments are supported:

| Name                 | Type     | Description                                                    | Default          | Required |
|----------------------|----------|----------------------------------------------------------------|------------------|----------|
| `endpoint`           | `string` | `host:port` to listen for traffic on.                          | `"0.0.0.0:6832"` | no       |
| `max_packet_size`    | `string` | Maximum UDP message size.                                      | `"65KiB"`        | no       |
| `queue_size`         | `number` | Maximum number of UDP messages that can be queued at once.     | `1000`           | no       |
| `socket_buffer_size` | `string` | Buffer to allocate for the UDP socket.                         |                  | no       |
| `workers`            | `number` | Number of workers to concurrently read from the message queue. | `10`             | no       |

### `thrift_compact`

The `thrift_compact` block configures a UDP server which can accept traces formatted to the Thrift compact protocol.
If the `thrift_compact` block isn't provided, a UDP server isn't started.

The following arguments are supported:

| Name                 | Type     | Description                                                    | Default          | Required |
|----------------------|----------|----------------------------------------------------------------|------------------|----------|
| `endpoint`           | `string` | `host:port` to listen for traffic on.                          | `"0.0.0.0:6831"` | no       |
| `max_packet_size`    | `string` | Maximum UDP message size.                                      | `"65KiB"`        | no       |
| `queue_size`         | `number` | Maximum number of UDP messages that can be queued at once.     | `1000`           | no       |
| `socket_buffer_size` | `string` | Buffer to allocate for the UDP socket.                         |                  | no       |
| `workers`            | `number` | Number of workers to concurrently read from the message queue. | `10`             | no       |

### `debug_metrics`

{{< docs/shared lookup="reference/components/otelcol-debug-metrics-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

## Exported fields

`otelcol.receiver.jaeger` doesn't export any fields.

## Component health

`otelcol.receiver.jaeger` is only reported as unhealthy if given an invalid configuration.

## Debug information

`otelcol.receiver.jaeger` doesn't expose any component-specific debug information.

## Example

This example creates a pipeline which accepts Jaeger-formatted traces and writes them to an OTLP server:

```alloy
otelcol.receiver.jaeger "default" {
  protocols {
    grpc {}
    thrift_http {}
    thrift_binary {}
    thrift_compact {}
  }

  output {
    traces = [otelcol.processor.batch.default.input]
  }
}

otelcol.processor.batch "default" {
  output {
    traces = [otelcol.exporter.otlp.default.input]
  }
}

otelcol.exporter.otlp "default" {
  client {
    endpoint = "my-otlp-server:4317"
  }
}
```

## Technical details

`otelcol.receiver.jaeger` supports [Gzip](https://en.wikipedia.org/wiki/Gzip) for compression.

## Enable authentication

You can create a `otelcol.receiver.jaeger` component that requires authentication for requests.
This is useful for limiting who can push data to the server.

{{< admonition type="note" >}}
This functionality is currently limited to the GRPC/HTTP blocks.
{{< /admonition >}}

{{< admonition type="note" >}}
Not all OpenTelemetry Collector authentication plugins support receiver authentication.
Refer to the [documentation](https://grafana.com/docs/alloy/<ALLOY_VERSION>/reference/components/otelcol/) for each `otelcol.auth.*` component to determine its compatibility.
{{< /admonition >}}

```alloy
otelcol.receiver.jaeger "default" {
  protocols {
    grpc {
      auth = otelcol.auth.basic.creds.handler
    }
    thrift_http {
      auth = otelcol.auth.basic.creds.handler
    }
  }
}

otelcol.auth.basic "creds" {
    username = sys.env("<USERNAME>")
    password = sys.env("<PASSWORD>")
}
```

<!-- START GENERATED COMPATIBLE COMPONENTS -->

## Compatible components

`otelcol.receiver.jaeger` can accept arguments from the following components:

- Components that export [OpenTelemetry `otelcol.Consumer`](../../../compatibility/#opentelemetry-otelcolconsumer-exporters)


{{< admonition type="note" >}}
Connecting some components may not be sensible or components may require further configuration to make the connection work correctly.
Refer to the linked documentation for more details.
{{< /admonition >}}

<!-- END GENERATED COMPATIBLE COMPONENTS -->
