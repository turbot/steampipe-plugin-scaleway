---
title: "Steampipe Table: scaleway_account_ssh_key - Query Scaleway Account SSH Keys using SQL"
description: "Allows users to query Scaleway Account SSH Keys, specifically the details of each SSH key associated with the account, providing insights into key management and security."
---

# Table: scaleway_account_ssh_key - Query Scaleway Account SSH Keys using SQL

Scaleway is a cloud computing company that provides a range of scalable computing resources, including bare metal servers, development platforms, and storage options. Among these resources, Scaleway Account SSH Keys are used to ensure secure shell access to instances. These keys are critical for maintaining the security and integrity of user data and operations within the Scaleway environment.

## Table Usage Guide

The `scaleway_account_ssh_key` table provides insights into the SSH keys associated with a Scaleway account. As a DevOps engineer, explore key-specific details through this table, including key fingerprints, creation dates, and associated metadata. Utilize it to uncover information about keys, such as their usage across instances, their creation and modification history, and the overall security posture of your Scaleway account.

**Important Notes**
- Unlike `scaleway_iam_api_key` and `scaleway_iam_user`, this table does not require the `organization_id` config argument. When it is set, the listing is scoped to that organization; when it is not, the IAM API scopes the request to the organization owning the credentials.

## Examples

### Basic info
Explore the creation dates and unique identifiers of your Scaleway account's SSH keys. This allows you to manage and track the keys more effectively.

```sql+postgres
select
  name,
  id,
  created_at,
  fingerprint
from
  scaleway_account_ssh_key;
```

```sql+sqlite
select
  name,
  id,
  created_at,
  fingerprint
from
  scaleway_account_ssh_key;
```

### List SSH keys older than 90 days
Discover the SSH keys that have been in use for more than 90 days. This is useful in understanding potential security risks associated with outdated keys and aids in maintaining optimal security practices.

```sql+postgres
select
  name,
  id,
  created_at,
  fingerprint,
  extract(day from current_timestamp - created_at) as age
from
  scaleway_account_ssh_key
where
  extract(day from current_timestamp - created_at) > 90;
```

```sql+sqlite
select
  name,
  id,
  created_at,
  fingerprint,
  cast((julianday('now') - julianday(created_at)) as integer) as age
from
  scaleway_account_ssh_key
where
  cast((julianday('now') - julianday(created_at)) as integer) > 90;
```

### List disabled SSH keys
Identify SSH keys that have been disabled. A disabled key is kept on the account but cannot be used to access instances, so reviewing them helps confirm that revoked access has actually been retired.

```sql+postgres
select
  name,
  id,
  created_at,
  fingerprint
from
  scaleway_account_ssh_key
where
  disabled;
```

```sql+sqlite
select
  name,
  id,
  created_at,
  fingerprint
from
  scaleway_account_ssh_key
where
  disabled = 1;
```