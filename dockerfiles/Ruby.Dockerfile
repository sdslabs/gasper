FROM ruby:2.6.5-alpine3.10

RUN apk --no-cache add curl build-base sqlite-dev mariadb-dev postgresql-dev nodejs tzdata

CMD tail -f /dev/null
