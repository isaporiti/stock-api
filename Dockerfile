FROM golang:1.20.6-alpine3.18

WORKDIR /app

COPY . .

RUN go build -o app cmd/main.go

EXPOSE 5001

CMD ["./app"]