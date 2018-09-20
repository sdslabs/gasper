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

- While developing we recommend can use a tool like [Fresh](https://github.com/pilu/fresh) or [Gomon](https://github.com/c9s/gomon)

  ```shell
  # Fresh
  $ go get github.com/pilu/fresh

  # Gomon
  $ go get -u github.com/c9s/gomon
  ```

- Run the following command to start the server

  ```shell
  # Using fresh:
  $ fresh

  # Using gomon:
  $ gomon src -- go run server.go
  ```

### Production

- Start server using
  ```shell
  $ go run server.go
  ```
