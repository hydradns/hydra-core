# Lightweight base image
FROM alpine:3.23

# Install sqlite and useful debugging tools
RUN apk add --no-cache \
    sqlite \
    bash \
    coreutils \
    busybox-extras

# Create default sqlite config for root
RUN printf ".headers on\n.mode column\n.nullvalue NULL\n" > /root/.sqliterc

# Set working directory to where the DB will be mounted
WORKDIR /app/data

