# Configuration

Configuration is done through environment variables.

If no environment variables are set, some defaults will be used that work well for local development.
One caveat is that all the data will be stored in memory, so it will be lost when the process shuts down.

For a list of environment variables and docs on what they all do, see
the [TestEnvironmentOverrides](../config/parsing_test.go) unit test.

All the config variables (excluding passwords) will be printed to the terminal on app startup.

For convenience, you may want to define your custom config variables in a `config.env` file in [the server directory](..):

```
export WKSPH_ARGS_SERVER_ADDR=127.0.0.1:8002
export WKSPH_ARGS_STORAGE_TYPE=postgres
export WKSPH_ARGS_SERVER_USE_SSL=true
```

You can then source this file to set your environment:

```bash
source ./config.env
```

Or, to avoid polluting your terminal's environment, run it in a bash subshell:

```bash
(source ./config.env && ./server)
```
