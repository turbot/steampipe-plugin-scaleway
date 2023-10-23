# Table: scaleway_iam_api_key

API keys allow you to securely connect to scaleway console in your organization.

This table requires the `organization_id` config argument to be set.

## Examples

### Basic info

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
