# Estágio 1: compilar o binário
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api ./cmd/main.go

# Estágio 2: imagem final mínima
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]
