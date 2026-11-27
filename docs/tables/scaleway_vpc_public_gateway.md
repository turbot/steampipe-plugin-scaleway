---
title: "Steampipe Table: scaleway_vpc_public_gateway - Query Scaleway Public Gateways using SQL"
description: "Allows users to query Scaleway Public Gateways, including their public IP, the Private Networks they serve and their SSH bastion settings."
---

# Table: scaleway_vpc_public_gateway - Query Scaleway Public Gateways using SQL

A Scaleway Public Gateway sits at the border of one or more Private Networks. It gives the resources on them access to the internet through NAT, and can expose an SSH bastion to reach them from outside.

Public Gateways are zonal resources.

## Table Usage Guide

The `scaleway_vpc_public_gateway` table provides insights into the Public Gateways of your Scaleway Projects. As a network or security engineer, use it to review which Private Networks reach the internet through which gateway, and which gateways expose an SSH bastion and to whom.

**Important Notes**
- Public Gateways are not available in every zone, so the zones without them are skipped.

## Examples

### Basic info
Explore the Public Gateways of your Projects, their type and their public IP.

```sql+postgres
select
  name,
  id,
  status,
  type,
  ipv4 ->> 'address' as public_ip,
  zone,
  project
from
  scaleway_vpc_public_gateway;
```

```sql+sqlite
select
  name,
  id,
  status,
  type,
  json_extract(ipv4, '$.address') as public_ip,
  zone,
  project
from
  scaleway_vpc_public_gateway;
```

### List the Private Networks served by each Public Gateway
Determine which Private Networks reach the internet through which Public Gateway, and whether the gateway masquerades their traffic.

```sql+postgres
select
  g.name as gateway_name,
  gn ->> 'private_network_id' as private_network_id,
  (gn ->> 'masquerade_enabled')::boolean as masquerade_enabled,
  (gn ->> 'push_default_route')::boolean as push_default_route
from
  scaleway_vpc_public_gateway as g,
  jsonb_array_elements(g.gateway_networks) as gn;
```

```sql+sqlite
select
  g.name as gateway_name,
  json_extract(gn.value, '$.private_network_id') as private_network_id,
  json_extract(gn.value, '$.masquerade_enabled') as masquerade_enabled,
  json_extract(gn.value, '$.push_default_route') as push_default_route
from
  scaleway_vpc_public_gateway as g,
  json_each(g.gateway_networks) as gn;
```

### List Public Gateways exposing an SSH bastion to any IP
Identify gateways whose SSH bastion accepts connections from anywhere, which widens the attack surface of the Private Networks behind them.

```sql+postgres
select
  name,
  id,
  bastion_port,
  bastion_allowed_ips,
  zone
from
  scaleway_vpc_public_gateway
where
  bastion_enabled
  and bastion_allowed_ips ? '0.0.0.0/0';
```

```sql+sqlite
select
  name,
  id,
  bastion_port,
  bastion_allowed_ips,
  zone
from
  scaleway_vpc_public_gateway
where
  bastion_enabled = 1
  and exists (
    select 1 from json_each(bastion_allowed_ips) where value = '0.0.0.0/0'
  );
```
