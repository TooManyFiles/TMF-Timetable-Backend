# Build Stage
FROM golang:1.22.1 AS builder

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .
# Download the Go module dependencies
RUN go mod download
# Build the Go application for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tmf-timetable-backend main.go


# Final Stage
FROM alpine:latest

# Install tzdata to get timezone information
RUN apk add --no-cache tzdata

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/tmf-timetable-backend .

# Expose the port your application runs on (adjust if necessary)
EXPOSE 8080
ENV CONFIG_FILE=./config/config.yml
# Command to run the executable
CMD ["./tmf-timetable-backend"]
