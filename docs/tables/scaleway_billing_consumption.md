---
title: "Steampipe Table: scaleway_billing_consumtion - Query Scaleway Billing comsumption using SQL"
description: "Allows users to query Scaleway Billing Consumption, specifically providing information about the billing."
---

# Table: scaleway_billing_consumtion - Query Scaleway Billing comsumption using SQL

Scaleway Billing Consumption tracks the costs of various scaleway products such as object_storage, storage, serverless, etc.

**Important Notes**
- This table requires an Organization ID to be configured in the `scaleway.spc` file.

## Table Usage Guide

The `scaleway_billing_consumtion` table offers insights into the Billing Consumption within Scaleway.

## Examples

### Basic info
Explore the billing for all scaleway services.

```sql+postgres
select
  category_name,
  product_name,
  project_id,
  resource_name,
  value
from
  scaleway_billing_consumption;
```

```sql+sqlite
select
  category_name,
  product_name,
  project_id,
  resource_name,
  value
from
  scaleway_billing_consumption;
```

### List consumption by category
Explore costs by product category of scaleway

```sql+postgres
select 
  category_name, 
  sum((value -> 'units')::float) AS units 
from 
  scaleway_billing_consumption 
group by 
  category_name;
```

```sql+sqlite
select 
  category_name, 
  sum((value -> 'units')::float) AS units 
from 
  scaleway_billing_consumption 
group by 
  category_name;
```

### List consumption by Project ID
Explore costs by product category of scaleway

```sql+postgres
select 
  project_id, 
  sum((value -> 'units')::float) AS units 
from 
  scaleway_billing_consumption 
group by 
  project_id;
```

```sql+sqlite
select 
  project_id, 
  sum((value -> 'units')::float) as units 
from 
  scaleway_billing_consumption 
group by 
  project_id;
```
