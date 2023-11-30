---
title: "Steampipe Table: scaleway_kubernetes_pool - Query Scaleway Kubernetes Pools using SQL"
description: "Allows users to query Scaleway Kubernetes Pools, specifically providing insights into the configuration, status, and metadata of each pool."
---

# Table: scaleway_kubernetes_pool - Query Scaleway Kubernetes Pools using SQL

Scaleway Kubernetes Pools are a feature within Scaleway's Kubernetes Service that allows for the grouping of nodes within a Kubernetes cluster. These pools provide a way to manage the distribution and scaling of workloads across different nodes. Kubernetes Pools are essential for maintaining the performance and reliability of applications running on a Kubernetes cluster.

## Table Usage Guide

The `scaleway_kubernetes_pool` table offers insights into the configuration and status of Kubernetes Pools within Scaleway's Kubernetes Service. As a DevOps engineer or system administrator, you can explore details about each pool, such as its size, autoscaling settings, and associated metadata. Utilize this table to monitor the state of your Kubernetes Pools, identify any potential scaling issues, and ensure optimal distribution of workloads across your Kubernetes cluster.

## Examples

### Basic info
Explore the status and details of your Kubernetes clusters in Scaleway. This allows you to assess the health and version of your clusters, helping with maintenance and troubleshooting.

```sql
select
  name,
  node_type,
  cluster_id,
  id,
  status,
  created_at,
  version
from
  scaleway_kubernetes_pool;
```

### List kubernetes pools with a specific node type
Determine the areas in which Kubernetes pools are utilizing a specific node type within the Scaleway platform. This can be useful for resource optimization and understanding the distribution of node types across your Kubernetes clusters.

```sql
select
  name,
  node_type,
  cluster_id,
  id,
  status,
  created_at,
  version
from
  scaleway_kubernetes_pool
where
  node_type = 'play2_nano';
```

### List kubernetes pools with auto-scaling disabled
Identify Kubernetes pools where auto-scaling is turned off. This can be useful in managing system resources and preventing unexpected scaling actions.

```sql
select
  name,
  node_type,
  cluster_id,
  id,
  status,
  autoscaling,
  version
from
  scaleway_kubernetes_pool
where
  autoscaling is false;
```

### List kubernetes pools with Kubernetes version inferior to 1.24
Determine the areas in which Kubernetes pools are operating on versions older than 1.24. This can help in identifying pools that may need to be updated for security or feature improvements.

```sql
select
  name,
  node_type,
  cluster_id,
  id,
  status,
  version
from
  scaleway_kubernetes_pool
where
  version < '1.24';
```