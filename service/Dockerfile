FROM golang:latest AS builder
WORKDIR /app

COPY  go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main/ .
COPY ./CDN /root/CDN

EXPOSE 3000
CMD ["./main"]