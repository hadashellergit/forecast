FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /kfc-server ./cmd/server

FROM alpine:3.19

RUN apk --no-cache add ca-certificates postgresql16-client

WORKDIR /app

COPY --from=builder /kfc-server .
COPY config.yaml .
COPY migrations/ ./migrations/
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]