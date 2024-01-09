# Use the official Golang image as the base image
FROM golang:1.21.3-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Use a smaller base image for the final image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the port that the application will run on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
