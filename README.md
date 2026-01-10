# Desafio Pós Go Expert - Rate limiter

> Este projeto contém a solução para o desafio técnico sobre rate limit da pós-graduação Go Expert da FullCycle.

(...) Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:

- Endereço IP: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
- Token de Acesso: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:

    - API_KEY: `<TOKEN>`

- As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.

Requisitos:

- O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
- O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.
- O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
- As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.
- Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.
- O sistema deve responder adequadamente quando o limite é excedido:
    - Código HTTP: 429
    - Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame
- Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.
- Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.
- A lógica do limiter deve estar separada do middleware.

Exemplos:

- Limitação por IP: Suponha que o rate limiter esteja configurado para permitir no máximo 5 requisições por segundo por IP. Se o IP 192.168.1.1 enviar 6 requisições em um segundo, a sexta requisição deve ser bloqueada.
- Limitação por Token: Se um token abc123 tiver um limite configurado de 10 requisições por segundo e enviar 11 requisições nesse intervalo, a décima primeira deve ser bloqueada.
- Nos dois casos acima, as próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Ex: Se o tempo de expiração é de 5 minutos, determinado IP poderá realizar novas requisições somente após os 5 minutos.

Dicas:

- Teste seu rate limiter sob diferentes condições de carga para garantir que ele funcione conforme esperado em situações de alto tráfego.

Entrega:

- O código-fonte completo da implementação.
- Documentação explicando como o rate limiter funciona e como ele pode ser configurado.
- Testes automatizados demonstrando a eficácia e a robustez do rate limiter.
- Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
- O servidor web deve responder na porta 8080.

# Funcionalidades

- Limite de requisições por IP ou por token de acesso
- Configuração de limites e tempos de bloqueio via variáveis de ambiente ou arquivo .env.
- Middleware para fácil integração com servidores HTTP
- Persistência dos dados de limitação utilizando Redis
- Resposta com código HTTP 429 quando o limite é excedido

# Executando a aplicação

### Configurando as variáveis de ambiente

1. Crie o arquivo de configurações, após baixar o projeto.

    ```sh
    cp .env.example .env
    ```

2. Edite as variáveis de ambiente no arquivo `.env`

    ```yaml
    REDIS_HOST=redis:6379 # Endereço onde Redis esta rodando dentro do docker
    REDIS_PASSWORD="123"
    RATE_LIMIT_REQUESTS_PER_SECOND_IP=100 # Quantidade de requisições permitidas por segundo por um mesmo IP
    RATE_LIMIT_REQUESTS_PER_SECOND_TOKEN=10 # Quantidade de requisições permitidas por segundo por um mesmo token
    ```

### Iniciando os serviços

1. Inicie os containers através do docker compose:

    ```sh
    #*lembre de configurar o .env antes de iniciar os containers
    docker-compose up -d
    ```

### Testando a aplicação

1. A aplicação roda na porta 8080 e expoe o endpoint `/sample` para testes:

    ```sh
    curl http://localhost:8080/sample

    ## ou passando uma `API_KEY`
    curl http://localhost:8080/sample -H 'API_KEY: somevalue'
    ```

    O sistema verificará a quantidade de acessos e aplicará as regras de limitação conforme configurado.

2. Testes automatizados

    ```sh
    go test ./test/api_test.go
    ```
