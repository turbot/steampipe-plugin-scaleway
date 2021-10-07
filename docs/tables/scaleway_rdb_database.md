# Table: scaleway_rdb_database

A RDB database is a logical databases on your instance.

## Examples

### Basic info

```sql
select
  name,
  instance_id,
  size,
  managed,
  region,
  project
from
  scaleway_rdb_database;
```

### List managed databases

```sql
select
  name,
  instance_id,
  size,
  managed,
  region,
  project
from
  scaleway_rdb_database
where
  managed;
```

### Get count of databases by instance

```sql
select
  instance_id,
  count(name)
from
  scaleway_rdb_database
group by
  instance_id;
```
