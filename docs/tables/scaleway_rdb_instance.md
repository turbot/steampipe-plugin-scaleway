# Table: scaleway_rdb_instance

A database instance is composed of one or more nodes, depending on the high-availability cluster setting.

## Examples

### Basic info

```sql
select
  name,
  id,
  status,
  engine,
  region,
  project
from
  scaleway_rdb_instance;
```

### Count instances by engine type

```sql
select
  engine,
  count(id) as instance_count
from
  scaleway_rdb_instance
group by
  engine;
```

### List instances older than 90 days

```sql
select
  name,
  id,
  status,
  engine,
  extract(day from current_timestamp - created_at) as age,
  region,
  project
from
  scaleway_rdb_instance
where
  extract(day from current_timestamp - created_at) > 90;
```

### List instances with automatic backup disabled

```sql
select
  name,
  id,
  status,
  backup_schedule ->> 'disabled' as automatic_backup,
  region,
  project
from
  scaleway_rdb_instance
where
  (backup_schedule ->> 'disabled')::boolean;
```
