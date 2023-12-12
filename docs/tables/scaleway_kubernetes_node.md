---
title: "Steampipe Table: scaleway_kubernetes_node - Query Scaleway Kubernetes Nodes using SQL"
description: "Allows users to query Scaleway Kubernetes Nodes, providing insights into the node details, including their status, versions, and associated metadata."
---

# Table: scaleway_kubernetes_node - Query Scaleway Kubernetes Nodes using SQL

Scaleway Kubernetes Nodes are the worker machines in a Kubernetes cluster that run containerized applications. Each node contains the necessary services to run Pods (the smallest and simplest unit in the Kubernetes object model that you create or deploy), including the container runtime, kubelet, and the kube-proxy. Nodes can be a virtual or physical machine, depending on the cluster.

## Table Usage Guide

The `scaleway_kubernetes_node` table provides insights into Kubernetes Nodes within Scaleway. As a DevOps engineer, explore node-specific details through this table, including status, versions, and associated metadata. Utilize it to uncover information about nodes, such as their health status, the Kubernetes version they are running, and their associated roles and labels.

## Examples

### Basic info
Analyze the status and creation date of your Scaleway Kubernetes nodes to understand their current operational state and longevity. This can be useful in assessing the overall health and maintenance needs of your Kubernetes infrastructure.

```sql+postgres
select
  name,
  cluster_id,
  id,
  status,
  created_at
from
  scaleway_kubernetes_node;
```

```sql+sqlite
select
  name,
  cluster_id,
  id,
  status,
  created_at
from
  scaleway_kubernetes_node;
```

### List kubernetes nodes where status is not ready
Identify Kubernetes nodes that are not in a 'ready' status. This query is useful in pinpointing potential issues within your Kubernetes cluster that may need attention or troubleshooting.

```sql+postgres
select
  name,
  cluster_id,
  id,
  status,
  error_message,
  created_at
from
  scaleway_kubernetes_node
where
  status <> 'ready';
```

```sql+sqlite
select
  name,
  cluster_id,
  id,
  status,
  error_message,
  created_at
from
  scaleway_kubernetes_node
where
  status <> 'ready';
```

### List kubernetes nodes with ipv6 public
Analyze the settings to understand the status and creation date of Kubernetes nodes on Scaleway that are utilizing IPv6 public IP addresses. This is useful for maintaining network configurations and ensuring optimal performance.

```sql+postgres
select
  name,
  cluster_id,
  id,
  status,
  public_ip_v6,
  created_at
from
  scaleway_kubernetes_node
where
  public_ip_v6 != '<nil>';
```

```sql+sqlite
select
  name,
  cluster_id,
  id,
  status,
  public_ip_v6,
  created_at
from
  scaleway_kubernetes_node
where
  public_ip_v6 != '<nil>';
```

### List kubernetes nodes created more than 90 days ago
Identify instances where Kubernetes nodes have been active for a prolonged period of time, specifically more than 90 days. This can be useful in managing resources and ensuring optimal performance within your system.

```sql+postgres
select
  name,
  cluster_id,
  id,
  status,
  updated_at,
  created_at
from
  scaleway_kubernetes_node
where
  created_at <= now() - interval '90' day;
```

```sql+sqlite
select
  name,
  cluster_id,
  id,
  status,
  updated_at,
  created_at
from
  scaleway_kubernetes_node
where
  created_at <= datetime('now','-90 day');
```