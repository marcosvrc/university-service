# University Service

Este é um serviço RESTful para gerenciamento de universidades, construído com Go, MongoDB e Kafka.

## Tecnologias Utilizadas

- Go 1.21
- MongoDB
- Apache Kafka
- Docker & Docker Compose
- Gin Web Framework

## Estrutura do Projeto

```
university-service/
├── api/
│   └── handlers.go
├── config/
│   ├── config.go
│   └── config.yaml
├── internal/
│   ├── models/
│   │   └── university.go
│   ├── repository/
│   │   └── mongodb.go
│   └── service/
│       └── kafka.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Pré-requisitos

- Docker
- Docker Compose
- Go 1.21 (para desenvolvimento local)

## Configuração

1. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/university-service.git
cd university-service
```

2. Configure as variáveis de ambiente no arquivo `config/config.yaml` (se necessário)

## Executando com Docker

1. Construa e inicie os containers:
```bash
docker-compose up -d
```

2. O serviço estará disponível em `http://localhost:8080`

## Desenvolvimento Local

1. Instale as dependências:
```bash
go mod download
```

2. Execute o serviço:
```bash
go run main.go
```

## API Endpoints

### Criar Universidade
```http
POST /universities
Content-Type: application/json

{
    "name": "Universidade Example",
    "address": "Rua Example, 123",
    "phone": "(11) 1234-5678",
    "email": "contato@example.edu",
    "website": "https://www.example.edu"
}
```

### Buscar Universidade por ID
```http
GET /universities/{id}
```

### Listar Todas as Universidades
```http
GET /universities
```

### Atualizar Universidade
```http
PUT /universities/{id}
Content-Type: application/json

{
    "name": "Universidade Example Atualizada",
    "address": "Rua Example, 456",
    "phone": "(11) 8765-4321",
    "email": "novo.contato@example.edu",
    "website": "https://www.example.edu"
}
```

### Deletar Universidade
```http
DELETE /universities/{id}
```

## Eventos Kafka

O serviço publica os seguintes eventos no tópico `university_events`:

- `university_created`: Quando uma nova universidade é criada
- `university_updated`: Quando uma universidade é atualizada
- `university_deleted`: Quando uma universidade é deletada

## Estrutura do Evento

```json
{
    "type": "university_created|university_updated|university_deleted",
    "university": {
        "id": "ObjectID",
        "name": "string",
        "address": "string",
        "phone": "string",
        "email": "string",
        "website": "string",
        "created_at": "timestamp",
        "updated_at": "timestamp"
    }
}
```

## Monitoramento

- MongoDB está disponível em `localhost:27017`
- Kafka está disponível em `localhost:9092`
- Zookeeper está disponível em `localhost:2181`

## Contribuindo

1. Fork o projeto
2. Crie sua feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para mais detalhes.