---
title: "Steampipe Table: scaleway_kubernetes_cluster - Query Scaleway Kubernetes Clusters using SQL"
description: "Allows users to query Scaleway Kubernetes Clusters, specifically providing information about the clusters' configuration, status, and metadata."
---

# Table: scaleway_kubernetes_cluster - Query Scaleway Kubernetes Clusters using SQL

A Scaleway Kubernetes Cluster is a managed service in the Scaleway ecosystem that allows users to deploy, manage, and scale containerized applications using Kubernetes, an open-source container orchestration platform. It provides a highly available and scalable infrastructure for deploying and running applications and services. The clusters can be customized according to the user's requirements, including the choice of the Kubernetes version, the size and type of worker nodes, and the geographical location of the cluster.

## Table Usage Guide

The `scaleway_kubernetes_cluster` table provides insights into Kubernetes Clusters within Scaleway. As a DevOps engineer, explore cluster-specific details through this table, including version, status, and associated metadata. Utilize it to uncover information about clusters, such as their configuration, the geographical location, and the current status of the clusters.

## Examples

### Basic info
Explore which Kubernetes clusters are currently active within your Scaleway account. This can help you understand the status and version of each cluster, which is useful for maintenance and upgrade planning.

```sql+postgres
select
  name,
  description,
  type,
  cluster_url,
  id,
  status,
  version
from
  scaleway_kubernetes_cluster;
```

```sql+sqlite
select
  name,
  description,
  type,
  cluster_url,
  id,
  status,
  version
from
  scaleway_kubernetes_cluster;
```

### List Kapsule clusters
Discover the segments that are utilizing Kapsule clusters within your Scaleway Kubernetes environment. This query is beneficial for gaining insights into the operational status and details of these specific clusters.

```sql+postgres
select
  name,
  description,
  type,
  cluster_url,
  id,
  status
from
  scaleway_kubernetes_cluster
where
  type = 'kapsule';
```

```sql+sqlite
select
  name,
  description,
  type,
  cluster_url,
  id,
  status
from
  scaleway_kubernetes_cluster
where
  type = 'kapsule';
```

### List Kosmos clusters
Determine the areas in which multicloud Kosmos clusters are being used. This query can be useful to understand the spread and utilization of multicloud resources, providing valuable insight for resource management and planning.

```sql+postgres
select
  name,
  description,
  type,
  cluster_url,
  id,
  status
from
  scaleway_kubernetes_cluster
where
  type = 'multicloud';
```

```sql+sqlite
select
  name,
  description,
  type,
  cluster_url,
  id,
  status
from
  scaleway_kubernetes_cluster
where
  type = 'multicloud';
```

### List clusters with Kubernetes version inferior to 1.24
Identify any clusters operating on a Kubernetes version less than 1.24. This is useful for pinpointing clusters that may need to be updated to maintain compatibility and security standards.

```sql+postgres
select
  name,
  description,
  type,
  cluster_url,
  id,
  status
from
  scaleway_kubernetes_cluster
where
  version < '1.24';
```

```sql+sqlite
select
  name,
  description,
  type,
  cluster_url,
  id,
  status
from
  scaleway_kubernetes_cluster
where
  version < '1.24';
```

### List clusters with upgrades available
Discover the segments that have upgrades available in your Kubernetes clusters on Scaleway. This can help in maintaining up-to-date environments, improving security and performance.

```sql+postgres
select
  name,
  type,
  id,
  version,
  auto_upgrade,
  upgrade_available
from
  scaleway_kubernetes_cluster
where
  upgrade_available is true;
```

```sql+sqlite
select
  name,
  type,
  id,
  version,
  auto_upgrade,
  upgrade_available
from
  scaleway_kubernetes_cluster
where
  upgrade_available = 1;
```

### List clusters with auto-upgrade enabled
Determine the areas in which clusters have the auto-upgrade feature enabled to ensure that they are always running the latest version and are not vulnerable to outdated software issues.

```sql+postgres
select
  name,
  type,
  id,
  version,
  auto_upgrade,
  upgrade_available
from
  scaleway_kubernetes_cluster
where
  auto_upgrade @> '{"enabled":true}';
```

```sql+sqlite
select
  name,
  type,
  id,
  version,
  auto_upgrade,
  upgrade_available
from
  scaleway_kubernetes_cluster
where
  json_extract(auto_upgrade, '$.enabled') = 1;
```