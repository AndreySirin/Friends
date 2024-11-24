FROM golang:1.23.1 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o myapp


FROM alpine:latest

WORKDIR /root

COPY --from=builder /app/myapp .

EXPOSE 8080

CMD ["./myapp"]