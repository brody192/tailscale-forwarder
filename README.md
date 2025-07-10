## Tailscale Forwarder

Tailscale Forwarder is a TCP proxy that allows you to connect through a Tailscale machine to the configured target address and port pair.

This allows you to connect to Railway services that are not accessible from the internet, for example, locking down access to your database to only those who are on your Tailscale network.

This also solves for the issue that you can only run one Tailscale subnet router per Tailscale account.

## Usage

1. Generate a Tailscale [auth key](https://tailscale.com/kb/1085/auth-keys).

   Make sure `Reusable` is enabled.

2. Enable [MagicDNS](https://tailscale.com/kb/1081/magicdns) for your Tailscale account.

   This is required so that your computer can resolve the Tailscale Forwarder machine's short hostname to the correct IP address.   

3. Deploy the Tailscale Forwarder service into your pre-existing Railway project.

   Set the `TS_AUTHKEY` environment variable to the auth key you generated in step 1.

   Set your first connection mapping, example:

   `CONNECTION_MAPPING_01=5432:${{Postgres.RAILWAY_PRIVATE_DOMAIN}}:${{Postgres.PGPORT}}`

   The format is `<Source Port>:<Target Host>:<Target Port>`.

   Note: You can set multiple connection mappings by incrementing the `CONNECTION_MAPPING_` prefix.

4. Get the machine's hostname.

   You should see a new machine in the Tailscale [dashboard](https://login.tailscale.com/admin/machines) with the format `<Project Name>-<Environment Name>-<Service Name>`.
   
   Copy this hostname.

5. Use the machine's hostname to connect to the service.

   Example: `postgresql://postgres:<Postgres Password>@<Tailscale Forwarder Hostname>:<Source Port From Desired Connection Mapping>/railway`

   While that example is for a PostgreSQL connection string, you can use the same `<Tailscale Forwarder Hostname>:<Source Port From Desired Connection Mapping>` format to connect to any service that listens on a TCP port, that you have setup a connection mapping for.

## Configuration

| Environment Variable     | Required | Default Value                                                                       | Description                                |
| ------------------------ | :------: | ----------------------------------------------------------------------------------- | ------------------------------------------ |
| `TS_AUTHKEY`             | Yes      | -                                                                                   | Tailscale auth key.                        |
| `TS_HOSTNAME`            | Yes      | `${{RAILWAY_PROJECT_NAME}}-${{RAILWAY_ENVIRONMENT_NAME}}-${{RAILWAY_SERVICE_NAME}}` | Hostname to use for the Tailscale machine. |
| `CONNECTION_MAPPING_[n]` | Yes      | -                                                                                   | Connection mapping for a service.          |

## Examples

For all these examples, lets assume that the Tailscale Forwarder machine is named `my-project-production-tailscale-forwarder`.

#### Redis

Set the connection mapping:

```shell
CONNECTION_MAPPING_01=6379:${{Redis.RAILWAY_PRIVATE_DOMAIN}}:${{Redis.REDISPORT}}
```

If your Redis service is named anything other than `Redis`, you can change the namespace in the reference variable.

Connect to Redis with:

```shell
redis://default:<password>@my-project-production-tailscale-forwarder:6379
```

#### ClickHouse

Set the connection mapping:

```shell
CONNECTION_MAPPING_01=8123:${{ClickHouse.RAILWAY_PRIVATE_DOMAIN}}:${{ClickHouse.PORT}}
```

If your ClickHouse service is named anything other than `ClickHouse`, you can change the namespace in the reference variable.

Connect to ClickHouse with:

```shell
http://clickhouse:<password>@my-project-production-tailscale-forwarder:8123/railway
```

#### A Web Server

Set the connection mapping:

```shell
CONNECTION_MAPPING_01=80:${{Web Server.RAILWAY_PRIVATE_DOMAIN}}:${{Web Server.PORT}}
```

If your web server service is named anything other than `Web Server`, you can change the namespace in the reference variable.

You may also need to add a `PORT` environment variable to the service, if it is not already set.

Connect to the web server with:

```shell
http://my-project-production-tailscale-forwarder:80
```

#### Multiple Services

Set the connection mappings:

```shell
CONNECTION_MAPPING_01=5432:${{Postgres.RAILWAY_PRIVATE_DOMAIN}}:${{Postgres.PGPORT}}
CONNECTION_MAPPING_02=6379:${{Redis.RAILWAY_PRIVATE_DOMAIN}}:${{Redis.REDISPORT}}
CONNECTION_MAPPING_03=8123:${{ClickHouse.RAILWAY_PRIVATE_DOMAIN}}:${{ClickHouse.PORT}}
CONNECTION_MAPPING_04=80:${{Web Server.RAILWAY_PRIVATE_DOMAIN}}:${{Web Server.PORT}}
```

Then you can connect to the services by substituting in the `my-project-production-tailscale-forwarder` hostname with the set source port from the connection mapping.