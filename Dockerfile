FROM golang:1.24.0
WORKDIR /app
ADD . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/faas_server ./cmd/server

FROM debian:12-slim
RUN apt update && apt install -y ca-certificates
COPY --from=0 /app/faas_server /app/faas_server
CMD ["/app/faas_server"]