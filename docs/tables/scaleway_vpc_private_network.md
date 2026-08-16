---
title: "Steampipe Table: scaleway_vpc_private_network - Query Scaleway VPC Private Networks using SQL"
description: "Allows users to query Scaleway VPC Private Networks, providing insights into the network's configuration and associated metadata."
---

# Table: scaleway_vpc_private_network - Query Scaleway VPC Private Networks using SQL

A Scaleway VPC Private Network is a resource within Scaleway that allows users to create isolated network environments. These networks are used to interconnect instances and other resources, providing a secure and efficient way to manage network traffic within a Scaleway environment. With this resource, users can control IP addressing, subnetting, and routing to provide secure, private communication between instances.

Private networks are regional resources: each one belongs to a VPC and carries the subnets attached to it.

## Table Usage Guide

The `scaleway_vpc_private_network` table provides insights into private networks within Scaleway's Virtual Private Cloud (VPC). As a network engineer, you can explore network-specific details through this table, including the VPC the network belongs to, its subnets, and related metadata. Use it to uncover information about private networks, such as their IP ranges and DHCP configuration.

## Examples

### Basic info
Explore which private networks have been created within your Scaleway VPCs. This is useful for keeping track of your network configurations and identifying any potential issues or areas for improvement.

```sql+postgres
select
  name,
  id,
  vpc_id,
  created_at,
  region,
  project
from
  scaleway_vpc_private_network;
```

```sql+sqlite
select
  name,
  id,
  vpc_id,
  created_at,
  region,
  project
from
  scaleway_vpc_private_network;
```

### List subnets attached to each private network
Determine the IP ranges in use across your private networks, which helps to plan new allocations and to spot overlapping ranges.

```sql+postgres
select
  n.name,
  n.id,
  n.region,
  s ->> 'subnet' as subnet
from
  scaleway_vpc_private_network as n,
  jsonb_array_elements(n.subnets) as s;
```

```sql+sqlite
select
  n.name,
  n.id,
  n.region,
  json_extract(s.value, '$.subnet') as subnet
from
  scaleway_vpc_private_network as n,
  json_each(n.subnets) as s;
```

### List private networks of the default VPC
Identify the private networks created in the default VPC of each project, which are the ones resources land in when no VPC is specified.

```sql+postgres
select
  n.name,
  n.id,
  n.region,
  v.name as vpc_name
from
  scaleway_vpc_private_network as n
  join scaleway_vpc as v on v.id = n.vpc_id
  and v.region = n.region
where
  v.is_default;
```

```sql+sqlite
select
  n.name,
  n.id,
  n.region,
  v.name as vpc_name
from
  scaleway_vpc_private_network as n
  join scaleway_vpc as v on v.id = n.vpc_id
  and v.region = n.region
where
  v.is_default = 1;
```

### List private networks with managed DHCP enabled
Assess which private networks hand out IP addresses themselves, since instances attached to the others need their addresses configured statically.

```sql+postgres
select
  name,
  id,
  vpc_id,
  region
from
  scaleway_vpc_private_network
where
  dhcp_enabled;
```

```sql+sqlite
select
  name,
  id,
  vpc_id,
  region
from
  scaleway_vpc_private_network
where
  dhcp_enabled = 1;
```
