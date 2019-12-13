FROM node:8.16.1-alpine

RUN apk --no-cache add curl

CMD tail -f /dev/null
