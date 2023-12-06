---
title: "Steampipe Table: scaleway_iam_user - Query Scaleway IAM Users using SQL"
description: "Allows users to query Scaleway Identity and Access Management (IAM) Users, providing insights into user permissions, roles, and associated metadata."
---

# Table: scaleway_iam_user - Query Scaleway IAM Users using SQL

Scaleway Identity and Access Management (IAM) is a service within Scaleway that helps manage access to Scaleway resources. It allows you to control who is authenticated and authorized to use resources. IAM makes it easy to manage users, security credentials, and permissions to access resources.

## Table Usage Guide

The `scaleway_iam_user` table provides insights into IAM users within Scaleway Identity and Access Management (IAM). As a DevOps engineer, explore user-specific details through this table, including permissions, roles, and associated metadata. Utilize it to uncover information about users, such as those with specific permissions, the roles assigned to each user, and the verification of user credentials.

**Important Notes**
- This table requires the `organization_id` config argument to be set.

## Examples

### Basic info
Explore the user profiles in your Scaleway IAM to understand their status and security settings. This can help identify if any users have an outdated login or if two-factor authentication is enabled, assisting in maintaining account security.

```sql
select
  email,
  created_at,
  last_login_at,
  id,
  status,
  two_factor_enabled
from
  scaleway_iam_user
```

### List all the users for whom MFA is not enabled
Explore which users have not activated Multi-Factor Authentication (MFA) to identify potential security risks and enhance user account protection measures.

```sql
select
  email,
  id,
  two_factor_enabled
from
  scaleway_iam_user
where
  not two_factor_enabled;
```

### List all the users not actived
Discover the segments that comprise users with an unknown status in the Scaleway IAM service. This allows you to pinpoint specific instances where user status may need investigation or clarification, enhancing your overall user management process.

```sql
select
  email,
  id,
  status
from
  scaleway_iam_user
where
   status = 'unknown_status';
```

### List all the users never connected
Identify users who have never logged in to your system, which can help in assessing inactive accounts and potentially freeing up resources.

```sql
select
  email,
  id,
  last_login_at
from
  scaleway_iam_user
where
   last_login_at is null;