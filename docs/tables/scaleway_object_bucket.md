# Table: scaleway_object_bucket

A Scaleway Object bucket is a public cloud storage resource available in Scaleway, an object storage offering.

## Examples

### Basic info

```sql
select
  name,
  region,
  project,
  bucket_policy_is_public
from
  scaleway_object_bucket;
```

### List buckets with versioning disabled

```sql
select
  name,
  region,
  project,
  versioning_enabled
from
  scaleway_object_bucket
where
  not versioning_enabled;
```

### List buckets with no lifecycle policy

```sql
select
  name,
  region,
  project,
  versioning_enabled
from
  scaleway_object_bucket
where
  lifecycle_rules is null;
```
