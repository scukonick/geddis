FROM ubuntu

WORKDIR /app
VOLUME ["/data"]

# adding config file
ADD config.example.toml /app/config.toml

# adding binary
ADD geddis /app/

ENTRYPOINT ["./geddis"]
