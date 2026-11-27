---
title: "Steampipe Table: scaleway_rdb_instance - Query Scaleway RDB Instances using SQL"
description: "Allows users to query Scaleway RDB Instances, specifically to retrieve information about each RDB instance like its name, status, region, and other related data."
---

# Table: scaleway_rdb_instance - Query Scaleway RDB Instances using SQL

Scaleway RDB Instances are a part of Scaleway's managed database services. These services provide a scalable, reliable, and easy to use database solution. They are fully managed by Scaleway, ensuring automatic updates, backups, and scalability without any manual intervention.

## Table Usage Guide

The `scaleway_rdb_instance` table provides insights into RDB instances within Scaleway's managed database services. As a database administrator or a DevOps engineer, you can explore instance-specific details through this table, including instance status, region, and other related data. Utilize it to manage and monitor your RDB instances effectively, ensuring optimal performance and security.

## Examples

### Basic info
Explore which Scaleway RDB instances are currently active or inactive, their associated engines, and the regions they're located in, all within the context of a specific project. This is beneficial for maintaining an overview of your database instances and their status, especially in larger projects.

```sql+postgres
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

```sql+sqlite
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
Analyze the distribution of instances based on their engine types to understand the usage patterns and preferences in your Scaleway RDB instances. This can help in making informed decisions for resource allocation and optimization.

```sql+postgres
select
  engine,
  count(id) as instance_count
from
  scaleway_rdb_instance
group by
  engine;
```

```sql+sqlite
select
  engine,
  count(id) as instance_count
from
  scaleway_rdb_instance
group by
  engine;
```

### List instances older than 90 days
Determine the areas in which Scaleway RDB instances have been running for over 90 days. This can be useful for identifying potential cost-saving opportunities by shutting down or resizing long-running instances.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  status,
  engine,
  julianday('now') - julianday(created_at) as age,
  region,
  project
from
  scaleway_rdb_instance
where
  julianday('now') - julianday(created_at) > 90;
```

### List instances with automatic backup disabled
Identify instances where automatic backup has been disabled, allowing you to assess risk and take necessary action to ensure data safety across different projects and regions. This is particularly useful in maintaining data integrity and preventing potential data loss.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  status,
  json_extract(backup_schedule, '$.disabled') as automatic_backup,
  region,
  project
from
  scaleway_rdb_instance
where
  json_extract(backup_schedule, '$.disabled') = 'true';
```
### List instances not encrypted at rest
Identify instances whose volume is not encrypted at rest, to check compliance with data protection requirements.

```sql+postgres
select
  name,
  id,
  engine,
  volume ->> 'type' as volume_type,
  region,
  project
from
  scaleway_rdb_instance
where
  not coalesce((encryption ->> 'enabled')::boolean, false);
```

```sql+sqlite
select
  name,
  id,
  engine,
  json_extract(volume, '$.type') as volume_type,
  region,
  project
from
  scaleway_rdb_instance
where
  coalesce(json_extract(encryption, '$.enabled'), 0) = 0;
```

### Get the backup schedule and storage of each instance
Review how often each instance is backed up, how long its backups are kept, and where they are stored.

```sql+postgres
select
  name,
  (backup_schedule ->> 'frequency')::int as backup_frequency_hours,
  (backup_schedule ->> 'retention')::int as backup_retention_days,
  backup_same_region,
  volume ->> 'type' as volume_type,
  (volume ->> 'size')::bigint as volume_size_bytes
from
  scaleway_rdb_instance;
```

```sql+sqlite
select
  name,
  json_extract(backup_schedule, '$.frequency') as backup_frequency_hours,
  json_extract(backup_schedule, '$.retention') as backup_retention_days,
  backup_same_region,
  json_extract(volume, '$.type') as volume_type,
  json_extract(volume, '$.size') as volume_size_bytes
from
  scaleway_rdb_instance;
```

### List instances reachable from any IP
Identify instances whose ACL rules let any address reach their public endpoint.

```sql+postgres
select
  name,
  id,
  rule ->> 'ip' as ip,
  rule ->> 'description' as description,
  region
from
  scaleway_rdb_instance,
  jsonb_array_elements(acl_rules) as rule
where
  rule ->> 'ip' = '0.0.0.0/0';
```

```sql+sqlite
select
  name,
  id,
  json_extract(rule.value, '$.ip') as ip,
  json_extract(rule.value, '$.description') as description,
  region
from
  scaleway_rdb_instance,
  json_each(acl_rules) as rule
where
  json_extract(rule.value, '$.ip') = '0.0.0.0/0';
```
