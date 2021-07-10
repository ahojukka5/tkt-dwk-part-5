FROM ubuntu:20.04
RUN apt-get update && apt-get install -y curl wget jq bash \
    && rm -rf /var/lib/apt/lists/*
COPY entrypoint.sh ./entrypoint.sh
CMD ./entrypoint.sh
