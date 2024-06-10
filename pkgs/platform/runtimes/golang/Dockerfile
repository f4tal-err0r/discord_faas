
FROM golang:1.22-alpine
WORKDIR /app
ADD . /app
RUN go mod download
RUN make test
RUN make build

FROM alpine:latest
LABEL runtime="golang"
WORKDIR /app
COPY --from=0 /app/gofunc /app/gofunc
CMD ["/app/route-generator"]