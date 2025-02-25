FROM golang:1.17-buster as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build

FROM debian:buster-slim

WORKDIR /app

COPY --from=builder /app/probed /app

ENTRYPOINT ["/app/probed"]
