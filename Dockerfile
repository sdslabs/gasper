FROM golang:1.16.0-alpine3.13 AS builder

WORKDIR /go/src/github.com/sdslabs/gasper

COPY . .

ARG vers

RUN apk update && \
  apk add make && \
  apk add bash

RUN make release VERSION=$vers

# Copy binary into actual image

FROM alpine:3.12.0

WORKDIR /go/bin

COPY --from=builder /go/src/github.com/sdslabs/gasper/releases/gasper_v1.0_linux_386/gasper .
COPY --from=builder /go/src/github.com/sdslabs/gasper/config.toml .
COPY --from=builder /go/src/github.com/sdslabs/gasper/public/frontend.go .

CMD [ "./gasper", "version" ]
