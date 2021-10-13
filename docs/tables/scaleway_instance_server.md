# Table: scaleway_instance_server

An instance, either virtual or physical, provides resources to run your applications on.

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

### List instance servers older than 90 days

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
