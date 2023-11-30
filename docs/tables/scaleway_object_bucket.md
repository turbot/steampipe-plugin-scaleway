---
title: "Steampipe Table: scaleway_object_bucket - Query Scaleway Object Storage Buckets using SQL"
description: "Allows users to query Scaleway Object Storage Buckets, providing insights into the configuration and metadata of each bucket."
---

# Table: scaleway_object_bucket - Query Scaleway Object Storage Buckets using SQL

Scaleway Object Storage is a service within Scaleway that allows you to store and retrieve any amount of data, at any time, from anywhere. It is designed to make web-scale computing easier by enabling you to store and retrieve any amount of data, at any time, from within Scaleway's computing environment. Scaleway Object Storage is a simple, scalable, and reliable object storage solution for managing data for applications, websites, backup, and restore, archive, and big data analytics.

## Table Usage Guide

The `scaleway_object_bucket` table provides insights into Object Storage Buckets within Scaleway. As a DevOps engineer, explore bucket-specific details through this table, including bucket name, region, creation date, and owner. Utilize it to uncover information about buckets, such as their ACLs, CORS configuration, and lifecycle configuration rules.

## Examples

### Basic info
Explore which of your Scaleway object storage buckets are publicly accessible. This can help in identifying potential security risks and ensuring that sensitive data is not exposed unintentionally.

```sql
select
  name,
  region,
  project,
  bucket_policy_is_public
from
  scaleway_object_bucket;
```

### List buckets with versioning disabled
Discover the segments that have versioning disabled in your Scaleway object storage. This can be useful to identify potential risks and ensure data integrity by enabling versioning.

```sql
select
  name,
  region,
  project,
  versioning_enabled
from
  scaleway_object_bucket
where
  not versioning_enabled;
```

### List buckets with no lifecycle policy
Explore which Scaleway object buckets are missing a lifecycle policy. This can be useful to identify potential areas of cost savings, as lifecycle policies can help manage storage costs by automatically archiving or deleting old data.

```sql
select
  name,
  region,
  project,
  versioning_enabled
from
  scaleway_object_bucket
where
  lifecycle_rules is null;
```