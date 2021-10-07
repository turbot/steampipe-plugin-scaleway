# Table: scaleway_instance_ip

A flexible IP address is an IP address which you hold independently of any server. You can attach it to any of your servers and do live migration of the IP address between your servers.

## Examples

### Basic info

```sql
select
  id,
  address,
  zone,
  project
from
  scaleway_instance_ip;
```

### List unused instance IPs

```sql
select
  id,
  address,
  zone,
  project
from
  scaleway_instance_ip
where
  server is null
  and reverse is null;
```
