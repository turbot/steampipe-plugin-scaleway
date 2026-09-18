# Table: scaleway_iam_organization

An Organization is the top level of the Scaleway hierarchy. It owns every Project, resource and user, and it is the level billing, IAM policies and console login settings are attached to.

The Organization API exposes no list endpoint, and an API key belongs to exactly one Organization, so this table returns a single row: the Organization configured in the scaleway.spc file.

This table requires an Organization ID to be configured in the scaleway.spc file.

## Examples

### Basic info

```sql
select
  id,
  name,
  alias
from
  scaleway_iam_organization;
```

### Organizations still allowing password login

```sql
select
  id,
  name,
  login_password_enabled,
  login_saml_enabled
from
  scaleway_iam_organization
where
  login_password_enabled;
```

### Projects of the Organization

```sql
select
  o.name as organization_name,
  p.id as project_id,
  p.name as project_name
from
  scaleway_iam_organization as o
  join scaleway_account_project as p on p.organization_id = o.id;
```
