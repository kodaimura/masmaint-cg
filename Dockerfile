FROM golang:1.22

WORKDIR /usr/src/app

RUN apt-get update \
    && apt-get install -y zip unzip \
    && apt-get install -y cron logrotate \
    && rm -rf /var/lib/apt/lists/*

COPY ./linux/cron.d/app /etc/cron.d/app
RUN chmod 0644 /etc/cron.d/app
RUN crontab /etc/cron.d/app