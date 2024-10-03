FROM golang:1.22.0-alpine as builder

WORKDIR /build
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go mod download
RUN go build -o ./notification_service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/notification_service /app/
COPY --from=builder /build/.env /app/

CMD ["./notification_service"]