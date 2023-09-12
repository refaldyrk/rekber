FROM golang:latest

WORKDIR /app

COPY .env .

COPY . .

RUN go build -o main .

EXPOSE 9090

CMD ["./main"]