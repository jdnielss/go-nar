# Use a multi-architecture base image
FROM --platform=$BUILDPLATFORM golang:latest AS builder

# Set necessary environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY . .

# Download and install Go modules
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o app .

# Final stage: the running container
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /app/app /

# Set the entry point for the container
ENTRYPOINT ["/app"]

