# Stage 1: Build
FROM golang:1.22 AS builder

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go build -o steam-exporter

FROM alpine:3.21.0

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the built binary and config file
COPY --from=builder /app/steam-exporter .

# Expose the metrics port
EXPOSE 8080

# Start the application
ENTRYPOINT ["./steam-exporter"]
