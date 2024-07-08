# Use the official Golang image as the build stage
FROM golang:latest AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o ./migration ./cmd/migration/main.go
RUN go build -o ./server ./cmd/server/main.go

# Start a new stage from scratch
FROM alpine:latest

# Set environment variables for the database connection
ENV POSTGRES_DB=WB_developer
ENV POSTGRES_USER=mbatimel
ENV POSTGRES_PASSWORD=wb_il
ENV POSTGRES_HOST=postgres
ENV POSTGRES_PORT=5432

# Create a directory for the app
WORKDIR /app

# Copy the pre-built binary files from the previous stage
COPY --from=build /app/migration .
COPY --from=build /app/server .

# Copy the entrypoint script
COPY entrypoint.sh .

# Make the entrypoint script executable
RUN chmod +x entrypoint.sh

# Expose the port the app runs on
EXPOSE 8080

# Run the entrypoint script
CMD ["./entrypoint.sh"]
