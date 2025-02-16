FROM golang:1.21.7-alpine
WORKDIR /app
ADD . /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" ./cmd/ -o /app/faas_server

FROM alpine:latest
CMD ["/app/faas_server"]