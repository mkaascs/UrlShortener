FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o url-shortener \
    ./cmd/url-shortener/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/url-shortener .
COPY --from=builder /app/config ./config
COPY --from=builder /app/docs ./docs

EXPOSE 5055

ENTRYPOINT ["./url-shortener"]