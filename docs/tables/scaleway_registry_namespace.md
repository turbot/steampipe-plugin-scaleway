# Table: scaleway_registry_namespace

Namespaces allow you to manage your Container Registry in a simple, clear and human-readable way.

## Examples

### Basic info

```sql
select
  name,
  id,
  status,
  created_at,
  region,
  project,
  organization
from
  scaleway_registry_namespace;
```

### List public registry namespaces

```sql
select
  name,
  id,
  status,
  created_at,
  region,
  project,
  organization
from
  scaleway_registry_namespace
where
  is_public = true;
```

