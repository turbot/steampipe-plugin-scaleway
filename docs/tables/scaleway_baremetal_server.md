# Table: scaleway_baremetal_server

An bare metal provides dedicate server to run your applications on.

## Examples

### Basic info

```sql
select
  name,
  id,
  status,
  updated_at,
  created_date,
  zone,
  project,
  organization
from
  scaleway_baremetal_server;
```

### List stopped baremetal servers

```sql
select
  name,
  id,
  status,
  zone
  project,
  organization
from
  scaleway_baremetal_server
where
  status = 'stopped';
```

### List baremetal servers older than 90 days

```sql
select
  name,
  id,
  status,
   extract(day from current_timestamp - created_date) as age,
  zone,
  project,
  organization
from
  scaleway_baremetal_server
where
  state = 'running'
  and extract(day from current_timestamp - created_date) > 90;
```
