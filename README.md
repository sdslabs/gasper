# SDS

> SDS Deployment Services

## Setup

- Clone the repository in your `$GOPATH/src/github.com/sdslabs`, hence the final path of package should be - `$GOPATH/src/github.com/sdslabs/SDS`

  ```shell
  # Create directory if not available
  $ mkdir -p $GOPATH/src/github.com/sdslabs
  # Clone the repo there
  $ git clone https://github.com/sdslabs/SDS  $GOPATH/src/github.com/sdslabs/SDS
  ```

- Make sure you have `dep` installed, if not, you can find it here: [github.com/golang/dep](https://github.com/golang/dep)

  ```shell
  $ cd $GOPATH/src/github.com/sdslabs/SDS
  $ dep ensure
  ```

### Development

- While developing we recommend can use a tool like [Fresh](https://github.com/pilu/fresh)

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
