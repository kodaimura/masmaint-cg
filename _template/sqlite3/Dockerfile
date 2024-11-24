FROM golang:1.22

WORKDIR /usr/src/app

RUN apt-get update \
    && apt-get install -y sqlite3 \
    && apt-get install -y logrotate \
    && rm -rf /var/lib/apt/lists/*