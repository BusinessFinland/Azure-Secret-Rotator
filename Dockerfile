# Build stage
FROM golang:1.23 AS builder

# Set the working directory in the builder stage
WORKDIR /app

# Copy go.mod and go.sum files for dependency management
COPY go.mod go.sum ./
RUN go mod download

# Copy the Go source files, including helpers and apis
COPY *.go ./
COPY helpers/*.go ./helpers/
COPY apis/*.go ./apis/

# Build the Go binary in a static way so that it has no dependencies on the
# C libraries and can run in a scratch container or Alpine.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /azure-secret-rotator .

# Final stage
FROM alpine:latest

# Install ca-certificates and cronie for cron functionality
RUN apk add --no-cache ca-certificates cronie

# Set the working directory
WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /azure-secret-rotator /azure-secret-rotator

# Copy the wrapper script and make it executable
COPY run_rotator.sh /run_rotator.sh
RUN chmod +x /run_rotator.sh

# Copy the entrypoint script and make it executable
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Copy necessary files with appropriate permissions
COPY azuresecretrotator.private-key.pem app/azuresecretrotator.private-key.pem
RUN chmod 0644 app/azuresecretrotator.private-key.pem

COPY crontab /etc/cron.d/azure-secret-rotator
RUN chmod 0644 /etc/cron.d/azure-secret-rotator && crontab /etc/cron.d/azure-secret-rotator

# (Optional) Expose the application port if needed
EXPOSE 8080

USER root

# Set the entrypoint to the entrypoint script
ENTRYPOINT ["/entrypoint.sh"]
