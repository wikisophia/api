# Configuration

Configuration uses [Viper](https://github.com/spf13/viper).

The service should work out of the box, storing arguments in memory.
Beware: using this default config, all data will be lost when the server shuts down.

To connect to an actual postgres databse, or change other config options, you
can either set environment variables or define a `config.yaml` file at the project root.

## Environment variables

|              Environment Variable            |   Description |
|----------------------------------------------|---------------|
| WKSPH_ARGS_SERVER_ADDR                       | TCP address to listen on. |
| WKSPH_ARGS_SERVER_READ_HEADER_TIMEOUT_MILLIS | The number of milliseconds that the server is willing to wait for a client to send the request headers |
| WKSPH_ARGS_SERVER_CORS_ALLOWED_ORIGINS       | Configures [CORS]() support so that the specified domains can use the API. If there are more than one, they should be space separated. |
| WKSPH_ARGS_STORAGE_TYPE                      | Can be "memory" or "postgres", depending on how the data should be stored. |

If using Postgres, the following params are allowed.

|          Environment Variable        |                      Description                       |
| ------------------------------------ | ------------------------------------------------------ |
| WKSPH_ARGS_STORAGE_POSTGRES_DBNAME   | The name of the postgres database to connect to.       |
| WKSPH_ARGS_STORAGE_POSTGRES_HOST     | The hostname of the machine where postgres lives.      |
| WKSPH_ARGS_STORAGE_POSTGRES_PORT     | The port of the machine where postgres listens.        |
| WKSPH_ARGS_STORAGE_POSTGRES_USER     | The name of the postgres user to connect as.           |
| WKSPH_ARGS_STORAGE_POSTGRES_PASSWORD | The password for the WKSPH_ARGS_STORAGE_POSTGRES_USER. |

## Using config.yaml

```yaml
server:
  addr: localhost:8001
  read_header_timeout_millis: 100
  cors_allowed_origins:
    - "localhost"
    - "some-domain.com"
storage:
  type: "postgres"
  postgres:
    dbname: "wikisophia"
    host: "localhost"
    port: 5432
    user: "whatever"
    password: "something-secret"
```
