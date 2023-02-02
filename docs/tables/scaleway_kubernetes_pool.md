# Table: scaleway_kubernetes_pool

A Scaleway Kubernetes is a public cloud mamanged kubernetes available in Scaleway.

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

### List kubernetes pool with specific node type

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

### List kubernetes pool with autoscaling false

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

### List kubernetes pool with Kubernetes version is inferior to 1.24

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
