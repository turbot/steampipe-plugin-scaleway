# Table: scaleway_kubernetes_cluster

A Scaleway Kubernetes is a public cloud mamanged kubernetes available in Scaleway.

## Examples

### Basic info

```sql
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

### List Kapsule cluster

```sql
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

### List Kosmos cluster

```sql
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

```sql
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

```sql
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

### List clusters with auto-upgrade enabled

```sql
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
