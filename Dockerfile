# Use the official Golang image as the build stage
FROM golang:latest AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . /app

# Build the Go app
RUN go build -o /app/server ./cmd/server/main.go

# Start a new stage from a Debian image
FROM debian:latest

# Install ca-certificates to avoid issues with secure connections
RUN apt-get update && apt-get install -y ca-certificates
RUN apt-get update && apt-get install -y curl postgresql-client
WORKDIR /root/

# Set environment variables for the database connection
ENV POSTGRES_DB=WB_developer
ENV POSTGRES_USER=mbatimel
ENV POSTGRES_PASSWORD=wb_il
ENV POSTGRES_HOST=postgres
ENV POSTGRES_PORT=5432


# Copy the pre-built binary files from the previous stage
COPY --from=builder /app/server .


# Copy the config file
COPY config/config.yaml config/config.yaml
COPY migrations/migrate.sql migrations/migrate.sql

# Expose the port the app runs on
EXPOSE 8080

# # Run the server by default
ENTRYPOINT ["./server"]

