---
title: "Steampipe Table: scaleway_redis_cluster - Query Scaleway Managed Database for Redis clusters using SQL"
description: "Allows users to query Scaleway Managed Database for Redis clusters, including their version, size, TLS setting and endpoints."
---

# Table: scaleway_redis_cluster - Query Scaleway Managed Database for Redis clusters using SQL

Scaleway Managed Database for Redis runs fully managed Redis clusters. A cluster is made of one or more nodes of a given node type, and is reachable through endpoints on Private Networks or on the public network.

Clusters are zonal resources.

## Table Usage Guide

The `scaleway_redis_cluster` table provides insights into the Managed Database for Redis clusters of your Scaleway Projects. As a database administrator or a security engineer, use it to review the versions and sizes of your clusters, to check which of them have TLS enabled, or to see which networks can reach them.

**Important Notes**
- Managed Database for Redis is not available in every zone, so the zones without it are skipped.

## Examples

### Basic info
Explore the clusters of your Projects, their version and size.

```sql+postgres
select
  name,
  id,
  status,
  version,
  node_type,
  cluster_size,
  zone,
  project
from
  scaleway_redis_cluster;
```

```sql+sqlite
select
  name,
  id,
  status,
  version,
  node_type,
  cluster_size,
  zone,
  project
from
  scaleway_redis_cluster;
```

### List clusters without TLS
Identify clusters whose traffic is not encrypted in transit.

```sql+postgres
select
  name,
  id,
  version,
  zone
from
  scaleway_redis_cluster
where
  not tls_enabled;
```

```sql+sqlite
select
  name,
  id,
  version,
  zone
from
  scaleway_redis_cluster
where
  tls_enabled = 0;
```

### List the endpoints of each cluster
Determine where each cluster can be reached from, on a Private Network or on the public network.

```sql+postgres
select
  c.name,
  e ->> 'port' as port,
  e -> 'private_network' ->> 'id' as private_network_id,
  e -> 'public_network' is not null as is_public,
  e -> 'ips' as ips
from
  scaleway_redis_cluster as c,
  jsonb_array_elements(c.endpoints) as e;
```

```sql+sqlite
select
  c.name,
  json_extract(e.value, '$.port') as port,
  json_extract(e.value, '$.private_network.id') as private_network_id,
  json_extract(e.value, '$.public_network') is not null as is_public,
  json_extract(e.value, '$.ips') as ips
from
  scaleway_redis_cluster as c,
  json_each(c.endpoints) as e;
```

### List clusters with an upgrade available
Find clusters that can be upgraded to a newer Redis version.

```sql+postgres
select
  name,
  id,
  version,
  upgradable_versions
from
  scaleway_redis_cluster
where
  jsonb_array_length(upgradable_versions) > 0;
```

```sql+sqlite
select
  name,
  id,
  version,
  upgradable_versions
from
  scaleway_redis_cluster
where
  json_array_length(upgradable_versions) > 0;
```
