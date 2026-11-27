---
title: "Steampipe Table: scaleway_lb - Query Scaleway Load Balancers using SQL"
description: "Allows users to query Scaleway Load Balancers, including their type, status, IP addresses and the number of frontends and backends they serve."
---

# Table: scaleway_lb - Query Scaleway Load Balancers using SQL

A Scaleway Load Balancer distributes incoming traffic across a set of backend servers, such as Instances, Elastic Metal servers or any server reachable by IP. Its frontends listen on its IP addresses and forward to its backends.

Load Balancers are zonal resources.

## Table Usage Guide

The `scaleway_lb` table provides insights into the Load Balancers of your Scaleway Projects. As a cloud architect or a DevOps engineer, use it to review the size and status of your Load Balancers, to find the ones that serve no traffic, or to check the TLS settings of their frontends.

**Important Notes**
- Load Balancers are not available in every zone, so the zones without them are skipped.

## Examples

### Basic info
Explore the Load Balancers of your Projects, their type and status.

```sql+postgres
select
  name,
  id,
  status,
  type,
  zone,
  project
from
  scaleway_lb;
```

```sql+sqlite
select
  name,
  id,
  status,
  type,
  zone,
  project
from
  scaleway_lb;
```

### Get the IP addresses of each Load Balancer
Find the addresses each Load Balancer listens on, for instance to set up DNS records.

```sql+postgres
select
  l.name,
  lb_ip ->> 'ip_address' as ip_address,
  lb_ip ->> 'reverse' as reverse
from
  scaleway_lb as l,
  jsonb_array_elements(l.ip) as lb_ip;
```

```sql+sqlite
select
  l.name,
  json_extract(lb_ip.value, '$.ip_address') as ip_address,
  json_extract(lb_ip.value, '$.reverse') as reverse
from
  scaleway_lb as l,
  json_each(l.ip) as lb_ip;
```

### List Load Balancers without any backend
Identify Load Balancers that forward traffic nowhere, which may be leftovers worth deleting.

```sql+postgres
select
  name,
  id,
  type,
  frontend_count,
  created_at,
  zone
from
  scaleway_lb
where
  backend_count = 0;
```

```sql+sqlite
select
  name,
  id,
  type,
  frontend_count,
  created_at,
  zone
from
  scaleway_lb
where
  backend_count = 0;
```

### List Load Balancers accepting older TLS versions
Identify Load Balancers whose SSL compatibility level still accepts older TLS versions and ciphers.

```sql+postgres
select
  name,
  id,
  ssl_compatibility_level,
  zone
from
  scaleway_lb
where
  ssl_compatibility_level <> 'ssl_compatibility_level_modern';
```

```sql+sqlite
select
  name,
  id,
  ssl_compatibility_level,
  zone
from
  scaleway_lb
where
  ssl_compatibility_level <> 'ssl_compatibility_level_modern';
```
