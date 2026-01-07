# Etapa 1: Build (Compilação)
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instala dependências do sistema necessárias para compilar (se houver)
RUN apk add --no-cache git

# Copia os arquivos de dependência primeiro (para aproveitar cache do Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o binário
# Aponta para o main.go que está em cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/server/main.go

# Etapa 2: Runtime (Imagem final)
FROM alpine:latest

WORKDIR /app

# Instala certificados de segurança (para fazer chamadas HTTPS externas, ex: Piston API)
RUN apk --no-cache add ca-certificates

# Copia o binário da etapa anterior
COPY --from=builder /app/api .

# Copia arquivos adicionais se necessário (ex: migrations sql)
# COPY --from=builder /app/schema.sql . 

# Expõe a porta que sua aplicação usa (conforme seu config.go)
EXPOSE 8080

# Comando para rodar
CMD ["./api"]