# Use an official Go runtime as a parent image
FROM golang:1.18 as builder

# Set the working directory inside the container
WORKDIR /app

# Pre-cache dependencies before main build to speed up builds
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the current directory contents into the container at /app
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main .

# Use a Docker multi-stage build to create a lean production image
# Start from scratch (empty) image
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Add CA certificates to allow SSL-based applications
RUN apk --no-cache add ca-certificates

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
