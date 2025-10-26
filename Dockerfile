FROM golang:1.25-alpine AS build
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod go mod download
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/faas_server ./cmd/server

FROM alpine:latest
COPY --from=build /app/faas_server /app/faas_server
RUN apk add --no-cache ca-certificates
CMD ["/app/faas_server"]
