# Use the official Golang image as the builder
FROM golang:1.20.3-alpine as builder

# Enable CGO to use C libraries (set to 0 to disable it)
# We set it to 0 to build a fully static binary for our final image
ENV CGO_ENABLED=0

# Set the working directory
WORKDIR /app

# Copy the Go Modules manifests (go.mod and go.sum files)
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application and install it
RUN go install .

# Use a scratch image as the final distroless image
FROM scratch

# Set the working directory
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /go/bin/ChatGPT-Proxy-V4 /app/ChatGPT-Proxy-V4

# Expose the port where the application is running
EXPOSE 8080

# Start the application
CMD [ "./ChatGPT-Proxy-V4" ]
