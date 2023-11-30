---
title: "Steampipe Table: scaleway_instance_ip - Query Scaleway Instance IPs using SQL"
description: "Allows users to query Instance IPs in Scaleway, specifically to obtain details about public and private IP addresses associated with instances, providing insights into network configurations and potential anomalies."
---

# Table: scaleway_instance_ip - Query Scaleway Instance IPs using SQL

Scaleway Instance IP is a resource within Scaleway that allows you to manage and configure IP addresses associated with your instances. It provides a centralized way to set up and manage IP addresses for various Scaleway resources, including Bare Metal servers, Development Instances, and more. Scaleway Instance IP helps you stay informed about the network configurations of your Scaleway resources, and take appropriate actions when necessary.

## Table Usage Guide

The `scaleway_instance_ip` table provides insights into IP addresses associated with instances within Scaleway. As a network administrator, explore IP-specific details through this table, including server IDs, public and private IP addresses, and associated metadata. Utilize it to uncover information about network configurations, such as those with specific server IDs, the relationships between servers and IP addresses, and the verification of IP status.

## Examples

### Basic info
Explore which instances in your Scaleway project are associated with specific IP addresses, to better manage your resources and ensure optimal project performance. This could be particularly useful for identifying potential bottlenecks or understanding the distribution of your resources.

```sql
select
  id,
  address,
  zone,
  project
from
  scaleway_instance_ip;
```

### List unused instance IPs
Discover the segments that contain unused IP addresses within your Scaleway instances. This can be useful for optimizing resource usage and reducing unnecessary costs.

```sql
select
  id,
  address,
  zone,
  project
from
  scaleway_instance_ip
where
  server is null
  and reverse is null;
```