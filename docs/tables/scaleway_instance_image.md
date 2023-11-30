---
title: "Steampipe Table: scaleway_instance_image - Query Scaleway Instance Images using SQL"
description: "Allows users to query Scaleway Instance Images, providing detailed information about each image available for use with Scaleway's Instances."
---

# Table: scaleway_instance_image - Query Scaleway Instance Images using SQL

Scaleway Instance Images are pre-configured operating system environments that can be used as a base to create new instances in Scaleway's cloud computing service. These images include a variety of Linux distributions, developer tools, and container technologies. They allow users to quickly launch and scale applications in the cloud.

## Table Usage Guide

The `scaleway_instance_image` table provides insights into Instance Images within Scaleway's cloud computing service. As a cloud architect or developer, explore image-specific details through this table, including their IDs, names, and the architectures they support. Utilize it to uncover information about images, such as their creation dates, modification dates, and the public visibility status.

## Examples

### Basic info
Explore the details of your Scaleway instances, such as their names, IDs, states, creation dates, and associated projects and organizations. This helps you to effectively manage and monitor your Scaleway resources.

```sql
select
  name,
  id,
  state,
  creation_date,
  zone,
  project,
  organization
from
  scaleway_instance_image;
```

### List custom (user-defined) images
Discover the segments that are utilizing custom images in your cloud infrastructure. This can assist in identifying areas for optimization and potential security risks.

```sql
select
  name,
  id,
  state,
  creation_date,
  zone,
  project,
  organization
from
  scaleway_instance_image
where
  not public;
```

### List images older than 90 days
Determine the areas in which images from Scaleway instances have been stored for more than 90 days. This can be useful in identifying outdated or potentially unnecessary data, enabling more efficient resource management and cost savings.

```sql
select
  name,
  id,
  state,
  extract(day from current_timestamp - creation_date) as age,
  zone,
  project,
  organization
from
  scaleway_instance_image
where
  extract(day from current_timestamp - creation_date) > 90;
```