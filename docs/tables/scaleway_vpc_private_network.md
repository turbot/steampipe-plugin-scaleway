# Table: scaleway_vpc_private_network

A VPC private network allows interconnecting your instances in an isolated and private network.

## Examples

### Basic info

```sql
select
  name,
  id,
  created_at,
  zone,
  project
from
  scaleway_vpc_private_network;
```
