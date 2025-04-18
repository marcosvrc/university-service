FROM golang:1.21-alpine

WORKDIR /app

# Instalar dependências do sistema
RUN apk add --no-cache git

# Copiar arquivos do módulo Go
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar o código fonte
COPY . .

# Compilar a aplicação
RUN go build -o main .

# Expor a porta
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"]