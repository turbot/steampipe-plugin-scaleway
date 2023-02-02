# Table: scaleway_kubernetes_pool

A Scaleway Kubernetes is a public cloud mamanged kubernetes available in Scaleway.

## Examples

### Basic info

```sql
select
  name,
  cluster_id,
  id,
  status,
  created_at
from
  scaleway_kubernetes_node;
```

### List kubernetes node with status is not ready

```sql
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

### List kubernetes node with ipv6 public

```sql
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

### List kubernetes node createad more than 90 days ago

```sql
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
  created_at <= now() - interval '90' day
```
