FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o fs-mcp .

FROM alpine
WORKDIR /app
COPY --from=builder /app/fs-mcp /app/fs-mcp

# Create the directory that will be mounted
RUN mkdir -p /baseDir

ENTRYPOINT ["/app/fs-mcp", "--docker", "--dir", "/baseDir"]