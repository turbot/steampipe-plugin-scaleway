# Table: scaleway_instance_volume

A volume is where you store your data inside your instance. It appears as a block device on Linux that you can use to create a filesystem and mount it.

## Examples

### Basic info

```sql
select
  name,
  id,
  state,
  size,
  volume_type,
  zone,
  project
from
  scaleway_instance_volume;
```

### Count of volumes by volume type

```sql
select
  volume_type,
  count(id)
from
  scaleway_instance_volume
group by
  volume_type;
```

### List of unattached volumes

```sql
select
  id,
  volume_type
from
  scaleway_instance_volume
where
  server is null;
```

### List volumes with size more than 100000000000 B (or 100 GB)

```sql
select
  name,
  id,
  state,
  size,
  volume_type,
  zone,
  project
from
  scaleway_instance_volume
where
  size > 100000000000;
```

### Find volumes attached to stopped instance servers

```sql
select
  v.name,
  v.id,
  v.state,
  v.volume_type,
  s.name as server_name,
  v.zone,
  v.project
from
  scaleway_instance_volume as v,
  scaleway_instance_server as s
where s.id = v.server ->> 'id';
```
