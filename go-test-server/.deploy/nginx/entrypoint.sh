#!/bin/sh

envsubst '$$APP_HOST$$APP_PORT' \
         < /etc/nginx/conf.d/default.conf.template \
         > /etc/nginx/conf.d/default.conf

nginx -g "daemon off;"
