![image](https://hub.steampipe.io/images/plugins/turbot/scaleway-social-graphic.png)

# Scaleway Plugin for Steampipe

Use SQL to query infrastructure servers, networks, databases and more from your Scaleway project.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/scaleway)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/scaleway/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-scaleway/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install scaleway
```

Configure your [credentials](https://hub.steampipe.io/plugins/turbot/scaleway#credentials) and [config file](https://hub.steampipe.io/plugins/turbot/scaleway#configuration).

Run a query:

```sql
select
  name,
  id,
  created_at
from
  scaleway_account_ssh_key;
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-scaleway.git
cd steampipe-plugin-scaleway
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```sh
make
```

Configure the plugin:

```sh
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/scaleway.spc
```

Try it!

```shell
steampipe query
> .inspect scaleway
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-scaleway/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Scaleway Plugin](https://github.com/turbot/steampipe-plugin-scaleway/labels/help%20wanted)
