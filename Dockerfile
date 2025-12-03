# Build stage
FROM golang:1.24-alpine AS builder

# Install ca-certificates (needed for TLS/HTTPS)
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary (static, stripped, minimal)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o andi-custodian ./cmd/demo

# Final stage
FROM scratch

# Import certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Create non-root user (security best practice)
# (scratch has no user management, so we rely on runtime user)
# In production, set USER 65532 (nonroot) if using distroless

# Copy binary
COPY --from=builder /app/andi-custodian /andi-custustodian

# Set entrypoint
ENTRYPOINT ["/andi-custodian"]