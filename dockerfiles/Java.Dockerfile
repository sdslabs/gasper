FROM openjdk:14-jdk-alpine3.10

RUN apk --no-cache add maven curl

CMD tail -f /dev/null
