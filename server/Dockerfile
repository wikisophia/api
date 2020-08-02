FROM golang:1.12 AS build
WORKDIR /go/src/github.com/wikisophia/api/server
ENV GOPATH /go
ENV CGO_ENABLED 0
ENV GOOS=linux
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY ./ ./
RUN dep ensure
RUN go build -o api-arguments .


FROM ubuntu:18.04 AS release
LABEL maintainer="admin@wikisophia.net"
WORKDIR /usr/local/bin/
COPY --from=build /go/src/github.com/wikisophia/api/server/api-arguments .
EXPOSE 8001
ENTRYPOINT ["/usr/local/bin/api-arguments"]
