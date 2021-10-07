# Table: scaleway_instance_server

A Compute Instance server is a virtual server in Scaleway.

## Examples

### Basic info

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
  scaleway_instance_server;
```

### List stopped instance servers

```sql
select
  name,
  id,
  state,
  zone
  project,
  organization
from
  scaleway_instance_server
where
  state = 'stopped';
```

### List old instance servers

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
  scaleway_instance_server
where
  state = 'running'
  and extract(day from current_timestamp - creation_date) > 90;
```
