# Table: scaleway_instance_image

Images are backups of your instances. You can reuse that image to restore your data or create a series of instances with a predefined configuration.

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
  scaleway_instance_image;
```

### List of custom (user-defined) images defined

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
  scaleway_instance_image
where
  not public;
```

### List images older than 90 days

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
  scaleway_instance_image
where
  extract(day from current_timestamp - creation_date) > 90;
```
