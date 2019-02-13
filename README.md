# SWS

> SDS Web Services

## Prerequisites

1. [GoLang](https://golang.org/)

   [Download](https://golang.org/dl/)

2. [Docker](https://www.docker.com/)

   [Download / Get Started](https://www.docker.com/get-started)

3. [MongoDB](https://www.mongodb.com/)

   [Download](https://www.mongodb.com/download-center/community)

4. [Nginx](https://nginx.org/en/download.html) or [Apache](https://httpd.apache.org/download.cgi) on your machine

## Setup

- Clone the repository
- `cp config.sample.json config.json`
- Start hacking

**Note:** The vendor is committed, to add another package as dependency, `go get ...` the package in your gopath and then run the command `go mod vendor` to add the dependency in the SWS package.

*To use go-modules you must have Golang version 1.11 or later. Also remember to set the environment variable `GO111MODULE=on`. For reference see - [https://github.com/golang/go/wiki/Modules](https://github.com/golang/go/wiki/Modules).*

### Development

- Set `debug` to `true` in `config.json`

- For development purposes we recommend using [Fresh](https://github.com/pilu/fresh)

  ```shell
  $ go get github.com/pilu/fresh
  ```

- Run the following command to start the server
  ```shell
  $ fresh
  ```

### Production

- Set `debug` to `false` in `config.json`

- Start server using
  ```shell
  $ go run server.go
  ```
