FROM golang:1.23.6-alpine AS builder

RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .

RUN go build -o main ./cmd/logfinder

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY .env .

ENV SERVER_PORT=8080
EXPOSE $SERVER_PORT
CMD ./main