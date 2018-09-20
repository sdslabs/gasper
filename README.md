# SDS

> SDS Deployment Services

## Setup

- Clone the repository in your `$GOPATH/src/github.com/sdslabs`, hence the final path of package should be - `$GOPATH/src/github.com/sdslabs/SDS`

- Make sure you have `dep` installed, if not, you can find it here: [github.com/golang/dep](https://github.com/golang/dep)

- `cd` into the `SDS` directory and run the following command:
  ```shell
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
