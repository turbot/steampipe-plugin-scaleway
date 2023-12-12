---
title: "Steampipe Table: scaleway_instance_server - Query Scaleway Instances using SQL"
description: "Allows users to query Scaleway Instances, specifically the server details, providing insights into server configurations, statuses, and other related information."
---

# Table: scaleway_instance_server - Query Scaleway Instances using SQL

Scaleway Instances is a service offered by Scaleway that allows you to deploy virtual instances in seconds. These instances are scalable, powerful, and reliable cloud servers designed for developers. They are equipped with SSD disks, high-end Intel CPUs, and come with a variety of OS and Apps.

## Table Usage Guide

The `scaleway_instance_server` table provides insights into the instances within Scaleway. As a system administrator or a developer, explore server-specific details through this table, including server configurations, statuses, and other related information. Utilize it to uncover information about servers, such as their commercial type, creation date, dynamic IP required status, and more.

## Examples

### Basic info
Explore which Scaleway servers are currently active, when they were created, and where they are located. This is useful for gaining insights into resource allocation and server management across different projects and organizations.

```sql+postgres
select
  name,
  id,
  state,
  creation_date,
  zone,
  project,
  organization
from
  scaleway_instance_server;
```

```sql+sqlite
select
  name,
  id,
  state,
  creation_date,
  zone,
  project,
  organization
from
  scaleway_instance_server;
```

### List stopped instance servers
Identify instances where Scaleway servers are in a 'stopped' state. This is useful to manage resources and maintain operational efficiency by pinpointing idle servers.

```sql+postgres
select
  name,
  id,
  state,
  zone
  project,
  organization
from
  scaleway_instance_server
where
  state = 'stopped';
```

```sql+sqlite
select
  name,
  id,
  state,
  zone,
  project,
  organization
from
  scaleway_instance_server
where
  state = 'stopped';
```

### List instance servers older than 90 days
Determine the areas in which instance servers have been running for more than 90 days. This is useful for identifying potential areas for resource optimization and cost-saving by assessing long-running servers.

```sql+postgres
select
  name,
  id,
  state,
   extract(day from current_timestamp - creation_date) as age,
  zone,
  project,
  organization
from
  scaleway_instance_server
where
  state = 'running'
  and extract(day from current_timestamp - creation_date) > 90;
```

```sql+sqlite
select
  name,
  id,
  state,
  julianday('now') - julianday(creation_date) as age,
  zone,
  project,
  organization
from
  scaleway_instance_server
where
  state = 'running'
  and julianday('now') - julianday(creation_date) > 90;
```