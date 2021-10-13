# Table: scaleway_instance_snapshot

Snapshots contain the data of a specific volume at a particular point in time. The data can include the instance's operating system, configuration information or files stored on the volume.

## Examples

### Basic info

```sql
select
  name,
  id,
  state,
  size,
  zone,
  project
from
  scaleway_instance_snapshot;
```

### List snapshots older than 90 days

```sql
select
  name,
  id,
  state,
  extract(day from current_timestamp - creation_date) as age,
  size,
  zone,
  project
from
  scaleway_instance_snapshot
where
  extract(day from current_timestamp - creation_date) > 90;
```

### List large snapshots (> 100GB or 100000000000 Bytes)

```sql
select
  name,
  id,
  state,
  size,
  zone,
  project
from
  scaleway_instance_snapshot
where
  size > 100000000000;
```
