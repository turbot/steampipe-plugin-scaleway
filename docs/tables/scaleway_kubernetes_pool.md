# Table: scaleway_kubernetes_pool

A Scaleway Kubernetes pool is a group of Scaleway Instances, organized by type. It represents the computing power of the cluster and contains the Kubernetes nodes, on which the containers run.

## Examples

### Basic info

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
