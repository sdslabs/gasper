# Gasper

> Your Cloud in a Binary

<img align="right" width="300px" src="./docs/assets/logo-11.svg">

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/sdslabs/gasper/blob/develop/LICENSE.md)

Gasper is an intelligent Platform as a Service (PaaS) used for deploying and managing 
applications and databases in any cloud topology.

## Overview


## Development

The following softwares are required for running Gasper:-

* [Golang 1.12.x](https://golang.org/dl/)
* [Docker](https://www.docker.com/)
    * [For Linux](https://runnable.com/docker/install-docker-on-linux)
    * [For MacOS](https://docs.docker.com/docker-for-mac/install/)
    * [For Windows](https://docs.docker.com/docker-for-windows/install/)
* [MongoDB](https://www.mongodb.com/download-center/community)
* [Redis](https://redis.io/download)

Once you are done installing the pre-requisites, open a terminal and start typing the following
```bash
$ go version
go version go1.12.7 darwin/amd64

$ git clone https://github.com/sdslabs/gasper

$ cd gasper && make help

 Gasper: Your cloud in a binary

  install   Install missing dependencies
  build     Build the project binary
  tools     Install development tools
  start     Start in development mode with hot-reload enabled
  clean     Clean build files
  fmt       Format entire codebase
  vet       Vet entire codebase
  help      Display this help
```


## Setup

- Clone the repository
- `cp config.sample.json config.json`
- Start hacking

**Note:** The vendor is committed, to add another package as dependency, `go get ...` the package in your gopath and then run the command `go mod vendor` to add the dependency in the gasper package.


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
