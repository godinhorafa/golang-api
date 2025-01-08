# Produto API

Esta é uma API REST para gerenciamento de produtos, construída em Go. A API permite operações CRUD (Criar, Ler, Atualizar e Deletar) em produtos.

## Pré-requisitos

Antes de começar, você precisará ter instalado em sua máquina:

- [Go](https://golang.org/dl/) (versão 1.22.6 ou superior)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)

## Configuração

1. **Clone o repositório:**

   ````bash
   git clone https://github.com/seu_usuario/produto-api.git
   cd produto-api   ```

   ````

2. **Crie um arquivo `.env`:**

   Copie o arquivo de exemplo `.env.example` para `.env` e ajuste as variáveis de ambiente conforme necessário.

   ````bash
   cp .env.example .env   ```

   As variáveis de ambiente necessárias são:

   - `DB_HOST`: Endereço do banco de dados (padrão: `db`)
   - `DB_PORT`: Porta do banco de dados (padrão: `3306`)
   - `DB_USER`: Usuário do banco de dados (padrão: `root`)
   - `DB_PASSWORD`: Senha do banco de dados (padrão: `senha123`)
   - `DB_NAME`: Nome do banco de dados (padrão: `produtos`)

   ````

3. **Inicie o banco de dados e a API usando Docker Compose:**

   Para iniciar os serviços, execute o seguinte comando no terminal:

   `````bash
   docker-compose up --build   ```

   - O comando `docker-compose up` inicia os serviços definidos no arquivo `docker-compose.yml`.
   - A opção `--build` força a reconstrução das imagens, garantindo que você esteja usando a versão mais recente do código.

   **Observação:** Se você quiser executar os serviços em segundo plano (modo detached), adicione a opção `-d`:

   ````bash
   docker-compose up --build -d   ```

   Para parar os serviços, você pode usar:

   ````bash
   docker-compose down   ```

   Isso irá parar e remover os contêineres, redes e volumes criados pelo `docker-compose up`.
   `````

## Execução

Após iniciar os serviços, a API estará disponível em `http://localhost:8080`.

### Endpoints

- `GET /produtos`: Retorna todos os produtos.
- `GET /produtos/{id}`: Retorna um produto específico pelo ID.
- `POST /produtos`: Cria um novo produto.
- `PUT /produtos/{id}`: Atualiza um produto existente pelo ID.
- `DELETE /produtos/{id}`: Deleta um produto pelo ID.
- `POST /produtos/importar`: Importa produtos em massa a partir de um arquivo JSON ou CSV.

## Licença

Este projeto está licenciado sob a MIT License - veja o arquivo [LICENSE](LICENSE) para mais detalhes.
