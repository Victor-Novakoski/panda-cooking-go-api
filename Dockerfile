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
ENTRYPOINT ["/api"]
