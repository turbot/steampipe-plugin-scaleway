---
title: "Steampipe Table: scaleway_rdb_database - Query Scaleway RDB Databases using SQL"
description: "Allows users to query Scaleway RDB Databases, specifically the database details, providing insights into database configurations, settings, and status."
---

# Table: scaleway_rdb_database - Query Scaleway RDB Databases using SQL

A Scaleway RDB Database is a managed relational database service that offers automated backups, high availability, and the ability to scale capacity up or down based on demand. It supports multiple database engines including MySQL, PostgreSQL, and Redis. With Scaleway RDB Database, you can focus on your application logic rather than managing database infrastructure.

## Table Usage Guide

The `scaleway_rdb_database` table provides insights into RDB Databases within Scaleway. As a Database Administrator, explore database-specific details through this table, including configurations, settings, and status. Utilize it to uncover information about databases, such as their current operational state, the engine used, and the associated instance information.

## Examples

### Basic info
Explore which Scaleway databases are managed and identify their respective sizes and locations. This can help in assessing resource allocation and optimizing database management across different projects.

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
Explore which databases are managed within your Scaleway RDB project. This query helps you to pinpoint specific locations and assess the elements within your project, providing insights into your data management and storage.

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
Explore the distribution of databases across different instances to understand how data is organized and managed within your Scaleway RDB environment. This can help optimize resource allocation and management strategies.

```sql
select
  instance_id,
  count(name)
from
  scaleway_rdb_database
group by
  instance_id;
```