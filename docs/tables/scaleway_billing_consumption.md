---
title: "Steampipe Table: scaleway_billing_consumtion - Query Scaleway Billing comsumption using SQL"
description: "Allows users to query Scaleway Billing Consumption, specifically providing information about the billing"
---

# Table: scaleway_billing_consumtion - Query Scaleway Billing comsumption using SQL

A Scaleway Billing Consumption track the costs of various scaleway products such as object_storage, storage, serverless, ect.

This table requires an Organization ID to be configured in the scaleway.spc file.

## Table Usage Guide

The `scaleway_billing_consumtion` table offers insights into the Billing Consumption within Scaleway.

## Examples

### Basic info
Explore the billing for all scaleway service.

```sql+postgres
select
  category,
  operation_path,
  project_id,
  description,
  value
from
  scaleway_billing_consumption;
```

```sql+sqlite
select
  category,
  operation_path,
  project_id,
  description,
  value
from
  scaleway_billing_consumption;
```

### List consumption by category
Explore costs by product category of scaleway

```sql+postgres
select 
  category, 
  sum((value -> 'units')::float) AS units 
from 
  scaleway_billing_consumption 
group by 
  category
```

```sql+sqlite
select 
  category, 
  sum((value -> 'units')::float) AS units 
from 
  scaleway_billing_consumption 
group by 
  category
```

### List consumption by project_id
Explore costs by product category of scaleway

```sql+postgres
select 
  project_id, 
  sum((value -> 'units')::float) AS units 
from 
  scaleway_billing_consumption 
group by 
  project_id
```

```sql+sqlite
select 
  project_id, 
  sum((value -> 'units')::float) AS units 
from 
  scaleway_billing_consumption 
group by 
  project_id
```
