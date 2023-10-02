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
