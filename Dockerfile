FROM golang:alpine AS FinalBuilder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install basic packages
RUN apk add \
    gcc \
    g++

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go mod download

# Build image
RUN go build .

FROM alpine:latest AS RUNNER

WORKDIR /app

COPY --from=FinalBuilder /app/IPFS-Upload-Relay /app/app

# This container exposes port 8080 to the outside world
EXPOSE 8080/tcp

ENV MODE=prod

# Run the executable
CMD ["/app/app"]
