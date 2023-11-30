---
title: "Steampipe Table: scaleway_vpc_private_network - Query Scaleway VPC Private Networks using SQL"
description: "Allows users to query Scaleway VPC Private Networks, providing insights into the network's configuration and associated metadata."
---

# Table: scaleway_vpc_private_network - Query Scaleway VPC Private Networks using SQL

A Scaleway VPC Private Network is a resource within Scaleway that allows users to create isolated network environments. These networks are used to interconnect instances and other resources, providing a secure and efficient way to manage network traffic within a Scaleway environment. With this resource, users can control IP addressing, subnetting, and routing to provide secure, private communication between instances.

## Table Usage Guide

The `scaleway_vpc_private_network` table provides insights into private networks within Scaleway's Virtual Private Cloud (VPC). As a network engineer, you can explore network-specific details through this table, including network configuration, associated instances, and related metadata. Use it to uncover information about private networks, such as their IP range, default configuration, and associated instances.

## Examples

### Basic info
Explore which private networks have been created within your Scaleway VPC. This is useful for keeping track of your network configurations and identifying any potential issues or areas for improvement.

```sql
select
  name,
  id,
  created_at,
  zone,
  project
from
  scaleway_vpc_private_network;
```