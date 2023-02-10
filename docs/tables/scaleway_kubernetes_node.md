# Table: scaleway_kubernetes_node

A Scaleway Kubernetes is a public cloud managed kubernetes available in Scaleway.

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

### List kubernetes nodes where status is not ready

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

### List kubernetes nodes created more than 90 days ago

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
