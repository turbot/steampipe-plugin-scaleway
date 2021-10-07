# Table: scaleway_instance_security_group

A security group is a set of firewall rules on a set of instances. Security groups enable to create rules that either drop or allow incoming traffic from certain ports of your instances.

## Examples

### Basic info

```sql
select
  name,
  id,
  creation_date,
  project_default,
  zone,
  project
from
  scaleway_instance_security_group;
```

### List default security group in a project

```sql
select
  name,
  id,
  creation_date,
  project_default,
  zone,
  project
from
  scaleway_instance_security_group
where
  project_default;
```
