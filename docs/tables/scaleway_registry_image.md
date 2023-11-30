---
title: "Steampipe Table: scaleway_registry_image - Query Scaleway Registry Images using SQL"
description: "Allows users to query Scaleway Registry Images, providing details about the image tags, visibility, status, and more."
---

# Table: scaleway_registry_image - Query Scaleway Registry Images using SQL

A Scaleway Registry Image is a resource in the Scaleway Container Registry, which is a fully-managed service to store and manage your Docker images. It allows you to store your images centrally for fast and reliable deployments in your applications. Each image in the registry includes a unique tag, status, visibility settings, and other metadata.

## Table Usage Guide

The `scaleway_registry_image` table provides insights into registry images within Scaleway's Container Registry. As a DevOps engineer, explore image-specific details through this table, including the image tags, visibility settings, status, and other metadata. Utilize it to manage and streamline your Docker image deployments and ensure the security and efficiency of your containerized applications.


## Examples

### Basic info
Explore which Scaleway registry images are active and when they were created to manage storage effectively. This helps in understanding the overall usage and aids in resource optimization.

```sql
select
  name,
  id,
  status,
  created_at,
  tags,
  size
from
  scaleway_registry_image;
```

### List images updated in last 10 days for a repository
Determine the images in a repository that have been updated recently, allowing you to stay informed about the latest changes and developments within your project.

```sql
select
  name,
  id,
  status,
  created_at,
  updated_at
  tags,
  size
from
  scaleway_registry_image
where
  updated_at >= now() - interval '10' day
```


### List images with a public visibility
Discover the segments that have images with public visibility, allowing you to assess potential security risks and manage access controls more effectively. This is useful for maintaining data privacy and ensuring only appropriate images are publicly accessible.

```sql
select
  name,
  id,
  status,
  created_at,
  updated_at
  tags,
  visibility
from
  scaleway_registry_image
where
  visibility = 'public'
```