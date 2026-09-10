---
title: "Steampipe Table: scaleway_vpc - Query Scaleway VPCs using SQL"
description: "Allows users to query Scaleway VPCs, providing insights into the regional networks their private networks belong to."
---

# Table: scaleway_vpc - Query Scaleway VPCs using SQL

A Scaleway VPC is a regional network that groups the private networks of a project. Every project has a default VPC per region, and resources attached to private networks of the same VPC can communicate with each other over private IP addresses.

## Table Usage Guide

The `scaleway_vpc` table provides insights into the VPCs of your Scaleway projects. As a network engineer, you can explore each VPC through this table, including whether it is the default one of its project and how many private networks it holds. Use it to review your network layout and to join private networks back to the VPC that contains them.

## Examples

### Basic info
Explore the VPCs available across your Scaleway projects to get an overview of your network layout.

```sql+postgres
select
  name,
  id,
  is_default,
  private_network_count,
  created_at,
  region,
  project
from
  scaleway_vpc;
```

```sql+sqlite
select
  name,
  id,
  is_default,
  private_network_count,
  created_at,
  region,
  project
from
  scaleway_vpc;
```

### List default VPCs
Identify the default VPC of each project and region, which is where resources are placed when no VPC is specified.

```sql+postgres
select
  name,
  id,
  region,
  project
from
  scaleway_vpc
where
  is_default;
```

```sql+sqlite
select
  name,
  id,
  region,
  project
from
  scaleway_vpc
where
  is_default = 1;
```

### List VPCs without any private network
Discover VPCs that hold no private network, which are candidates for clean up.

```sql+postgres
select
  name,
  id,
  region,
  project
from
  scaleway_vpc
where
  private_network_count = 0;
```

```sql+sqlite
select
  name,
  id,
  region,
  project
from
  scaleway_vpc
where
  private_network_count = 0;
```

### List the private networks of each VPC
Determine the private networks contained in each VPC to understand how a project's network is split up.

```sql+postgres
select
  v.name as vpc_name,
  v.id as vpc_id,
  v.region,
  n.name as private_network_name,
  n.id as private_network_id
from
  scaleway_vpc as v
  left join scaleway_vpc_private_network as n on n.vpc_id = v.id
  and n.region = v.region
order by
  v.name,
  n.name;
```

```sql+sqlite
select
  v.name as vpc_name,
  v.id as vpc_id,
  v.region,
  n.name as private_network_name,
  n.id as private_network_id
from
  scaleway_vpc as v
  left join scaleway_vpc_private_network as n on n.vpc_id = v.id
  and n.region = v.region
order by
  v.name,
  n.name;
```
