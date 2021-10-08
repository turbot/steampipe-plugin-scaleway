---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/scaleway.svg"
brand_color: "#4F0599"
display_name: "Scaleway"
short_name: "scaleway"
description: "Steampipe plugin to query servers, networks, databases and more from Scaleway account."
og_description: "Query Scaleway with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/scaleway-social-graphic.png"
---

# Scaleway + Steampipe

[Scaleway](https://www.scaleway.com) is a cloud platform, offering BareMetal and Virtual SSD Cloud Servers for any workload.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

List VPC Private Networks in your Scaleway project:

```sql
select
  name,
  id,
  zone,
  project
from
  scaleway_vpc_private_network;
```

```
+---------------------+--------------------------------------+----------+--------------------------------------+
| name                | id                                   | zone     | project                              |
+---------------------+--------------------------------------+----------+--------------------------------------+
| pvn-peaceful-diffie | 63cd190c-aef9-4ef0-8958-0a6e36f977ff | fr-par-1 | ad52df9b-0a0e-48d4-b8d2-148e95606004 |
| pvn-silly-cori      | 4691cf48-9cee-4f99-a633-e3e7c15eb5e2 | fr-par-1 | 3d4f5adb-450a-407d-a7e8-8481a6aa97d6 |
+---------------------+--------------------------------------+----------+--------------------------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/scaleway/tables)**

## Get started

### Install

Download and install the latest Scaleway plugin:

```bash
steampipe plugin install scaleway
```

### Credentials

| Item | Description |
| - | - |
| Credentials | [Get your credentials](https://console.scaleway.com/project/credentials) from [Scaleway console](https://console.scaleway.com). |
| Radius | Each connection represents a single Scaleway project. |
| Resolution | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/scaleway.spc`).<br />2. Credentials specified in environment variables e.g. `SCW_ACCESS_KEY` and `SCW_SECRET_KEY`. |
| Region Resolution | 1. Regions set for the connection via the regions argument in the config file (~/.steampipe/config/scaleway.spc).<br />2. The region specified in the `SCW_DEFAULT_REGION` environment variable. |

### Configuration

Installing the latest scaleway plugin will create a config file (`~/.steampipe/config/scaleway.spc`) with a single connection named `scaleway`:

```hcl
connection "scaleway" {
  plugin  = "scaleway"

  # You may connect to one or more regions. If `regions` is not specified,
  # Steampipe will use a single default region using:
  # The `SCW_DEFAULT_REGION` environment variable
  # regions     = ["fr-par", "nl-ams"]

  # Set the static credential with the `access_key` and `secret_key` arguments
  # Alternatively, if no creds passed in config, you may set the environment variables using
  # `SCW_ACCESS_KEY` and `SCW_SECRET_KEY` arguments
  access_key = "YOUR_ACCESS_KEY"
  secret_key = "YOUR_SECRET_ACCESS_KEY"
}
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-scaleway
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)

## Multi-Region Connections

You may also specify one or more regions with the `regions` argument:

```hcl
connection "scaleway" {
  plugin  = "scaleway"
  regions = ["fr-par", "nl-ams", "pl-waw"]
}
```

The `region` argument supports wildcards:

- All regions

  ```hcl
  connection "scaleway" {
    plugin  = "scaleway"
    regions = ["*"]
  }
  ```

Scaleway multi-region connections are common, but be aware that performance may be impacted by the number of regions and the latency to them.

## Multi-Project Connections

You may create multiple scaleway connections:

```hcl
connection "scaleway_01" {
  plugin      = "scaleway" 
  regions     = ["fr-par", "nl-ams"]
}

connection "scaleway_02" {
  plugin      = "scaleway"
  regions     = ["pl-waw"]
}
```

Each connection is implemented as a distinct [Postgres schema](https://www.postgresql.org/docs/current/ddl-schemas.html). As such, you can use qualified table names to query a specific connection:

```sql
select * from scaleway_02.scaleway_vpc_private_network;
```

Alternatively, can use an unqualified name and it will be resolved according to the [Search Path](https://steampipe.io/docs/using-steampipe/managing-connections#setting-the-search-path):

```sql
select * from scaleway_vpc_private_network;
```

You can multi-project connections by using an [**aggregator** connection](https://steampipe.io/docs/using-steampipe/managing-connections#using-aggregators).Aggregators allow you to query data from multiple connections for a plugin as if they are a single connection:

```hcl
connection "scaleway_all" {
  plugin      = "scaleway"
  type        = "aggregator"
  connections = ["scaleway_01", "scaleway_02"]
}
```

Querying tables from this connection will return results from the `scaleway_01` and `scaleway_02` connections:

```sql
select * from scaleway_all.scaleway_vpc_private_network;
```

Steampipe supports the `*` wildcard in the connection names. For example, to aggregate all the Scaleway plugin connections whose names begin with `scaleway_`:

```hcl
connection "scaleway_all" {
  type        = "aggregator"
  plugin      = "scaleway"
  connections = ["scaleway_*"]
}
```

Aggregators are powerful, but they are not infinitely scalable. Like any other steampipe connection, they query APIs and are subject to API limits and throttling.Consider as an example and aggregator that includes 3 Scaleway connections, where each connection queries 3 regions. This means you essentially run the same list API calls 9 times! When using aggregators, it is especially important to:

- Query only what you need! `select * from scaleway_object_bucket` must make a list API call in each connection, and then 9 API calls *for each bucket*, where `select name, versioning_enabled from scaleway_object_bucket` would only require a single API call per bucket.
- Consider extending the [cache TTL](https://steampipe.io/docs/reference/config-files#connection-options). The default is currently 300 seconds (5 minutes).Obviously, anytime steampipe can pull from the cache, its is faster and less impactful to the APIs. If you don't need the most up-to-date results, increase the cache TTL!

## Configuring Scaleway Credentials

### Credentials from Environment Variables

The Scaleway plugin will use the standard Scaleway environment variables to obtain credentials **only if other arguments (`access_key`, `secret_key`, `regions`) are not specified** in the connection:

```sh
export SCW_ACCESS_KEY=<YOUR_ACCESS_KEY>
export SCW_SECRET_KEY=<YOUR_SECRET_KEY>
```

```hcl
connection "scaleway" {
  plugin = "scaleway"
}
```
