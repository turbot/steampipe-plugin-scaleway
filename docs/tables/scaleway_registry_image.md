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
  tags,
  size
from
  scaleway_registry_image;
```

### List images updated in last 10 days for a repository

```sql
select
  name,
  id,
  status,
  created_at,
  updated_at
  tags,
  size
from
  scaleway_registry_image
where
  updated_at >= now() - interval '10' day
```


### List images with a public visibility

```sql
select
  name,
  id,
  status,
  created_at,
  updated_at
  tags,
  visibility
from
  scaleway_registry_image
where
  visibility = 'public'
```
