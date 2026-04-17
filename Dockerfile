FROM golang:alpine AS builder

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Delete migrate.go to prevent 'main redeclared' build error
RUN rm -f migrate.go

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a minimal Alpine image for runtime
FROM alpine:latest
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the API port
EXPOSE 8080

# Command to run
CMD ["./main"]
