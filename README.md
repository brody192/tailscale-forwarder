## Tailscale Forwarder

Tailscale Forwarder is a TCP proxy that allows you to connect through a Tailscale machine to the configured target address and port pair.

This allows you to connect to Railway services that are not accessible from the internet, for example, locking down access to your database to only those who are on your Tailscale network.

This also solves for the issue that you can only run one Tailscale subnet router per Tailscale account.

## Usage

1. Generate a Tailscale auth key.

   Make sure `Reusable` is enabled.

2. Deploy the Tailscale Forwarder service into your pre-existing Railway project.

   Set the `TS_AUTHKEY` environment variable to the auth key you generated in step 1.

   Set your first connection mapping, example:

   `CONNECTION_MAPPING_01=5432:${{Postgres.RAILWAY_PRIVATE_DOMAIN}}:${{Postgres.PGPORT}}`

   The format is `<Source Port>:<Target Host>:<Target Port>`.

   Note: You can set multiple connection mappings by incrementing the `CONNECTION_MAPPING_` prefix.

3. Get the machine's hostname.

   You should see a new machine in the dashboard with the format `<Project Name>-<Environment Name>-<Service Name>`, copy this hostname.

4. Use the machine's hostname in your connection string.

   Example: `postgresql://postgres:<Postgres Password>@<Tailscale Forwarder Hostname>:<Source Port From Desired Connection Mapping (5432)>/railway`

   While that example is for a PostgreSQL connection string, you can use the same `<Tailscale Forwarder Hostname>:<Source Port From Desired Connection Mapping (5432)>` format to connect to any service that listens on a TCP port, that you have setup a connection mapping for.

## Configuration

| Environment Variable     | Required | Default Value                                                                       | Description                                |
| ------------------------ | :------: | ----------------------------------------------------------------------------------- | ------------------------------------------ |
| `TS_AUTHKEY`             | Yes      | -                                                                                   | Tailscale auth key.                        |
| `TS_HOSTNAME`            | Yes      | `${{RAILWAY_PROJECT_NAME}}-${{RAILWAY_ENVIRONMENT_NAME}}-${{RAILWAY_SERVICE_NAME}}` | Hostname to use for the Tailscale machine. |
| `CONNECTION_MAPPING_[n]` | Yes      | -                                                                                   |                                            |
