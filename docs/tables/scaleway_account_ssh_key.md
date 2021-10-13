# Table: scaleway_account_ssh_key

SSH keys allow you to securely connect to your compute instances without using a password.

## Examples

### Basic info

```sql
select
  name,
  id,
  created_at,
  fingerprint
from
  scaleway_account_ssh_key;
```

### List SSH keys older than 90 days

```sql
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
