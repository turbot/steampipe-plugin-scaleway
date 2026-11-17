---
title: "Steampipe Table: scaleway_iam_organization - Query Scaleway IAM Organizations using SQL"
description: "Allows users to query the Scaleway Organization configured for a connection, such as its alias and the login methods it allows."
---

# Table: scaleway_iam_organization - Query Scaleway IAM Organizations using SQL

An Organization is the top level of the Scaleway hierarchy. It owns every Project, resource and user, and it is the level billing, IAM policies and console login settings are attached to.

## Table Usage Guide

The `scaleway_iam_organization` table provides insights into the Scaleway Organization used by a connection. As a security analyst, use it to check which login methods are still allowed for the Organization, or join it with other tables like `scaleway_account_project`.

**Important Notes**
- This table requires the `organization_id` config argument to be set, unless you query it by `id`.
- The Organization API has no list endpoint, so this table returns a single row for the configured Organization.

## Examples

### Basic info
Explore the name and alias of the Organization.

```sql+postgres
select
  id,
  name,
  alias
from
  scaleway_iam_organization;
```

```sql+sqlite
select
  id,
  name,
  alias
from
  scaleway_iam_organization;
```

### Organizations still allowing password login
Identify Organizations where password login is still enabled, to check how far single sign-on has been rolled out.

```sql+postgres
select
  id,
  name,
  login_password_enabled,
  login_saml_enabled
from
  scaleway_iam_organization
where
  login_password_enabled;
```

```sql+sqlite
select
  id,
  name,
  login_password_enabled,
  login_saml_enabled
from
  scaleway_iam_organization
where
  login_password_enabled;
```

### Projects of the Organization
List the Projects the Organization owns, to see how its resources are grouped.

```sql+postgres
select
  o.name as organization_name,
  p.id as project_id,
  p.name as project_name
from
  scaleway_iam_organization as o
  join scaleway_account_project as p on p.organization_id = o.id;
```

```sql+sqlite
select
  o.name as organization_name,
  p.id as project_id,
  p.name as project_name
from
  scaleway_iam_organization as o
  join scaleway_account_project as p on p.organization_id = o.id;
```
