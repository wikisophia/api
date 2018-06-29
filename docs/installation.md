# Installation

To run locally, you'll first need to install:

- [Golang](https://golang.org/doc/install)
- [Dep](https://github.com/golang/dep#installation)

Clone the project with:

```sh
export GOPATH=$(go env GOPATH)
mkdir -p $GOPATH
git clone https://github.com/wikisophia/api-arguments.git $GOPATH/src/github.com/wikisophia/api-arguments
cd $GOPATH/src/github.com/wikisophia/api-arguments
```

Install the dependencies, or update them after a `git pull`:

```sh
$GOPATH/bin/dep ensure
```

Then build and run the app:

```sh
go build .
./api-arguments
```

You can access the API through `http://localhost:8081`.
```

## Using Postgres

By default, the app will store all the state changes from the website in memory.
If you have [Postgres](https://www.postgresql.org/) set up, it can also be configured
to use that.

```sh
WKSPH_STORAGE_TYPE=postgres .server
```

Just make sure you've created the user & database, and run
[the init script](../postgres/scripts/clear.sql) on it
to create all the tables.
