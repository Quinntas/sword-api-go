FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/.env .

USER appuser

EXPOSE 3000

CMD ["./main"]