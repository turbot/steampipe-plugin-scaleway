---
title: "Steampipe Table: scaleway_iam_api_key - Query Scaleway IAM API Keys using SQL"
description: "Allows users to query Scaleway IAM API Keys, specifically providing details about each key, including ID, state, creation and expiration dates, and associated user information."
---

# Table: scaleway_iam_api_key - Query Scaleway IAM API Keys using SQL

Scaleway IAM API Keys are unique identifiers associated with a Scaleway account. These keys are used to authenticate and authorize actions in the Scaleway API. They can be created, deleted, and managed through the Scaleway console or API.

## Table Usage Guide

The `scaleway_iam_api_key` table provides insights into IAM API keys within Scaleway Identity and Access Management (IAM). As a security analyst, explore key-specific details through this table, including key states, creation and expiration dates, and associated user information. Utilize it to monitor and manage the lifecycle of API keys, ensuring they are rotated regularly and expired keys are deleted, enhancing the security posture of your Scaleway environment.

## Examples

### Basic info
Explore which Scaleway IAM API keys have been created, when they were made, and their expiration dates. This allows for efficient management of API keys, ensuring none are expired or unused.

```sql
select
  access_key,
  created_at,
  user_id,
  expires_at,
  default_project_id
from
  scaleway_iam_api_key
```

### List API keys older than 90 days
Determine the areas in which API keys are older than 90 days to manage and update them as necessary. This helps in maintaining security and ensuring all keys are up-to-date.

```sql
select
  access_key,
  created_at,
  user_id,
  expires_at,
  default_project_id,
  extract(day from current_timestamp - created_at) as age
from
  scaleway_iam_api_key
where
  extract(day from current_timestamp - created_at) > 90;
```