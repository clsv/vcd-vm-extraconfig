# Stage 1: Build the Go binary
FROM golang:1.24.2-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum (if using modules)
# If not using modules, skip these lines or initialize a module
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY vcd-vm-extraconfig.go .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -s' -o vcd-vm-extraconfig vcd-vm-extraconfig.go

# Stage 2: Create minimal runtime image
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /app/vcd-vm-extraconfig /app/vcd-vm-extraconfig

# Set working directory
WORKDIR /app

# Command to run the binary
ENTRYPOINT ["/app/vcd-vm-extraconfig"]
