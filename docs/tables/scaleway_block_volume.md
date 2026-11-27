---
title: "Steampipe Table: scaleway_block_volume - Query Scaleway Block Storage volumes using SQL"
description: "Allows users to query Scaleway Block Storage volumes, including their size, performance class, encryption key and the resources they are attached to."
---

# Table: scaleway_block_volume - Query Scaleway Block Storage volumes using SQL

Scaleway Block Storage provides network-attached volumes for Instances. Each volume has a size and a performance class, can be encrypted with a Key Manager key, and records the resources it is attached to.

Block Storage volumes are zonal resources.

## Table Usage Guide

The `scaleway_block_volume` table provides insights into the Block Storage volumes of your Scaleway Projects. As a cloud engineer or a security analyst, use it to review the size and performance of the volumes of your Instances, to check which volumes are encrypted with a customer-managed key, or to find volumes attached to nothing.

**Important Notes**
- The Instance API only returns the ID and type of a Block Storage volume attached to a server, so this table is where its size and performance come from: join it with `scaleway_instance_server` on the volume ID.

## Examples

### Basic info
Explore the Block Storage volumes of your Projects, their type and size.

```sql+postgres
select
  name,
  id,
  type,
  size,
  status,
  zone,
  project
from
  scaleway_block_volume;
```

```sql+sqlite
select
  name,
  id,
  type,
  size,
  status,
  zone,
  project
from
  scaleway_block_volume;
```

### Get the Block Storage volumes of each Instance
Find the size and performance of the volumes attached to each Instance.

```sql+postgres
select
  s.name as server_name,
  v.name as volume_name,
  v.type,
  v.size,
  v.perf_iops
from
  scaleway_instance_server as s,
  jsonb_each(s.volumes) as sv
  join scaleway_block_volume as v on v.id = sv.value ->> 'id'
where
  sv.value ->> 'volume_type' = 'sbs_volume';
```

```sql+sqlite
select
  s.name as server_name,
  v.name as volume_name,
  v.type,
  v.size,
  v.perf_iops
from
  scaleway_instance_server as s,
  json_each(s.volumes) as sv
  join scaleway_block_volume as v on v.id = json_extract(sv.value, '$.id')
where
  json_extract(sv.value, '$.volume_type') = 'sbs_volume';
```

### List volumes not encrypted with a customer-managed key
Identify volumes that are not encrypted with a Key Manager key of your own.

```sql+postgres
select
  name,
  id,
  type,
  zone,
  project
from
  scaleway_block_volume
where
  kms_key_id is null;
```

```sql+sqlite
select
  name,
  id,
  type,
  zone,
  project
from
  scaleway_block_volume
where
  kms_key_id is null;
```

### List volumes attached to nothing
Identify volumes that are not attached to any resource, which may be leftovers worth deleting.

```sql+postgres
select
  name,
  id,
  size,
  last_detached_at,
  zone
from
  scaleway_block_volume
where
  jsonb_array_length(coalesce(attached_resources, '[]'::jsonb)) = 0;
```

```sql+sqlite
select
  name,
  id,
  size,
  last_detached_at,
  zone
from
  scaleway_block_volume
where
  json_array_length(coalesce(attached_resources, '[]')) = 0;
```
