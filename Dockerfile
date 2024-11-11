FROM golang:1.22

WORKDIR /usr/src/app

RUN apt-get update \
    && apt-get install -y zip unzip \
    && rm -rf /var/lib/apt/lists/*