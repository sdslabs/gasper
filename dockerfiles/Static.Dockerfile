FROM nginx:1.17.4-alpine

RUN rm /etc/nginx/conf.d/default.conf \
	&& apk --no-cache add curl
