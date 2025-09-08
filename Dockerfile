FROM golang:1.26-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/faas_server ./cmd/server

FROM debian:bullseye-slim
COPY --from=build /app/faas_server /app/faas_server
CMD ["/app/faas_server"]
