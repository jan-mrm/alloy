---
canonical: https://grafana.com/docs/alloy/latest/reference/components/otelcol/otelcol.auth.headers/
aliases:
  - ../otelcol.auth.headers/ # /docs/alloy/latest/reference/components/otelcol.auth.headers/
description: Learn about otelcol.auth.headers
labels:
  stage: general-availability
  products:
    - oss
title: otelcol.auth.headers
---

# `otelcol.auth.headers`

`otelcol.auth.headers` exposes a `handler` that other `otelcol` components can use to authenticate requests using custom headers.

This component only supports client authentication.

{{< admonition type="note" >}}
`otelcol.auth.headers` is a wrapper over the upstream OpenTelemetry Collector [`headerssetter`][] extension.
Bug reports or feature requests will be redirected to the upstream repository, if necessary.

[`headerssetter`]: https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/{{< param "OTEL_VERSION" >}}/extension/headerssetterextension
{{< /admonition >}}

You can specify multiple `otelcol.auth.headers` components by giving them different labels.

## Usage

```alloy
otelcol.auth.headers "<LABEL>" {
  header {
    key   = "<HEADER_NAME>"
    value = "<HEADER_VALUE>"
  }
}
```

## Arguments

The `otelcol.auth.headers` component doesn't support any arguments. You can configure this component with blocks.

## Blocks

You can use the following blocks with `otelcol.auth.headers`:

| Block                            | Description                                                                | Required |
|----------------------------------|----------------------------------------------------------------------------|----------|
| [`header`][header]               | Custom header to attach to requests.                                       | yes      |
| [`debug_metrics`][debug_metrics] | Configures the metrics that this component generates to monitor its state. | no       |

[header]: #header
[debug_metrics]: #debug_metrics

### `header`

{{< badge text="Required" >}}

The `header` block defines a custom header to attach to requests.
It's valid to provide multiple `header` blocks to set more than one header.

| Name             | Type                 | Description                                                  | Default    | Required |
|------------------|----------------------|--------------------------------------------------------------|------------|----------|
| `key`            | `string`             | Name of the header to set.                                   |            | yes      |
| `action`         | `string`             | An action to perform on the header.                          | `"upsert"` | no       |
| `from_attribute` | `string`             | Authentication attribute name used to retrieve header value. |            | no       |
| `from_context`   | `string`             | Metadata name used to retrieve header value.                 |            | no       |
| `value`          | `string` or `secret` | Value of the header.                                         |            | no       |

The supported values for `action` are:

* `insert`: Inserts the new header if it doesn't exist.
* `update`: Updates the header value if it exists.
* `upsert`: Inserts a header if it doesn't exist and updates the header if it exists.
* `delete`: Deletes the header.

Exactly one of `value`, `from_context`, or `from_attribute` must be provided for each `header` block.

The `value` attribute sets the value of the header directly.
Alternatively, you can use `from_context` to dynamically retrieve the header value from request metadata, or you can use `from_attribute` to dynamically retrieve the header value from request authentication metadata.

For `from_context` to work, other components in the pipeline also need to be configured appropriately:

* If an `otelcol.processor.batch` is present in the pipeline, it must be configured to preserve client metadata.
  Do this by adding the value that `from_context` needs to the `metadata_keys` of the batch processor.
* `otelcol` receivers must be configured with `include_metadata` set to `true` so that metadata keys are available to the pipeline.

`from_attribute` metadata can't, at this time, be preserved through an `otelcol.processor.batch` component, and is only provided from the `otelcol.auth.basic` extension.

### `debug_metrics`

{{< docs/shared lookup="reference/components/otelcol-debug-metrics-block.md" source="alloy" version="<ALLOY_VERSION>" >}}

## Exported fields

The following fields are exported and can be referenced by other components:

| Name      | Type                       | Description                                                     |
|-----------|----------------------------|-----------------------------------------------------------------|
| `handler` | `capsule(otelcol.Handler)` | A value that other components can use to authenticate requests. |

## Component health

`otelcol.auth.headers` is only reported as unhealthy if given an invalid configuration.

## Debug information

`otelcol.auth.headers` doesn't expose any component-specific debug information.

## Example

This example configures [`otelcol.exporter.otlp`][otelcol.exporter.otlp] to use custom headers:

```alloy
otelcol.receiver.otlp "default" {
  http {
    include_metadata = true
  }
  grpc {
    include_metadata = true
  }

  output {
    metrics = [otelcol.processor.batch.default.input]
    logs    = [otelcol.processor.batch.default.input]
    traces  = [otelcol.processor.batch.default.input]
  }
}

otelcol.processor.batch "default" {
  // Preserve the tenant_id metadata.
  metadata_keys = ["tenant_id"]

  output {
    metrics = [otelcol.exporter.otlp.production.input]
    logs    = [otelcol.exporter.otlp.production.input]
    traces  = [otelcol.exporter.otlp.production.input]
  }
}

otelcol.auth.headers "creds" {
  header {
    key          = "X-Scope-OrgID"
    from_context = "tenant_id"
  }

  header {
    key   = "User-ID"
    value = "user_id"
  }
}

otelcol.exporter.otlp "production" {
  client {
    endpoint = sys.env("<OTLP_SERVER_ENDPOINT>")
    auth     = otelcol.auth.headers.creds.handler
  }
}
```

[otelcol.exporter.otlp]: ../otelcol.exporter.otlp/
