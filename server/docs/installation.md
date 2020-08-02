# Installation

Make sure to install [Golang](https://golang.org/doc/install).

To clone the project:

```sh
git clone https://github.com/wikisophia/api.git
```

To build the project:

```sh
cd api/server
go build .
```

To run the server:

```sh
./api
```

By default, the API will listen on `http://localhost:8081` and store all state in memory only.

To use Postgres or set other config options, see [the config docs](./configuration.md).
