# Auction system

Esse repositório adiciona a funcionalidade de fechamento automático do leilão.
Repositório base: https://github.com/devfullcycle/labs-auction-goexpert

### Iniciando o projeto
Para rodar o servidor execute o comando `docker-compose up --build` que atualiza as imagens
e inicia as dependênicias junto com a aplicação.

### Sistema de fechamento do leilão
Todas as `auctions` criadas iniciam uma goroutine responsável
por atualizar o status para `Completed`. \
O status é autalizado após o tempo de duração do leilão,
definido em `AUCTION_INTERVAL` no arquivo `.env`. 

### Testando a aplicação
1. Crie uma nova `auction`, exemplo de request:
```http request
POST http://localhost:8080/auction
Content-Type: application/json

{
  "product_name": "car",
  "category": "vehicles",
  "description": "super fast car",
  "condition": 0
}
```

2. Para ver todos os detalhes da action utilize o seguinte endpoint:
```http request
GET http://localhost:8080/auction/f671cc31-8943-474f-b5ba-d8090fd014db
```

3. Após o tempo definido no arquivo `.env` a auction terá o status atualizado para `1`, o que significa `Completed`.

#### Note
Todas as alterações para adicionar a lógica de fechamento do leilão
foram feitas no arquivo `internal/infra/database/auction/create_auction.go`
