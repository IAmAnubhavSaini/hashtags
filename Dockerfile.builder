FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy go files
COPY src/hashtags.go .
COPY go.mod ./

# Build the Go app
RUN go build -o /hashtags hashtags.go

# Final image
FROM alpine:latest  

WORKDIR /app

# Copy binary from builder
COPY --from=builder /hashtags /app/hashtags

# Create dist directory
RUN mkdir -p /app/dist

# Set entrypoint
ENTRYPOINT ["/app/hashtags"]

# By default, print help if no arguments are provided
CMD [] 