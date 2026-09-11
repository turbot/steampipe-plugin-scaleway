## v2.0.1 [2026-09-11]

_Bug fixes_

- Fixed the plugin module path to include the `/v2` suffix required for v2 releases, which caused the v2.0.0 Steampipe Anywhere (Postgres FDW, SQLite, and export) builds to fail.

## v2.0.0 [2026-09-11]

Thanks to [@pdecat](https://github.com/pdecat) for contributing all of the table, API, and SDK changes in this release!

_Breaking changes_

- Updated the `scaleway_vpc_private_network` table to use the Scaleway VPC v2 API, in which private networks are regional resources. The `zone` column has been replaced by `region`, and the get key columns are now `id` and `region` instead of `id` and `zone`. The table now iterates over the `regions` connection config argument instead of `zones` and connections that set `zones` without `regions` will only return private networks from the default region until `regions` is added. Queries that filter on or select `zone` must be updated. ([#171](https://github.com/turbot/steampipe-plugin-scaleway/pull/171))
- Removed the `creation_info` column from the `scaleway_account_ssh_key` table, as it is not returned by the IAM API. ([#170](https://github.com/turbot/steampipe-plugin-scaleway/pull/170))

_What's new?_

- New tables added
  - [scaleway_vpc](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_vpc) ([#171](https://github.com/turbot/steampipe-plugin-scaleway/pull/171))

_Enhancements_

- Added `vpc_id`, `dhcp_enabled`, and `subnets` columns to the `scaleway_vpc_private_network` table, with `vpc_id` available as an optional filter. ([#171](https://github.com/turbot/steampipe-plugin-scaleway/pull/171))
- Updated the `scaleway_account_ssh_key` table to use the Scaleway IAM API, as the Account v2alpha1 SSH key API has been removed and the table no longer returned data. Added the `disabled` column. ([#170](https://github.com/turbot/steampipe-plugin-scaleway/pull/170))

_Bug fixes_

- Fixed the `scaleway_account_ssh_key`, `scaleway_iam_api_key`, and `scaleway_iam_user` tables to return an error instead of causing a plugin panic when the API call fails. ([#169](https://github.com/turbot/steampipe-plugin-scaleway/pull/169))
- Fixed the plugin to correctly use credentials from the `SCW_ACCESS_KEY` and `SCW_SECRET_KEY` environment variables when `access_key` and `secret_key` are not set in the connection config. ([#169](https://github.com/turbot/steampipe-plugin-scaleway/pull/169))

_Dependencies_

- Recompiled plugin with Go version `1.26`. ([#159](https://github.com/turbot/steampipe-plugin-scaleway/pull/159))
- Recompiled plugin with [steampipe-plugin-sdk v6.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v610-2026-08-18) that fixes data races in the SDK's in-process execute path.

## v1.3.0 [2025-12-02]

_Enhancements_

- Added `billing_period` as a column in the `scaleway_billing_consumption` table and enabled it as an optional filter (qual) for queries. ([#154](https://github.com/turbot/steampipe-plugin-scaleway/pull/154)) (Thanks [@tdannenmuller](https://github.com/tdannenmuller) for the contribution!)

## v1.2.0 [2025-10-13]

_Dependencies_

- Recompiled plugin with Go version `1.24`. ([#147](https://github.com/turbot/steampipe-plugin-scaleway/pull/147))
- Recompiled plugin with [steampipe-plugin-sdk v5.13.1](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5131-2025-09-25) that addresses critical and high vulnerabilities in dependent packages. ([#152](https://github.com/turbot/steampipe-plugin-scaleway/pull/152))

## v1.1.0 [2025-04-17]

_What's new?_

- New tables added
  - [scaleway_billing_invoice](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_billing_invoice) ([#122](https://github.com/turbot/steampipe-plugin-scaleway/pull/122)) (Thanks [@tdannenmuller](https://github.com/tdannenmuller) for the contribution!)

_Dependencies_

- Recompiled plugin with Go version `1.23.1`. ([#137](https://github.com/turbot/steampipe-plugin-scaleway/pull/137))
- Recompiled plugin with [steampipe-plugin-sdk v5.11.5](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.11.5/CHANGELOG.md#v5115-2025-03-31) that addresses critical and high vulnerabilities in dependent packages. ([#137](https://github.com/turbot/steampipe-plugin-scaleway/pull/137))

## v1.0.0 [2024-10-22]

There are no significant changes in this plugin version; it has been released to align with [Steampipe's v1.0.0](https://steampipe.io/changelog/steampipe-cli-v1-0-0) release. This plugin adheres to [semantic versioning](https://semver.org/#semantic-versioning-specification-semver), ensuring backward compatibility within each major version.

_Dependencies_

- Recompiled plugin with Go version `1.22`. ([#118](https://github.com/turbot/steampipe-plugin-scaleway/pull/118))
- Recompiled plugin with [steampipe-plugin-sdk v5.10.4](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5104-2024-08-29) that fixes logging in the plugin export tool. ([#118](https://github.com/turbot/steampipe-plugin-scaleway/pull/118))

## v0.11.1 [2024-02-16]

_Bug fixes_

- Fixed the typo in the `scaleway_billing_consumption` table docs to use `consumption` instead of `consumtion`. ([#80](https://github.com/turbot/steampipe-plugin-scaleway/pull/80))

## v0.11.0 [2024-02-16]

_What's new?_

- New tables added
  - [scaleway_account_project](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_account_project) ([#53](https://github.com/turbot/steampipe-plugin-scaleway/pull/53)) (Thanks [@jplanckeel](https://github.com/jplanckeel) for the contribution!)
  - [scaleway_billing_consumption](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_billing_consumption) ([#70](https://github.com/turbot/steampipe-plugin-scaleway/pull/70)) (Thanks [@jplanckeel ](https://github.com/jplanckeel) for the contribution!)

## v0.10.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#67](https://github.com/turbot/steampipe-plugin-scaleway/pull/67))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#67](https://github.com/turbot/steampipe-plugin-scaleway/pull/67))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-scaleway/blob/main/docs/LICENSE). ([#67](https://github.com/turbot/steampipe-plugin-scaleway/pull/67))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#66](https://github.com/turbot/steampipe-plugin-scaleway/pull/66))

## v0.9.1 [2023-10-05]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#45](https://github.com/turbot/steampipe-plugin-scaleway/pull/45))

## v0.9.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#41](https://github.com/turbot/steampipe-plugin-scaleway/pull/41))
- Recompiled plugin with Go version `1.21`. ([#41](https://github.com/turbot/steampipe-plugin-scaleway/pull/41))

## v0.8.0 [2023-07-17]

_Enhancements_

- Updated the `docs/index.md` file to include multi-project configuration examples. ([#28](https://github.com/turbot/steampipe-plugin-scaleway/pull/28))

## v0.7.0 [2023-03-31]

_What's new?_

- New tables added
  - [scaleway_baremetal_server](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_baremetal_server) ([#17](https://github.com/turbot/steampipe-plugin-scaleway/pull/17)) (Thanks [@jplanckeel](https://github.com/jplanckeel) for the contribution!)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#26](https://github.com/turbot/steampipe-plugin-scaleway/pull/26))

## v0.6.0 [2023-02-16]

_What's new?_

- New tables added
  - [scaleway_kubernetes_cluster](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_kubernetes_cluster) ([#23](https://github.com/turbot/steampipe-plugin-scaleway/pull/23)) (Thanks to [@jplanckeel](https://github.com/jplanckeel) for the contribution!)
  - [scaleway_kubernetes_node](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_kubernetes_node) ([#23](https://github.com/turbot/steampipe-plugin-scaleway/pull/23)) (Thanks to [@jplanckeel](https://github.com/jplanckeel) for the contribution!)
  - [scaleway_kubernetes_pool](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_kubernetes_pool) ([#23](https://github.com/turbot/steampipe-plugin-scaleway/pull/23)) (Thanks to [@jplanckeel](https://github.com/jplanckeel) for the contribution!)

## v0.5.0 [2023-01-25]

_What's new?_

- New tables added
  - [scaleway_registry_image](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_registry_image) ([#20](https://github.com/turbot/steampipe-plugin-scaleway/pull/20)) (Thanks to [@jplanckeel](https://github.com/jplanckeel) for the contribution!)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.11](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v4111-2023-01-24) which fixes the issue of non-caching of all the columns of the queried table. ([#21](https://github.com/turbot/steampipe-plugin-scaleway/pull/21))

## v0.4.0 [2022-12-08]

_What's new?_

- New tables added
  - [scaleway_registry_namespace](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_registry_namespace) ([#19](https://github.com/turbot/steampipe-plugin-scaleway/pull/19)) (Thanks to [@jplanckeel](https://github.com/jplanckeel) for the contribution!)

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.8](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v418-2022-09-08) which increases the default open file limit. ([#18](https://github.com/turbot/steampipe-plugin-scaleway/pull/18))

## v0.3.0 [2022-09-27]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#14](https://github.com/turbot/steampipe-plugin-scaleway/pull/14))
- Recompiled plugin with Go version `1.19`. ([#14](https://github.com/turbot/steampipe-plugin-scaleway/pull/14))

## v0.2.1 [2022-05-23]

_Bug fixes_

- Fixed the Slack community links in README and docs/index.md files. ([#10](https://github.com/turbot/steampipe-plugin-scaleway/pull/10))

## v0.2.0 [2022-04-28]

_Enhancements_

- Added support for native Linux ARM and Mac M1 builds. ([#8](https://github.com/turbot/steampipe-plugin-scaleway/pull/8))
- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#7](https://github.com/turbot/steampipe-plugin-scaleway/pull/7))

## v0.1.0 [2021-12-16]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) ([#4](https://github.com/turbot/steampipe-plugin-scaleway/pull/4))
- Recompiled plugin with Go version 1.17 ([#4](https://github.com/turbot/steampipe-plugin-scaleway/pull/4))

## v0.0.2 [2021-10-13]

_Bug fixes_

- Fixed: Typo in `scaleway_instance_snapshot` table document filename

## v0.0.1 [2021-10-12]

_What's new?_

- New tables added

  - [scaleway_account_ssh_key](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_account_ssh_key)
  - [scaleway_instance_image](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_instance_image)
  - [scaleway_instance_ip](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_instance_ip)
  - [scaleway_instance_security_group](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_instance_security_group)
  - [scaleway_instance_server](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_instance_server)
  - [scaleway_instance_snapshot](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_instance_snapshot)
  - [scaleway_instance_volume](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_instance_volume)
  - [scaleway_object_bucket](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_object_bucket)
  - [scaleway_rdb_database](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_rdb_database)
  - [scaleway_rdb_instance](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_rdb_instance)
  - [scaleway_vpc_private_network](https://hub.steampipe.io/plugins/turbot/scaleway/tables/scaleway_vpc_private_network)
