FROM golang:1.23.1 AS builder

WORKDIR /app

COPY . ./

RUN go mod tidy

WORKDIR /app/cmd

RUN go build -o /app/friend


FROM ubuntu:latest

WORKDIR /root

COPY --from=builder /app/friend .

COPY config.yaml /root/

COPY internal/dirMigrite /app/internal/dirMigrite

COPY internal/htmlFile /app/internal/htmlFile

EXPOSE 8080

CMD ["./friend"]