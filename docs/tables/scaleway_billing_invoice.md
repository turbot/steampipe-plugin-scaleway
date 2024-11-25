---
title: "Steampipe Table: scaleway_billing_invoice - Query Scaleway Invoices using SQL"
description: "Enables users to query Scaleway invoices, offering comprehensive billing and usage details for Scaleway cloud services."
---

# Table: scaleway_billing_invoice - Query Scaleway invoices using SQL

Scaleway invoices provide detailed records of charges for using Scaleway's cloud services. These invoices include a complete breakdown of costs for various resources and services within a Scaleway account.

## Table Usage Guide

The `scaleway_billing_invoice` table offers insights into billing information in Scaleway. It allows finance managers or cloud administrators to query invoice-specific details such as total amounts, billing periods, and associated organizations. Use this table to track expenses, verify charges, and manage cloud spending across different projects and timeframes.

## Examples

### Explore basic details of Scaleway invoices  
Retrieve invoice identifiers, associated organizations, and billing periods to track cloud expenses effectively.

```sql+postgres
select
  id,
  organization_id,
  billing_period,
  total_taxed_amount,
  state,
  currency
from
  scaleway_billing_invoice;
```

```sql+sqlite
select
  id,
  organization_id,
  billing_period,
  total_taxed_amount,
  state,
  currency
from
  scaleway_billing_invoice;
```

### Get total billed amount for each organization  
Calculate the total amount billed for each organization to analyze spending across different entities.

```sql+postgres
select
  organization_id,
  sum(total_taxed_amount) as total_billed,
  currency
from
  scaleway_billing_invoice
group by
  organization_id,
  currency;
```

```sql+sqlite
select
  organization_id,
  sum(total_taxed_amount) as total_billed,
  currency
from
  scaleway_billing_invoice
group by
  organization_id,
  currency;
```

### Find invoices with high discount amounts  
Identify invoices with substantial discounts to understand cost-saving opportunities.

```sql+postgres
select
  id,
  billing_period,
  total_discount_amount,
  total_taxed_amount,
  currency
from
  scaleway_billing_invoice
where
  total_discount_amount > 1000
order by
  total_discount_amount desc;
```

```sql+sqlite
select
  id,
  billing_period,
  total_discount_amount,
  total_taxed_amount,
  currency
from
  scaleway_billing_invoice
where
  total_discount_amount > 1000
order by
  total_discount_amount desc;
```

### List invoices within a specific date range  
Retrieve invoices for a defined period to assist with financial reviews or audits.

```sql+postgres
select
  id,
  billing_period,
  total_taxed_amount,
  issued_date,
  currency
from
  scaleway_billing_invoice
where
  issued_date between '2023-01-01' and '2023-12-31'
order by
  issued_date;
```

```sql+sqlite
select
  id,
  billing_period,
  total_taxed_amount,
  issued_date,
  currency
from
  scaleway_billing_invoice
where
  issued_date between '2023-01-01' and '2023-12-31'
order by
  issued_date;
```

### Get the average invoice amount by month  
Analyze monthly spending patterns by calculating the average invoice amount.

```sql+postgres
select
  date_trunc('month', issued_date) as month,
  avg(total_taxed_amount) as average_invoice_amount,
  currency
from
  scaleway_billing_invoice
group by
  date_trunc('month', issued_date),
  currency
order by
  month;
```

```sql+sqlite
select
  strftime('%Y-%m', issued_date) as month,
  avg(total_taxed_amount) as average_invoice_amount,
  currency
from
  scaleway_billing_invoice
group by
  strftime('%Y-%m', issued_date),
  currency
order by
  month;
```

### Compare total taxed and untaxed amounts  
Analyze the tax impact by comparing taxed and untaxed amounts.

```sql+postgres
select
  id,
  total_untaxed_amount,
  total_taxed_amount,
  total_taxed_amount - total_untaxed_amount as tax_amount,
  currency
from
  scaleway_billing_invoice
order by
  tax_amount desc;
```

```sql+sqlite
select
  id,
  total_untaxed_amount,
  total_taxed_amount,
  total_taxed_amount - total_untaxed_amount as tax_amount,
  currency
from
  scaleway_billing_invoice
order by
  tax_amount desc;
```

### Examine discounts and their impact  
Evaluate the effect of discounts on invoices by comparing undiscounted and final taxed amounts.

```sql+postgres
select
  id,
  total_undiscount_amount,
  total_discount_amount,
  total_taxed_amount,
  total_discount_amount / total_undiscount_amount * 100 as discount_percentage,
  currency
from
  scaleway_billing_invoice
where
  total_discount_amount > 0
order by
  discount_percentage desc;
```

```sql+sqlite
select
  id,
  total_undiscount_amount,
  total_discount_amount,
  total_taxed_amount,
  total_discount_amount / total_undiscount_amount * 100 as discount_percentage,
  currency
from
  scaleway_billing_invoice
where
  total_discount_amount > 0
order by
  discount_percentage desc;
```
