# Documentação Payment Processor API

Descrição da solução para o desafio do Processador de Pagamentos em Golang([Back-end](https://github.com/pagarme/vagas/tree/master/desafios/software-engineer-golang))

- [O Problema](#O-Problema)
- [Back-end](#Back-end)
    - [Dependências](#Dependências)
        - [Boilerplate Code](#Boilerplate-Code)
- [Documentação API](#Documentação-API)
    - [Segurança](#Segurança)
- [Serviço Adquirentes](#Serviço-Adquirentes)
- [Instalação e Execução](#Instalação-e-Execução)
- [Ambiente Desenvolvimento](#Ambiente-Desenvolvimento)
    - [Instalação das Dependências](#Instalação-das-Dependências)
    - [Configurações](#Configurações)
    - [Execução](#Execução)
    - [Infra Desenvolvimento](#Infra-Desenvolvimento)
    - [Infra Testes](#Infra-Testes)
    - [Infra Documentação](#Infra-Documentação)
- [Testes](#Testes)
    - [Feedback rápido](#Feedback-rápido)
        - [Unitários](#Unitários)
        - [Integração](#Integração)
        - [Unitários e Integração](#Unitários-e-Integração)
    - [Carga](#Carga)
        - [Jmeter](#Jmeter)
        - [ApacheBench](#ApacheBench)
- [Empacotamento](#Empacotamento)
- [Instrumentação](#Instrumentação)
    - [Links Úteis](#Links-Úteis)
    - [Grafana](#Grafana)
- [Comandos Make](#comandos-make)
- [Melhorias](#Melhorias)

## O Problema

Foi implementada solução que permite fazer processamento de solicitação de autorização de pagamento, em colaboração com Adquirentes.

Processador foi implementado como API Rest que recebe requisições em um endpoint, enriquece requisição com dados sensíveis de cartão armazenados em banco de dados Redis, e depois encaminha requisição para serviço Rest da Adquirente.

A seguir temos exemplo de utilização da API. Para uma documentação mais completa verifique o tópico [swagger](#Documentação-API).

 - Solicitação de Autorização
 
    Ex: Autorizar compra do *João* de valor *1000,00* com *1* parcela, usando cartão *xpto121a*, através da Adquirente *Stone*.

    ```sh
    curl -X POST \
        http://localhost:80/api/transactions \
        -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjoicGF5bWVudC1wcm9jZXNzb3ItYXBpIn0.uw-8pECPeJbme82nptMI-bsP8f4GvCx9x6b_GzM5wws' \
        -H 'X-ACQUIRER-ID: Stone' \
        -d '{
            "token": "xpto121a",
            "holder": "João",
            "total": 1000,
            "installments": 1
        }'
    ```

    Teremos a resposta JSON produzida na Adquirente que a transação foi Aprovada.

    ```json
    {
        "message": "Transação autorizada"
    }
    ```

## Back-end

Implementação do Processador como microserviço e API Rest.

Aplicação tem um endpoint HTTP Rest */transactions*.
Foi definida uma camada de negócio no pacote internal/pkg/***`processor`*** que tem 
um ponto de entrada, `PaymentProcessorService`. 

Uma chamada HTTP típica ao endpoint /transactions para Autorização segue este caminho, iniciando pelo método [`Controller.Process()`](internal/pkg/processor/api/api.go) da API.

- `Controller` recebe requisição HTTP, faz algumas validações e aciona camada de negócio.

- `PaymentProcessorService` do negócio, colabora com `AcquirerActors` para enviar a transação para os Workers do Adquirente correto(Fan-in).

- Na sequência algum Worker `Acquirer` recebe a transação(Fan-out), a enriquece com informações sensíveis armazenadas no Redis, e depois envia esses dados para o Serviço da Adquirente por meio de outra chamada HTTP.

- Por último a resposta da Adquirente de Autorizada ou Negada é devolvida pela API. O Status Code do HTTP é 
usado como indicativo de Autorizada(200) ou Negada(400).

A seguir é representado as interações entre as interfaces que os componentes implementam. `Controller` é o único participante concreto.

![Diagrama de Sequência do Processamento de uma Transação](docs/ProcessTransaction.png "Processar Transação")

Segue descrição dos principais pacotes e arquivos da solução.

 ```sh
 tree -L 5
.
├── cmd
│   └── processorAPI
│       ├── main.go                 # Main da Aplicação
│       ├── wire_gen.go             # Boilerplate Code do Framework de Injeção de Dependências
│       └── wire.go                 # Provedores do Framework de Injeção de Dependências da aplicão
├── configs
│   ├── config.env                  # Arquivo ENV de configuração usado em 'make run'
├── go.mod                          # Dependências da aplicação
├── internal
│   ├── app
│   │   ├── app.go                  # Inicialização da Aplicação
│   │   └── wireSets.go             # Declaração dos componentes manipulados pelo Framework de Injeção de Dependências
│   └── pkg
│       ├── commom                  # Package de utilitários
│       ├── infra                   # Package de infraestrutura Redis, Http, JWT
│       ├── instrumentation         # Package 'DevOps'
│       │   ├── info
│       │   └── metrics
│       ├── processor               # Package principal do contexto de negócio
│       │   ├── acquireractors.go   # Faz roteamento para Acquirer correto
│       │   ├── acquireractors_test.go
│       │   ├── acquirer.go         # Serviço processamento, que enriquece e envia dados para Adquirentes
│       │   ├── acquirer_test.go
│       │   ├── api                 # Package da API HTTP da aplicação
│       │   │   ├── api.go
│       │   │   ├── api_it_test.go   # Testes de integração
│       │   │   ├── api_test.go     # Testes de unidade
│       │   │   ├── errors.go       # Mapeia erros de negócio para erros HTTP Status Code
│       │   ├── card.go             # Entidade com dados sensíveis persitida no Redis
│       │   ├── dto.go              # Implementa DTO usado para recebimento de dados
│       │   ├── errors.go           # Erros de negócio
│       │   ├── processor.go        # Ponto de entrada do negócio. Atende solicitações da API
│       │   ├── processor_test.go   # Testes de unidade com Mock Objects
│       │   ├── repository.go       # Meio de Persitência com Redis
│       │   ├── repository_it_test.go
│       └── webserver
│           ├── home.go             # Gera HTML de página símples com links úteis
│           ├── webserver.go        # Implementa servidor HTTP que expões API de negócio
 ```

### Dependências

- **[Go](http://golang.org/)** - Liguagem usada na implementação da API. 
- **[Wire](http://github.com/google/wire)** - Framework de Injeção de Dependências.
    - Com Wire aplicação não quebra em runtime por falha na declaração de dependências.
    - Ciclo de desenvolvimento é mais rápido por não precisar iniciar aplicação para testar dependências.
    - Mensagens de falha na resolução do grafo de dependências são claras.
    - Ponto desfavorável: Gera arquivos Boilerplate Code ***`wire_gen.go`*** que devem ser commitados.
- **[Redigo](http://github.com/gomodule/redigo)** - Driver performárico, que mantém API do Redis.
- **[slLog](http://github.com/kimprado/sllog)** - Escrevi esta lib para configurar logging como no Spring Boot.
- **[Configor](http://github.com/jinzhu/configor)** - Lib flexível para carregar configuração, via Variáveis de Ambiente e outros.
- **[HttpRouter](http://github.com/julienschmidt/httprouter)** - HTTP mux performático e flexível.
- **[Prometheus](http://github.com/julienschmidt/httprouter)** - Excelente ferramenta para publicar métricas da aplicação.
- **[Testify](http://github.com/stretchr/testify)** - Lib que uso para simplificar assertions.

#### Boilerplate Code

A dependência Wire gera arquivos Boilerplate Code.

- `wire_gen.go`
- `wire_gen_test.go` (renomeados por script)

## Documentação API

Documentação disponibilizada no arquivo [`swagger.yml`](api/swagger.yml). 
Para acessar a documentação interativa execute o ambiente como descrito a seguir([Instalação e Execução](#Instalação-e-Execução)), e depois siga as instruções em [Infra Documentação](#Infra-Documentação).

Consultas no formato Postman estão no arquivo [`ProcessTransactions.postman_collection.json`](tools/ProcessTransactions.postman_collection.json).

### Segurança

Autorização da API é feita com [JWT](http://jwt.io/) token.
Deve ser enviado Header HTTP Authorization como a seguir.
```yaml
Authorization: Bearer <jwt_token>
```

Exemplo.

```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjoicGF5bWVudC1wcm9jZXNzb3ItYXBpIn0.uw-8pECPeJbme82nptMI-bsP8f4GvCx9x6b_GzM5wws
```

A chave JWT '***challenge***' usada nos testes e em desenvolvimento é parametrizada via 
variável de ambiente. 
Nos testes o Container Docker recebe variáveis de ambiente definidas no arquivo
de configuração config.env.
```INI
PROCESSOR_SECURITY_JWTKEY=challenge
```

Para criar um token válido é preciso informar o atributo ***aud** = payment-processor-api*.
```json
{
"sub": "1234567890",
"aud": "payment-processor-api"
}
```

Para verificar o token acima use a chave 'challenge' no [link](https://jwt.io/#debugger-io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjoicGF5bWVudC1wcm9jZXNzb3ItYXBpIn0.uw-8pECPeJbme82nptMI-bsP8f4GvCx9x6b_GzM5wws).

## Serviço Adquirentes

Foi criado serviço para simular Adquirentes. Implementado em go([`main.go`](test/_acquirers_server_mock/main.go)), responde requisições HTTP. 
Foram simuladas apenas Stone([`stone.go`](test/_acquirers_server_mock/stone.go)) e Cielo([`cielo.go`](test/_acquirers_server_mock/cielo.go)).

Para ficar um pouco mais real foi aplicado um delay(parametrizável) de 100ms, em todas as requisições recebidas.

 - Configurações

```INI
ACQUIRERS_PORT=8092
ACQUIRERS_DELAY=100
ACQUIRERS_LOGGING=DEBUG
```

Transações Autorizadas retornam HTTP Status Code *200*, e Negadas *400*, usados pelo [Processador](#Back-end).

 - Todas as transações são autorizadas, exceto os casos descritos a seguir.

    Adquirente | Campo | Comp. | Valor
    --- | ---: | :---: | :---
    Stone | Total | > | 1000
    Stone | Portador | == |  João Antônio
    Stone | Parcelas | > |  12
    Cielo | Total | > | 500
    Cielo | Portador | == |  João Antônio
    Cielo | Parcelas | > |  6

O serviço é iniciado pelo Docker Compose juntamente com o Processador, e demais dependências.

## Instalação e Execução

Para fazer deploy e execução do projeto rode os seguintes comandos.

```sh
./configure
```

Após executar feche e abra outro terminal. Alguns comandos talvez ainda precisem ser executados com sudo, ou reinicie o computador.

Execute aplicaçao com comando

```sh
make run
```

Ao final o comando printa no console as URLs para acesso aos serviços.

- http://localhost:80/              (nginx)      - Página web com links úteis
- http://0.0.0.0:3000/              (API)        - URL da API
- http://localhost:80/docs          (swagger)    - Documentação interativa
- http://localhost:3001/d/kKd-m3qiz (grafana)    - URL do Grafana
- http://localhost:9090             (prometheus) - URL do Prometheus

Exemplo:

```sh
make run
...
Acesse nginx:
http://localhost:80/
Acesse API:
http://0.0.0.0:3000/
Acesse swagger:
http://localhost:80/docs
Acesse grafana:
http://localhost:3001/d/kKd-m3qiz/requisicoes-http-payment-processor-api?orgId=1&kiosk=tv
Acesse prometheus:
http://localhost:9090
...
```

## Ambiente Desenvolvimento

Segue como instalar e configurar o ambiente e ferramentas de desenvolvimento do projeto.

### Instalação das Dependências

Execute script [`configure`](configure) presente na raiz do repositório para fazer download e instalação das dependências. O script também cria, caso necessário, a configuração da IDE e de execução da aplicação.

```sh
./configure
```

Após executar feche e abra outro terminal. Alguns comandos talvez ainda precisem ser executados com sudo, ou reinicie o computador.

As seguintes ferramentas serão provisionadas, caso necessário.

- **Docker** - Ferramenta usada para containerização.
- **Docker Compose** - Ferramenta usada para orquestração em ambiende de dev.
- **Go** - Linguagem de programação.
- **Wire** - Framework de Injeção de Dependências.
- **Jmeter** - Ferramenta usada para testes de carga. Testes podem ser executados a partir de containers, e também localmente.

Os seguintes arquivos são criados, caso necessário.

- .vscode/settings.json - Arquivo da IDE VSCode
- .vscode/launch.json - Arquivo da IDE VSCode
- configs/config.env - Configurações injetadas como Variáveis de Ambiente no deploy com [Docker Compose](#Execução).
- configs/config-dev.env - Configurações padão de desenvolvimento.
- configs/config-dev.json - Configuração opcional da aplicação em tempo de desenvolvimento.

### Configurações

Para execução containerizada as Variáveis de Ambiente da aplicação devem ser declaradas no arquivo [`config.env`](configs/config.env), criado no passo anterior.

```INI
PROCESSOR_ENVIRONMENT_NAME="local"
PROCESSOR_SERVER_PORT="3000"

PROCESSOR_REDISDB_HOST="redisdb"
PROCESSOR_REDISDB_PORT=6379

PROCESSOR_STONEACQUIRER_URL="http://acquirers:8092/stone"
PROCESSOR_STONEACQUIRER_CONCURRENTWORKERS=200

PROCESSOR_CIELOACQUIRER_URL="http://acquirers:8092/cielo"
PROCESSOR_CIELOACQUIRER_CONCURRENTWORKERS=200

# PROCESSOR_SECURITY_ENABLE=false
PROCESSOR_SECURITY_JWTKEY=challenge

PROCESSOR_LOGGING_LEVEL={ROOT: INFO, infra.redis.db: INFO, webserver: INFO}

```

### Execução

 - Executar solução

    ```sh
    make run
    ```

 - Interromper a execução

    ```sh
    make stop
    ```

### Infra Desenvolvimento

- Iniciar infra de Desenvolvimento

    ```sh
    make infra-start
    ```

- Interromper infra de Desenvolvimento

    ```sh
    make infra-stop
    ```

### Infra Testes

- Iniciar infra de Testes

    ```sh
    make infra-test-start
    ```

- Interromper infra de Testes

    ```sh
    make infra-test-stop
    ```

### Infra Documentação

- http://localhost:80/docs - Modo interativo com back-end dockerizado.

    ```sh
    make run
    ```

- http://localhost:8080/ - Somente documentação.

    ```sh
    docker-compose up -d swagger
    ```

- http://localhost:80/docs - Modo interativo com back-end rodando pela IDE, por exemplo.

    ```sh
    make infra-start
    ```

## Testes

Foram implementados diversos tipos de testes, divididos em Feedback rápido e Carga.

### Feedback rápido

Estão dividos em *[Unitários](#Unitários)* e *[Integração](#Integração)*. Os grupos de testes são separados em arquivos de testes diferentes. Usei o conceito de [Build Constraints ou Build Tag](http://golang.org/pkg/go/build/#hdr-Build_Constraints) para selecionar quais testes queremos executar.

Para especificar um grupo de teste executamos o comando `go test` com o parâmetro `-tags`.

```sh
go test ./internal/pkg/commom/config -tags="unit"
```

Neste exemplo o pacote *config* possui os seguintes arquivos de teste:

- [`config_test.go`](internal/pkg/commom/config/config_test.go)

    ```go
    // +build test unit

    package config
    // ...
    ```

- [`config_envvars_it_test.go`](internal/pkg/commom/config/config_envvars_it_test.go)

    ```go
    // +build test integration

    package config
    // ...
    ```

- [`config_it_test.go`](internal/pkg/commom/config/config_it_test.go)

    ```go
    // +build test integration

    package config
    // ...
    ```

Apenas os testes do arquivo [`config_test.go`](internal/pkg/commom/config/config_test.go), com a build tag *"// +build test **unit**"*, serão executados pois no comando informamos *-tags="**unit**"*.

#### Unitários

Testes unitários que não dependem da infra para executar, são mais rápidos, podendo conter Mock Objects conforme necessário.

Use os seguintes comandos para executar os testes unitários.

- Ambiente Containerizado.

    ```sh
    make test-unit-container
    ```

- Ambiente Local.

    ```sh
    make test-unit
    ```

Estes comandos são atalho para a execução do script [`test.sh`](scripts/test.sh) com parâmetro *unit*, que resulta em:

```sh
go test ./... -tags="unit" -cover -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -n 1
```

Para configurar um arquivo como Unit Test:

- Sufixo - *_test.go
- Build Tag - *unit*
    - Ex: // +build test unit
    - Ex: arquivo [`config_test.go`](internal/pkg/commom/config/config_test.go)
    - Ex: arquivo [`processor_test.go`](internal/pkg/processor/processor_test.go)

#### Integração

Testes de integração dependem do [deploy da infra de testes](#Infra-Testes). Acessam os serviços de dependência sem Mock Objects. Procuramos acelerar sua execução habilitando o paralelismo com *t.Parallel()*, e para permitir isso cada teste tem seu próprio prefixo nas chaves do Redis. 

Chaves do Redis mudam de acordo com o deploy. Ex: ***[prefixo]:card:hash_id_card***

- Em teste recebe prefixo *TestFindCards* e *TestIntegrationProcessTransaction*.

    ```
    TestFindCards:card:xpto121a
    TestIntegrationProcessTransaction:card:xpto121a
    ```

- Em produção recebe prefixo *processor*.

    ```
    processor:card:xpto121a
    ```

Use os seguintes comandos para executar os testes de integração.

- Ambiente Containerizado.

    ```sh
    make test-integration-container
    ```

- Ambiente Local.

    ```sh
    make test-integration
    ```

Estes comandos são atalho para a execução do script [`test.sh`](scripts/test.sh) com parâmetro *integration*, que resulta em:

```sh
go test -parallel 10 -timeout 1m30s ./... -tags="integration" -cover -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -n 1
```

Para configurar um arquivo como Integration Test:

- Sufixo - *_it_test.go
- Build Tag - *integration*
    - Ex: // +build test integration
    - Ex: arquivo [`api_it_test.go`](internal/pkg/processor/api/api_it_test.go)
    - Ex: arquivo [`repository_it_test.go`](internal/pkg/processor/repository_it_test.go)

#### Unitários e Integração

Permite executar ao mesmo tempo testes de [Unidade](#Unitários) e de [Integração](#Integração). O benefício é ter maior cobertura e estatística unificada.

Os testes devem ser configurados com a build tag ***test***, sendo:

- Unitários
    ```go
    // +build test unit
    ```
- Integração 
    ```go
    // +build test integration
    ```

Use os seguintes comandos para executar os testes.

- Ambiente Containerizado.

    ```sh
    make test-all-container
    ```

- Ambiente Local.

    ```sh
    make test-all
    ```

Estes comandos são atalho para a execução do script [`test.sh`](scripts/test.sh) com parâmetro *all*, que resulta em:

```sh
go test -parallel 10 -timeout 1m30s ./... -tags="test" -cover -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -n 1
```

### Carga

Testes de carga totalmente automatizados com containers para os testes e aplicações(targets). 
Também é possível executar fluxos alternativos tendo componetes rodando localmente e outros containerizados.

Basicamente os testes apontam para o target na porta 80 onde executa o nginx, que encaminha as requests 
para a porta 3000, onde roda a API Processador.

```
Jmeter/ApacheBench
        \----- HTTP|POST|:80 -----> ngix(proxy)
                                                \ ----- HTTP|POST|:3000 -----> API Processador
```

O nginx é iniciado pelo Docker Compose com a configuração
```yaml
    network_mode: host
```
que permite testar API Processador executando containerizada ou localmente, talvez em modo debug por uma IDE.

No cenário acima podemos ter os seguintes ambientes de execução.

Serviço | Ambientes
--- | --- 
*Jmeter/ApacheBench* | Containerizado/Local
*ngix(proxy)* | Containerizado
*API Processador* | Containerizado/Local

Como exemplo de execução de testes Jmeter rodando containerizado, e API Processador executando localmente podemos fazer o seguinte.

```sh
make infra-start && make infra-monitoring-start # Inicia (Redis ... etc, Grafana, Prometheus)
make build                                      # Compila em './payment-processor-api.bin'
PROCESSOR_LOGGING_FILE=./processor.log \
./payment-processor-api.bin \
-config-location=./configs/config-dev.json &    # Executa com saída de logging './processor.log'

### Aperte 'Enter' para liberar o console ###

docker-compose up --build test-load-jmeter      # Executa container de testes com binário do Jmeter
```

#### Jmeter

Execução de testes de carga com Jmeter. Definido com 50 threads, mantém uma taxa de aproximadamente 500 req/sec.

Executa vários cenários de processamento de transações, sendo algumas Autorizadas e outras Negadas por diversos motivos. Executa requisições contra *http://localhost:80/api/transactions*.

A definição do plano de testes pode ser encontrada em  [`jmeter-load-test-plan.jmx`](test/jmeter-load-test-plan.jmx). 

![Screenshot do template de Requisição Jmeter](docs/jmeter-request.png "Template de Requisição")

A execução do teste é realizda com o script [`load-test-jmeter.sh`](test/load-test-jmeter.sh).

As configurações do plano de testes são carregadas dinamicamente por meio de Variáveis de Ambiente.

- Configuração com Variáveis de Ambiente

    ```sh
    export JMETER_API_HOST=localhost    # Host do target 
    export JMETER_API_PORT=80           # Porta do target 
    export JMETER_API_CONTEXT=/api      # Contexto do target 
    export JMETER_TEST_CASES_PATH=./jmeter-requests-config.csv # Arquivo com definição de cenários
    export JMETER_EXIT_ON_ERROR=true    # Força saída do programa com código 1 (Indicativo de erro)
    ```

- Cenários de teste

    São definidos no arquivo [`jmeter-requests-config.csv`](test/jmeter-requests-config.csv), e descrevem parâmetros e respostas esperadas para cada requisição.
    Segue exemplo da definição de um cenário.

    ```csv
    testCaseName;requestBody;requestHeaderAcquirerID;responseStatusCode;responseBody
    Stone - Valid - dados básicos;{"token":"xpto121a","holder":"João","total":1000,"installments":1};Stone;200;{"message":"Transação autorizada"}
    ```

    testCaseName | requestBody | requestHeaderAcquirerID | responseStatusCode | responseBody
    --- | --- | --- | --- | --- 
    Stone - Valid - dados básicos | {"token":"xpto121a","holder":"João","total":1000,"installments":1} | Stone | 200 | {"message":"Transação autorizada"}
    Stone - Invalid - Parâmetros vazios | {"token":"xpto121a"} | Stone | 400 | {"title":"Validação da Adquirente ao Processar Transação","detail":"Valor inválido"}
    Cielo - Valid - dados básicos | {"token":"xpto121a","holder":"João","total":500,"installments":1} | Cielo | 200 | {"message":"Transação autorizada"}
    Rede - Invalid - Adquirente não cadastrada | {"token":"xpto121a","holder":"João","total":500,"installments":1} | Rede | 404 | {"title":"Falha ao enviar requisição para Adquirente","detail":"Adquirente \"Rede\" inexistente"}

- Validação dos Cenários

    É usado programação Groovy para aplicar assertions dos cenários. 
    Caso alguma condição falhe é printado na tela e o processo encerrado com status 1, que indica falha na validação de algum cenário.
    Segue script Groovy extraído da [configuração do Jmeter](test/jmeter-load-test-plan.jmx).

    ```groovy
    def responseStatusCode = vars.get("responseStatusCode"); // Obtém assertion vinda do CSV
    def responseBody = vars.get("responseBody"); // Obtém assertion vinda do CSV
    def failureMessage = "";

    if(!responseStatusCode.equals(prev.getResponseCode())){
        failureMessage += "Expected <status code> ["+ responseStatusCode +"] but we got instead [" + prev.getResponseCode() + "]\n" ;
    }

    if (!prev.getResponseDataAsString().startsWith(responseBody)) {
        failureMessage += "Expected <response body> ["+ responseBody + "]. Response body is:[" + prev.getResponseDataAsString() + "]\n";
    }

    if (failureMessage.trim()) {
        // code ...
        AssertionResult.setFailureMessage(failureMessage);
        AssertionResult.setFailure(true);
        println "Erro ao executar teste\n"
        println failureMessage;

        if (vars.get("exit-on-error").toBoolean()) { // Encerra com erro caso exit-on-error=true
            System.exit(1);
        }
    }
    ```

    Ao final da execução containerizada é mostrado no console '*exited with code **0***', que indica que todos os testes passaram com sucesso,
    ou '*exited with code **1***', em casso de falha.

    ```sh
    payment-processor_test-load-jmeter_1 exited with code 0 # Sucesso
    payment-processor_test-load-jmeter_1 exited with code 1 # Falha
    ```

Use os seguintes comandos para executar os testes de carga.

- Ambiente Containerizado(Headless CLI mode).

    ```sh
    make test-load-jmeter-container
    ```

- Ambiente Local(Headless CLI mode).

    ```sh
    make test-load-jmeter-local
    ```

- Ambiente Local em modo gráfico(GUI mode).

    Inicializa o Jmeter em modo gráfico para edição do script .jmx.

    ```sh
    make test-load-jmeter-local gui=s
    ```

#### ApacheBench

Execução de testes de carga com ApacheBench. Definido com 103 threads, mantém uma taxa de aproximadamente 1000 req/sec.

Executa requisições com apenas um tipo de request, de sucesso de autorização contra *http://localhost:80/api/transactions*.

Segue parametrização do teste, que pode ser encontrada em  [`load-test-apachebench.sh`](test/load-test-apachebench.sh).

```sh
ab -n 610000 -c 103 -p $REQUEST_BODY_FILE -T application/json \
    -H 'X-ACQUIRER-ID: Stone' \
    -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjoicGF5bWVudC1wcm9jZXNzb3ItYXBpIn0.uw-8pECPeJbme82nptMI-bsP8f4GvCx9x6b_GzM5wws' \
    -m POST http://localhost:80/api/transactions
```

Use os seguintes comandos para executar os testes de carga.

- Ambiente Containerizado.

    ```sh
    make test-load-ab-container
    ```

- Ambiente Local.

    ```sh
    make test-load-ab-local
    ```

## Instrumentação

Aplicação possui diversos endpoints com informações e exibição de métricas.

### Links Úteis

Página Web na raiz da API com diversos links para endpoints com informações de versão, documentação, build e configurações carregadas.

Execute
```sh
make run
```
e depois acesse http://localhost:80/.

### Grafana

Exibe informações e métricas da aplicação coletadas por meio do Prometheus. Use o [link](http://localhost:3001/d/kKd-m3qiz/requisicoes-http-payment-processor-api?orgId=1&kiosk=tv).

- Métricas expostas pela aplicação.

    - Informações de Build e versão
    - Informações de nome e versões das dependências no go.mod
    - Tempo de Resposta das Requisições
    - Total de requisições
    - Total de falhas graves(panic)

Com as métricas os seguintes gráficos e informações podem ser exibidos.

![Screenshot de Monitoração com Grafana](docs/grafana-monitoracao.png "Monitoração com Grafana")

No print acima é possível observar o comportamento ao executar o teste de carga Jmeter,
e logo depois também iniciar o teste com ApacheBench. 

Vemos que os recursos da máquina local ficaram limitados e algumas indisponibilidades ocorreram(Respostas com Status Code 503).

Outras informaçoes podem ser verificadas.

- Versão da plicação na tag **v0.0.1**.
- Revião(commit) da plicação **c901452**.
- Aplicação compilada com **go1.13.4**.
- Dependência indireta **github.com/prometheus/procfs** está na versão **v0.0.5**.
- Não há errors graves, sendo **0**.
- No gráfico 'Pedidos de Autorização por segundo - status/endpoint/adquirente' a **Stone** recebeu muitas requisições, chegando a fazer quase **900 autorizações de sucesso** por segundo(linha vermelha no gráfico que representa Status 200 para Stone).
- Também é possível ver que **Rede** não é adquirente válida(Status 404), e Stone teve 20 indisponibilidades/segundo(Status 503).
- **524** foi a média de requisições/segundo dos últimos 5 minutos.
- **1%** das requisições começaram a levar mais de 400ms, e atrasando cada vez mais. As médias também ficaram pouco piores.



## Empacotamento

Para empacotar como imagem Docker.

- Base Alpine Linux
    ```sh
    make package # cria imagem challenge/payment-processor-api:latest
    ```

- Base Golang Official. Alternativa criada devido a algumas indisponibilidades percebinas no repositório apk durante desenvolvimento.
    ```sh
    make package-safe # cria imagem challenge/payment-processor-api:latest
    ```

## Comandos Make

Todos comandos para facilitar o desenvolvimento estão no [Makefile](Makefile).

Para listar comandos disponíveis use o seguinte comando.

```sh
make help
```

- Execução
```yaml
run                            : Executa aplicação empacotada em imagem Alpine Linux.
run-safe                       : Executa aplicação empacotada com imagem Golang Official(pesada).
stop                           : Pára aplicação.
```

- Compilação
```yaml
build                          : Compila aplicação. Gera arquivo './payment-processor-api.bin'.
build-static                   : Compila aplicação com lincagem estática.
                                    Ex. 'make build-static path=./'.
wire                           : Gera códigos(wire_gen*.go) do framework de Injeção de Dependências.
```

- Testes
```yaml
test-unit                      : Testes de unidade.
test-integration               : Testes de integração.
test-all                       : Executa testes de unidade e integração.
test-unit-container            : Executa testes de unidade em ambiente containerizado.
test-integration-container     : Executa testes de integração em ambiente containerizado.
test-all-container             : Executa testes de unidade e integração em ambiente containerizado.
test-load-ab-container         : Executa testes de carga com ApacheBench e API containerizada.
test-load-ab-local             : Executa testes de carga com ApacheBench localmente. 
                                    Alternativa a test-load-ab-container.
test-load-ab-stop              : Interrompe containers de testes ApacheBench.
test-load-jmeter-container     : Executa testes de carga com Jmeter e API containerizada.
test-load-jmeter-local         : Executa testes de carga localmente com Jmeter.
                                    Ex. modo gráfico 'make test-load-jmeter-local gui=s'
```

- Deploy Serviços de Dependência
```yaml
infra-start                    : Inicia serviços de dependência containerizados.
infra-stop                     : Interrompe serviços de dependência containerizados.
infra-test-start               : Inicia serviços de dependência de testes containerizados.
infra-test-stop                : Interrompe containers de testes.
infra-monitoring-start         : Inicia serviços de dependência de monitoramento.
infra-monitoring-stop          : Interrompe containers de monitoramento.
```

- Empacotamento Docker
```yaml
package                        : Empacota API na imagem challenge/payment-processor-api:latest usando
                                    Alpine Linux.
package-safe                   : Empacota API na imagem challenge/payment-processor-api:latest usando
                                    Golang Official(pesada).
```

## Melhorias

### Design e Concorrência

- Componente [`Acquirer`](internal/pkg/processor/acquirer.go) que é um Worker está executando muito trabalho de forma procedural.
    É feito o de-para com enriquecimento de dados sensíveis usando Redis e depois é disparada requisição para Adquirente.

    O ponto ruim desta implementação é que não conseguimos controlar individualmente a quantidade de concorrência.
    Não tem como, por exemplo, limitar um número global de 1000 consultas simultâneas ao Redis, e para as Adiquirentes Stone e Cielo, 500 e 200 requisições respectivamente.

    Estes dois passos poderiam ser quebrados, permitindo controle individual. Inclusive `Acquirer` já implementa as interfaces `AcquirerProcessor` e `AcquirerTransactionMapper`.

    `AcquirerProcessor` atual seria renomeado para `AcquirerSender` e teria a seguinte definição.
    
    ```go
    type AcquirerProcessor interface {
        Process(r *AuthorizationRequest)
    }

    // AcquirerSender` envia requisição com dados sensíveis da Transação presentes em r.
    type AcquirerSender interface {
        Send(r *AuthorizationRequest)
    }
    ```
