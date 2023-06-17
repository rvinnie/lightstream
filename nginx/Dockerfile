FROM alpine:3.11.5

ENV NGINX_VERSION=1.16.1-r8

EXPOSE 80 443

RUN apk add --no-cache nginx=${NGINX_VERSION} \
    && apk add openssl \
    && mkdir -p /run/nginx \
    && mkdir -p /var/www/lightstream/website

COPY ./start.sh /var/www

RUN sh /var/www/start.sh

CMD nginx -g 'daemon off;';
