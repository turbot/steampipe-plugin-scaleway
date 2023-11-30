---
title: "Steampipe Table: scaleway_registry_namespace - Query Scaleway Registry Namespaces using SQL"
description: "Allows users to query Scaleway Registry Namespaces, providing insights into the details of each namespace, including its ID, name, endpoint, and creation date."
---

# Table: scaleway_registry_namespace - Query Scaleway Registry Namespaces using SQL

Scaleway Registry Namespace is a resource that belongs to Scaleway's Container Registry service. It allows users to create isolated spaces to host their container images. Each namespace provides a unique endpoint where users can push and pull images.

## Table Usage Guide

The `scaleway_registry_namespace` table provides insights into the registry namespaces within Scaleway's Container Registry service. As a DevOps engineer or system administrator, you can explore the details of each namespace through this table, including its unique identifier, name, endpoint, and the date it was created. This table is useful for managing and tracking your container images, and for ensuring the organization and security of your container registry.

## Examples

### Basic info
Explore the status and creation date of your Scaleway registry namespaces, which can help you track and manage your resources more effectively. This is particularly useful for maintaining organization and project information across multiple regions.

```sql
select
  name,
  id,
  status,
  created_at,
  region,
  project,
  organization
from
  scaleway_registry_namespace;
```

### List public registry namespaces
Explore which registry namespaces are publicly accessible. This can help in understanding the level of data exposure and potential security risks.

```sql
select
  name,
  id,
  status,
  created_at,
  region,
  project,
  organization
from
  scaleway_registry_namespace
where
  is_public = true;
```