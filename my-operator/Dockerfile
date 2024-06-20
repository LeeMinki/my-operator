# Stage 1: Build the Go application
FROM golang:1.21 as builder

# Set the working directory inside the container
WORKDIR /workspace

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the operator binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager ./cmd/main.go

# Stage 2: Create the final image
FROM gcr.io/distroless/static:nonroot

# Set the working directory inside the container
WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /workspace/manager .

# Use non-root user
USER 65532:65532

# Command to run the operator
ENTRYPOINT ["/manager"]
