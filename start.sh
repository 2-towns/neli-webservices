#!/bin/bash

seeds --config ./docker.ini
neli-webservices --config ./docker.ini
cp crontab /etc/cron.d/hello-cron
cp cleaning.sh /cleaning.sh
chmod a+x /cleaning.sh
touch /var/log/cron.log
cron && tail -f /var/log/cron.log