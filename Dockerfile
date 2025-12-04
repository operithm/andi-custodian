# Build stage
FROM golang:1.24-alpine AS builder

# Install CA certs
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download modules
RUN go mod download

# Copy ALL source (including generated .pb.go files!)
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -trimpath -o custody-server ./cmd/server

# Final stage
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/custody-server /custody-server

ENTRYPOINT ["/custody-server"]