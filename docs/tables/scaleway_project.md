# Table: scaleway_project

A Project is a grouping of Scaleway resources. Each Scaleway Organization comes with a default Project, and you can create new Projects if necessary. Projects are cross-region, meaning resources located in different regions can be grouped in one single Project. When grouping resources into different Projects, you can use IAM to define custom access rights for each Project.

This table requires an Organization ID to be configured in the scaleway.spc file.

## Examples

### Basic info

```sql
select
  name,
  project_id,
  created_at,
  organization
from
  scaleway_project;
```
