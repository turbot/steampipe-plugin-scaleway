---
title: "Steampipe Table: scaleway_invoice - Query Scaleway Invoices using SQL"
description: "Allows users to query Scaleway Invoices, providing detailed information about billing and usage for Scaleway services."
---

# Table: scaleway_invoice - Query Scaleway Invoices using SQL

Scaleway Invoices are detailed records of charges for the use of Scaleway's cloud services. These invoices provide a comprehensive breakdown of costs associated with various resources and services used within a Scaleway account.

## Table Usage Guide

The `scaleway_invoice` table provides insights into billing information within Scaleway. As a finance manager or cloud administrator, explore invoice-specific details through this table, including total amounts, billing periods, and associated organizations. Utilize it to track expenses, verify charges, and manage cloud spending across different projects and timeframes.

## Examples

### Basic info
Explore the basic details of your Scaleway invoices, including their unique identifiers, associated organizations, and billing periods. This can help in tracking and managing your cloud expenses effectively.

```sql
SELECT
  id,
  organization_id,
  billing_period,
  total_taxed_amount,
  state,
  currency
FROM
  scaleway_invoices;
```

### Get total billed amount for each organization
Calculate the total amount billed to each organization. This provides an overview of cloud spending across different entities within your Scaleway account.

```sql
SELECT
  organization_id,
  SUM(total_taxed_amount) as total_billed,
  currency
FROM
  scaleway_invoices
GROUP BY
  organization_id,
  currency;
```

### Find invoices with high discount amounts
Identify invoices with significant discounts. This can help in understanding which billing periods or services are providing the most cost savings.

```sql
SELECT
  id,
  billing_period,
  total_discount_amount,
  total_taxed_amount,
  currency
FROM
  scaleway_invoices
WHERE
  total_discount_amount > 1000
ORDER BY
  total_discount_amount DESC;
```

### List invoices for a specific date range
Retrieve invoices within a specific time frame. This is useful for periodic financial reviews or audits.

```sql
SELECT
  id,
  billing_period,
  total_taxed_amount,
  issued_date,
  currency
FROM
  scaleway_invoices
WHERE
  issued_date BETWEEN '2023-01-01' AND '2023-12-31'
ORDER BY
  issued_date;
```

### Get the average invoice amount by month
Calculate the average invoice amount for each month. This helps in understanding monthly spending patterns and budgeting for cloud services.

```sql
SELECT
  DATE_TRUNC('month', issued_date) AS month,
  AVG(total_taxed_amount) AS average_invoice_amount,
  currency
FROM
  scaleway_invoices
GROUP BY
  DATE_TRUNC('month', issued_date),
  currency
ORDER BY
  month;
```

### Compare total taxed and untaxed amounts
Analyze the difference between taxed and untaxed amounts for each invoice to understand the tax impact on your cloud spending.

```sql
SELECT
  id,
  total_untaxed_amount,
  total_taxed_amount,
  total_taxed_amount - total_untaxed_amount AS tax_amount,
  currency
FROM
  scaleway_invoices
ORDER BY
  tax_amount DESC;
```

### Examine discounts and their impact
Investigate how discounts affect your invoices by comparing the undiscounted amount to the final taxed amount.

```sql
SELECT
  id,
  total_undiscount_amount,
  total_discount_amount,
  total_taxed_amount,
  total_discount_amount / total_undiscount_amount * 100 AS discount_percentage,
  currency
FROM
  scaleway_invoices
WHERE
  total_discount_amount > 0
ORDER BY
  discount_percentage DESC;
```
