# SDS

> SDS Deployment Services

## Prerequisites

1. [GoLang](https://golang.org/)
   [Download](https://golang.org/dl/)
2. [Docker](https://www.docker.com/)
   [Download / Get Started](https://www.docker.com/get-started)
3. [Nginx Docker Image](https://hub.docker.com/_/nginx/)

   ```shell
   # Install using
   $ docker pull nginx

   # For installation on MAC, the latest image might not be pulled successfully
   # Use a specific version tag instead, like so
   $ docker pull nginx:1.15.3
   ```

4. [Nginx](https://nginx.org/en/download.html) or [Apache](https://httpd.apache.org/download.cgi) on your machine

## Setup

- Clone the repository in your `$GOPATH/src/github.com/sdslabs`, hence the final path of package should be - `$GOPATH/src/github.com/sdslabs/SDS`

  ```shell
  # Create directory if not available
  $ mkdir -p $GOPATH/src/github.com/sdslabs

  # Clone the repo there
  $ git clone https://github.com/sdslabs/SDS $GOPATH/src/github.com/sdslabs/SDS

  # cd into the directory
  $ cd $GOPATH/src/github.com/sdslabs/SDS

  # Install all dependencies using dep
  $ dep ensure
  ```

### Development

- For development purposes we recommend using [Fresh](https://github.com/pilu/fresh)

  ```shell
  $ go get github.com/pilu/fresh
  ```

- Run the following command to start the server
  ```shell
  $ fresh
  ```

### Production

- Start server using
  ```shell
  $ go run server.go
  ```
