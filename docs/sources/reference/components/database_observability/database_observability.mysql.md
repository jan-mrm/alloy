---
canonical: https://grafana.com/docs/alloy/latest/reference/components/database_observability.mysql/
description: Learn about database_observability.mysql
title: database_observability.mysql
labels:
  stage: experimental
  products:
    - oss
---

# `database_observability.mysql`

{{< docs/shared lookup="stability/experimental.md" source="alloy" version="<ALLOY_VERSION>" >}}

## Usage

```alloy
database_observability.mysql "<LABEL>" {
  data_source_name = <DATA_SOURCE_NAME>
  forward_to       = [<LOKI_RECEIVERS>]
}
```

## Arguments

You can use the following arguments with `database_observability.mysql`:

| Name                               | Type                 | Description                                                                                    | Default | Required |
|------------------------------------|----------------------|------------------------------------------------------------------------------------------------|---------|----------|
| `data_source_name`                 | `secret`             | [Data Source Name][] for the MySQL server to connect to.                                       |         | yes      |
| `forward_to`                       | `list(LogsReceiver)` | Where to forward log entries after processing.                                                 |         | yes      |
| `collect_interval`                 | `duration`           | How frequently to collect information from database.                                           | `"1m"`  | no       |
| `disable_collectors`               | `list(string)`       | A list of collectors to disable from the default set.                                          |         | no       |
| `disable_query_redaction`          | `bool`               | Collect unredacted SQL query text including parameters.                                        | `false` | no       |
| `enable_collectors`                | `list(string)`       | A list of collectors to enable on top of the default set.                                      |         | no       |
| `explain_plan_collect_interval`    | `duration`           | How frequently to collect explain plan information from database.                              | `"1m"`  | no       |
| `explain_plan_per_collect_ratio`   | `float`              | Ratio of explain plan queries to collect per collect interval.                                 | `1.0`   | no       |
| `explain_plan_initial_lookback`    | `duration`           | How far back to look for explain plan queries on the first collection interval.                | `"24h"` | no       |
| `locks_collect_interval`           | `duration`           | How frequently to collect locks information from database.                                     | `"30s"` | no       |
| `locks_threshold`                  | `duration`           | Threshold for locks to be considered slow. If a lock exceeds this duration, it will be logged. | `"1s"`  | no       |
| `setup_consumers_collect_interval` | `duration`           | How frequently to collect `performance_schema.setup_consumers` information from the database.    | `"1h"`  | no       |
| `allow_update_performance_schema_settings` | `boolean`     | Whether to allow updates to `performance_schema` settings in any collector. | `false` | no |
| `query_sample_auto_enable_setup_consumers` | `boolean`     | Whether to allow the `query_sample` collector to enable some specific `performance_schema.setup_consumers` settings. | `false` | no |

The following collectors are configurable:

| Name              | Description                                              | Enabled by default |
|-------------------|----------------------------------------------------------|--------------------|
| `query_tables`    | Collect query table information.                         | yes                |
| `schema_table`    | Collect schemas and tables from `information_schema`.    | yes                |
| `query_sample`    | Collect query samples.                                   | yes                |
| `setup_consumers` | Collect enabled `performance_schema.setup_consumers`.    | yes                |
| `locks`           | Collect queries that are waiting/blocking other queries. | no                 |
| `explain_plan`    | Collect explain plan information.                        | no                 |

## Blocks

The `database_observability.mysql` component doesn't support any blocks. You can configure this component with arguments.

## Example

```alloy
database_observability.mysql "orders_db" {
  data_source_name = "user:pass@tcp(mysql:3306)/"
  forward_to = [loki.write.logs_service.receiver]
}

prometheus.scrape "orders_db" {
  targets = database_observability.mysql.orders_db.targets
  honor_labels = true // required to keep job and instance labels
  forward_to = [prometheus.remote_write.metrics_service.receiver]
}

prometheus.remote_write "metrics_service" {
  endpoint {
    url = sys.env("<GRAFANA_CLOUD_HOSTED_METRICS_URL>")
    basic_auth {
      username = sys.env("<GRAFANA_CLOUD_HOSTED_METRICS_ID>")
      password = sys.env("<GRAFANA_CLOUD_RW_API_KEY>")
    }
  }
}

loki.write "logs_service" {
  endpoint {
    url = sys.env("<GRAFANA_CLOUD_HOSTED_LOGS_URL>")
    basic_auth {
      username = sys.env("<GRAFANA_CLOUD_HOSTED_LOGS_ID>")
      password = sys.env("<GRAFANA_CLOUD_RW_API_KEY>")
    }
  }
}
```

Replace the following:

* _`<GRAFANA_CLOUD_HOSTED_METRICS_URL>`_: The URL for your Grafana Cloud hosted metrics.
* _`<GRAFANA_CLOUD_HOSTED_METRICS_ID>`_: The user ID for your Grafana Cloud hosted metrics.
* _`<GRAFANA_CLOUD_RW_API_KEY>`_: Your Grafana Cloud API key.
* _`<GRAFANA_CLOUD_HOSTED_LOGS_URL>`_: The URL for your Grafana Cloud hosted logs.
* _`<GRAFANA_CLOUD_HOSTED_LOGS_ID>`_: The user ID for your Grafana Cloud hosted logs.

[Data Source Name]: https://github.com/go-sql-driver/mysql#dsn-data-source-name

<!-- START GENERATED COMPATIBLE COMPONENTS -->

## Compatible components

`database_observability.mysql` can accept arguments from the following components:

- Components that export [Loki `LogsReceiver`](../../../compatibility/#loki-logsreceiver-exporters)

`database_observability.mysql` has exports that can be consumed by the following components:

- Components that consume [Targets](../../../compatibility/#targets-consumers)

{{< admonition type="note" >}}
Connecting some components may not be sensible or components may require further configuration to make the connection work correctly.
Refer to the linked documentation for more details.
{{< /admonition >}}

<!-- END GENERATED COMPATIBLE COMPONENTS -->
