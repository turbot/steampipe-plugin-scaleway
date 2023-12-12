---
title: "Steampipe Table: scaleway_baremetal_server - Query Scaleway Baremetal Servers using SQL"
description: "Allows users to query Scaleway Baremetal Servers, specifically providing information about the server's status, location, and configuration."
---

# Table: scaleway_baremetal_server - Query Scaleway Baremetal Servers using SQL

A Scaleway Baremetal Server is a physical server provided by Scaleway, a cloud computing company. These servers offer high-performance capabilities and full control over the hardware, making them ideal for compute-intensive workloads. They can be customized according to the user's needs, with a variety of CPUs, memory, and storage options available.

## Table Usage Guide

The `scaleway_baremetal_server` table offers insights into the Baremetal Servers within Scaleway. As a system administrator, you can explore server-specific details through this table, including server status, location, and configuration. Use this table to monitor the performance and health of your servers, verify their configurations, and ensure they are optimally located for their intended workloads.

## Examples

### Basic info
Explore the status and timeline of your Scaleway baremetal servers across different projects and organizations. This aids in understanding the distribution and upkeep of your servers, assisting in resource management and operational efficiency.

```sql+postgres
select
  name,
  id,
  status,
  updated_at,
  created_date,
  zone,
  project,
  organization
from
  scaleway_baremetal_server;
```

```sql+sqlite
select
  name,
  id,
  status,
  updated_at,
  created_date,
  zone,
  project,
  organization
from
  scaleway_baremetal_server;
```

### List stopped bare metal servers
Discover the segments that include inactive bare metal servers across various zones, projects, and organizations. This can be useful for assessing resource utilization and identifying potential areas for cost savings.

```sql+postgres
select
  name,
  id,
  status,
  zone
  project,
  organization
from
  scaleway_baremetal_server
where
  status = 'stopped';
```

```sql+sqlite
select
  name,
  id,
  status,
  zone,
  project,
  organization
from
  scaleway_baremetal_server
where
  status = 'stopped';
```

### List bare metal servers older than 90 days
Determine the areas in which bare metal servers have been running for over 90 days. This is useful for assessing long-term usage and identifying potential areas for resource optimization.

```sql+postgres
select
  name,
  id,
  status,
  extract(day from current_timestamp - created_date) as age,
  zone,
  project,
  organization
from
  scaleway_baremetal_server
where
  state = 'running'
  and extract(day from current_timestamp - created_date) > 90;
```

```sql+sqlite
select
  name,
  id,
  status,
  julianday('now') - julianday(created_date) as age,
  zone,
  project,
  organization
from
  scaleway_baremetal_server
where
  state = 'running'
  and julianday('now') - julianday(created_date) > 90;
```