FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_PATH

RUN go build -o /app/main $SERVICE_PATH

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache tzdata
ENV TZ=Europe/Istanbul
# -----------------------

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080 8081 8082 3000 50051

CMD ["./main"]