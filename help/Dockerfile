FROM golang:1.24.0-alpine
WORKDIR /app
ADD . /app
RUN go mod download
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" ./cmd/ -o /app/faas_server

FROM alpine:latest
EXPOSE 8080
CMD ["/app/faas"]