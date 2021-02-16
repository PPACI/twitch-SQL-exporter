FROM golang:1.15 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg pkg
COPY main.go .
RUN go build -o twitch-sql-exporter .

FROM debian:buster-slim

COPY --from=builder /app/twitch-sql-exporter /usr/local/bin/twitch-sql-exporter

ENTRYPOINT ["twitch-sql-exporter"]