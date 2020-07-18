# Installation

Make sure to install [Golang](https://golang.org/doc/install).

Clone the project with:

```sh
export GOPATH=$(go env GOPATH)
mkdir -p $GOPATH
git clone https://github.com/wikisophia/api-arguments.git $GOPATH/src/github.com/wikisophia/api-arguments
cd $GOPATH/src/github.com/wikisophia/api-arguments/server
```

Then build and run the app:

```sh
go build .
./api-arguments
```

You can hit the API through `http://localhost:8081`.

Beware that the default app config will store everything in memory.
To use Postgres, set [the config options](./configuration.md).
