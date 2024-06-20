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
definido em `AUCTION_INTERVAL` no arquivo `.env`. \

#### Note
Todas as alterações para adicionar a lógica de fechamento do leilão
foram feitas no arquivo `internal/infra/database/auction/create_auction.go`
