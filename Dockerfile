FROM golang:1.21.0-alpine3.18 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -ldflags "-X 'main.Version=Docker'"  -o main .

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main /app/main

EXPOSE 6789

CMD ["./main","--config=/app/data/config.json","--basedirectory=/app/data/"]
