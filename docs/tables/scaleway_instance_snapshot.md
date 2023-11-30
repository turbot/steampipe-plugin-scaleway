---
title: "Steampipe Table: scaleway_instance_snapshot - Query Scaleway Instance Snapshots using SQL"
description: "Allows users to query Scaleway Instance Snapshots, primarily revealing details about the snapshot's state, size, creation date, and associated instance."
---

# Table: scaleway_instance_snapshot - Query Scaleway Instance Snapshots using SQL

Scaleway Instance Snapshots are a resource within Scaleway's cloud services that allow users to create a point-in-time copy of their instances. This is particularly useful for creating backups, migrating data, or testing changes without affecting the original instance. These snapshots contain all the information necessary to restore your instance (system settings, applications, and data) from the moment the snapshot was taken.

## Table Usage Guide

The `scaleway_instance_snapshot` table provides insights into Instance Snapshots within Scaleway's cloud services. As a system administrator or DevOps engineer, explore snapshot-specific details through this table, including snapshot state, size, creation date, and the associated instance. Utilize it to manage and understand your instance backups, verify snapshot details, and ensure the integrity and safety of your data.

## Examples

### Basic info
Explore which Scaleway instance snapshots are currently active, by assessing their state and size. This can help manage resources and plan projects more efficiently.

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
Assess the elements within your Scaleway instances by identifying snapshots that have been stored for more than 90 days. This can be useful for managing storage and ensuring efficient use of resources.

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
Discover the segments that contain large snapshots, specifically those exceeding 100GB. This can be beneficial in managing storage space and optimizing resource allocation within your project.

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