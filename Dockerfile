FROM golang:alpine AS Builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install basic packages
RUN apk add \
    gcc \
    g++ \
    curl

# Download ytdl
RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /app/youtube-dl

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go mod download

# Build image
RUN go build .

FROM alpine:latest AS Runner

WORKDIR /app

COPY --from=Builder /app/IPFS-Upload-Relay /app/app
COPY --from=Builder /app/youtube-dl /app/youtube-dl

# Set correct attribte
RUN chmod a+rx /app/youtube-dl

# This container exposes port 8080 to the outside world
EXPOSE 8080/tcp

ENV MODE=prod

# Run the executable
CMD ["/app/app"]
