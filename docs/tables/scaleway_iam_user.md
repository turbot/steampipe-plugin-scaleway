# Table: scaleway_iam_user

Users allow you to connect to scaleway console in your organization.

## Examples

### Basic info

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

```sql
select
  email,
  id,
  last_login_at
from
  scaleway_iam_user
where
   last_login_at is null;
