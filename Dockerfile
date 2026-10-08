# Desenvolvimento: código montado como volume e hot reload com air (docker compose up)
FROM golang:1.27-alpine AS dev
RUN go install github.com/air-verse/air@v1.67.4
WORKDIR /app
CMD ["air"]

# Estágio 1: compilar o binário
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd

# Estágio 2: imagem final mínima, sem shell e rodando sem root
FROM gcr.io/distroless/static-debian12:nonroot AS prod
COPY --from=build /out/api /api
EXPOSE 8080
# a imagem não tem curl: o próprio binário confere o /health
HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 CMD ["/api", "healthcheck"]
ENTRYPOINT ["/api"]
