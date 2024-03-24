ARG VERSION

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -ldflags "-X 'main.Version=$VERSION'"  -o main .

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main /app/main

EXPOSE 6789

CMD ["./main","--config=/app/data/config.json","--basedirectory=/app/data/","--databasepath=/app/data/database.db"]
