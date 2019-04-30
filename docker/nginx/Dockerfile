FROM nginx:1.13-alpine

RUN apk update && apk upgrade && \
    \
    apk add bash supervisor logrotate && \
    \
    rm /etc/logrotate.d/*

COPY ./supervisor.conf /etc/supervisor.d/app.ini
COPY ./logrotate_nginx.conf /provision/logrotate_nginx.conf
COPY ./provision.sh /provision/provision.sh

ENTRYPOINT ["/provision/provision.sh"]
CMD ["/usr/bin/supervisord", "--nodaemon"]