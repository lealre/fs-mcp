FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o fs-mcp .

FROM alpine
WORKDIR /app
COPY --from=builder /app/fs-mcp /app/fs-mcp

ENV FS_MCP_DOCKER_MODE=true
# Hardcode Docker mode but allow volume override
ENTRYPOINT ["/app/fs-mcp", "-t","http"]