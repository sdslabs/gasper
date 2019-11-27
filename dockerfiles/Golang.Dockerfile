FROM golang:alpine3.10

RUN apk --no-cache add curl

CMD tail -f /dev/null
