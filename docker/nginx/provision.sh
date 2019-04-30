#!/bin/sh

echo "Starting container provisioning"

# crontab setup
if [[ -f /provision/cron/crontab ]]; then
    cp /provision/cron/crontab /var/spool/cron/crontabs/root

    chown root:root /var/spool/cron/crontabs/root
    chmod 600 /var/spool/cron/crontabs/root

    echo "crontab setup successfully"
fi

# logs rotation setup
cp /provision/logrotate_nginx.conf /etc/logrotate.d
chown -R root:root /etc/logrotate.d
chmod -R 644 /etc/logrotate.d

echo "Container provisioning finished"

exec "$@"