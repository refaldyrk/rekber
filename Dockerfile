FROM golang:1.20-alpine as builder

WORKDIR /app

COPY . .

# RUN go mod download

RUN go build -o main .

FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app .

EXPOSE 9090

CMD ["./main", "--type=docker"]