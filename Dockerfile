# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21 AS build-stage

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o rms-paystack

# Deploy the application binary into a lean image
FROM ubuntu:22.04 AS build-release-stage

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /server

# Copy the binary from the build stage
COPY --from=build-stage /app/rms-paystack /server/

# Create a non-root user and switch to it
RUN useradd -m nonroot
USER nonroot

EXPOSE 9001

ENV DOCKERIZE_VERSION v0.7.0

USER root
RUN apt-get update \
    && apt-get install -y wget \
    && wget -O - https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz | tar xzf - -C /usr/local/bin \
    && apt-get remove -y wget \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/*

USER nonroot

CMD ["dockerize", "-wait", "tcp://postgres:5432", "/server/rms-paystack"]