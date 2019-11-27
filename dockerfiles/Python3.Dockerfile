FROM python:3.7.4-alpine3.10

RUN apk --no-cache add curl

CMD tail -f /dev/null
