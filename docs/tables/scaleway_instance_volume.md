---
title: "Steampipe Table: scaleway_instance_volume - Query Scaleway Instance Volumes using SQL"
description: "Allows users to query Scaleway Instance Volumes, providing detailed insights into the storage capabilities and configurations of Scaleway instances."
---

# Table: scaleway_instance_volume - Query Scaleway Instance Volumes using SQL

Scaleway Instance Volumes are block storage devices that you can attach to your Scaleway Instances. They offer reliable, scalable, and high-performance storage for your cloud servers. Instance Volumes can be used for primary storage of data, to provide additional storage capacity, or to increase I/O performance.

## Table Usage Guide

The `scaleway_instance_volume` table provides insights into the storage capabilities and configurations of Scaleway instances. As a system administrator or DevOps engineer, you can explore volume-specific details through this table, including size, type, and state. Utilize it to monitor storage usage, verify configurations, and ensure optimal storage performance for your Scaleway instances.

## Examples

### Basic info
Explore which instances are active within your Scaleway project, along with their respective sizes and types. This can help you manage resources and identify areas for potential optimization or scaling.

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
Analyze the distribution of volume types in your Scaleway instance to better understand your storage utilization. This could potentially help optimize storage resources by identifying which volume types are most commonly used.

```sql
select
  volume_type,
  count(id)
from
  scaleway_instance_volume
group by
  volume_type;
```

### List unattached volumes
Discover the segments that consist of unused storage volumes within your Scaleway instances. This can aid in optimizing storage utilization and reducing unnecessary costs.

```sql
select
  id,
  volume_type
from
  scaleway_instance_volume
where
  server is null;
```

### List volumes with size more than 10 GB (10000000000 Bytes)
Identify instances where your Scaleway volumes exceed 10 GB to help manage your storage resources more effectively.

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
  size > 10000000000;
```

### Find volumes attached to stopped instance servers
Determine the areas in which storage volumes are attached to servers that are not currently active. This is useful for optimizing resource usage and managing costs, as unused volumes may be unnecessarily incurring charges.

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
where
  s.id = v.server ->> 'id';
```