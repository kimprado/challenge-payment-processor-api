# Documentação Payment Processor API

Descrição da solução para o desafio do Processador de Pagamentos em Golang([Back-end](https://github.com/pagarme/vagas/tree/master/desafios/software-engineer-golang))

- [O Problema](#O-Problema)
- [Back-end](#Back-end)
    - [Dependências](#Dependências)
        - [Boilerplate Code](#Boilerplate-Code)
- [Documentação API](#Documentação-API)
- [Instalação e Execução](#Instalação-e-Execução)
- [Ambiente Desenvolvimento](#Ambiente-Desenvolvimento)
    - [Primeira Execução](#Primeira-Execução)
        - [Instalação das Dependências](#Instalação-das-Dependências)
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
- [Comandos Make](#comandos-make)
- [Melhorias](#Melhorias)

## O Problema

Foi implementada solução que permite fazer processamento de solicitação de autorização de pagamento, em colaboração com Adquirentes.

Processador foi implementado como API Rest que recebe requisições em um endpoint, enriquece requisição com dados sensíveis de cartão armazenados em banco de dados Redis, e depois encaminha requisição para serviço Rest da Adquirente.

A seguir temos exemplos de utilização da API. Para uma documentação mais completa verifique o tópico [swagger](#Documentação-API).

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

- Na sequência algum Worker `Acquirer` recebe a transação(Fan-out), a enriquece com informações sensíveis armazenadas no Redis, e depois envia esses dados para o Serviço da Adquirente por meio de outra chamada HTTTP.

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
│       │   │   ├── apiIT_test.go   # Testes de integração
│       │   │   ├── api_test.go     # Testes de unidade
│       │   │   ├── errors.go       # Mapeia erros de negócio para erros HTTP Status Code
│       │   ├── card.go             # Entidade com dados sensíveis persitida no Redis
│       │   ├── dto.go              # Implementa DTO usado para recebimento de dados
│       │   ├── errors.go           # Erros de negócio
│       │   ├── processor.go        # Ponto de entrada do negócio. Atende solicitações da API
│       │   ├── processor_test.go   # Testes de unidade com Mock Objects
│       │   ├── repository.go       # Meio de Persitência com Redis
│       │   ├── repositoryIT_test.go
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

Documentação disponibilizada no [arquivo](api/swagger.yml) implementada com swagger. Para acessar a documentação interativa execute o ambiente como descrito a seguir([Instalação e Execução](#Instalação-e-Execução)), e depois siga as instruções em [Infra Documentação](#Infra-Documentação).


## Instalação e Execução

Para fazer deploy e execução do projeto rode os seguintes comandos:

```sh
./configure
make run
```

Ao final na execução o comando printa no console as URLs para acesso aos serviço.

- http://localhost:80/              (nginx)      - Página web com links úteis
- http://0.0.0.0:3000/              (API)        - URL da API
- http://localhost:80/docs          (swagger)    - DOcumentação interativa
- http://localhost:3001/d/kKd-m3qiz (grafana)    - URL do Grafana
- http://localhost:9090             (prometheus) - URL do Prometheus

Ex:

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

### Primeira Execução

#### Instalação das Dependências

Execute script [`configure`](configure) presente na raiz do repositório para fazer download e instalação das dependências. O script também cria, caso necessário, a configuração da IDE e de execução da aplicação.

```sh
./configure
```

Após executar feche e abra outro terminal. Alguns comandos talvez ainda precisem ser executados com sudo, ou reinicie o computador.

As seguintes ferramentas serão provisionadas, caso necessário.

- **Docker** - Ferramenta usada para containerização.
- **Docker Compose** - Ferramenta usada para orquestração em ambiende de dev.
- **Go** - Linguagem de programação.
- **Jmeter** - Ferramenta usada para testes de carga. Testes podem ser executados a partir de containers, e também localmente.

Os seguintes arquivos são criados, caso necessário.

- .vscode/settings.json - Arquivo da IDE VSCode
- .vscode/launch.json - Arquivo da IDE VSCode
- configs/config.env - Configurações injetadas como Variáveis de Ambiente no deploy com [Docker Compose](#Execução).
- configs/config-dev.env - Configurações padão de desenvolvimento.
- configs/config-dev.json - Configuração opcional da aplicação em tempo de desenvolvimento.

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

### Feedback rápido

Fizemos separação dos testes em dois grupos principais, *[Unitários](#Unitários)* e *[Integração](#Integração)*. Os grupos de testes são separados em arquivos de testes diferentes. Usei o conceito de [Build Constraints ou Build Tag](http://golang.org/pkg/go/build/#hdr-Build_Constraints) para selecionar quais testes queremos executar.

Para especificar um grupo de teste executamos o comando *go test* com o parâmetro *-tags*.

```sh
go test ./internal/pkg/commom/config -tags="unit"
```

Neste exemplo o pacote *config* possui os seguintes arquivos de teste:

- [config_test.go](internal/pkg/commom/config/config_test.go)

    ```go
    // +build test unit

    package config
    // ...
    ```

- [configEnvVarsIT_test.go](internal/pkg/commom/config/configEnvVarsIT_test.go)

    ```go
    // +build testenvvars

    package config
    // ...
    ```

- [configIT_test.go](internal/pkg/commom/config/configIT_test.go)

    ```go
    // +build test integration

    package config
    // ...
    ```

Apenas os testes do arquivo [config_test.go](internal/pkg/commom/config/config_test.go), com a build tag *"// +build test **unit**"*, serão executados pois no comando informamos *-tags="**unit**"*.

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

Estes comandos são atalho para a execução do script [test.sh](scripts/test.sh) com parâmetro *unit*, que resulta em:

```sh
go test ./... -tags="unit" -cover -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -n 1
```

Para configurar um arquivo como Unit Test:

- Sufixo - *_test.go
- Build Tag - *unit*
    - Ex: // +build test unit
    - Ex: arquivo [config_test.go](internal/pkg/commom/config/config_test.go)
    - Ex: arquivo [processor_test.go](internal/pkg/processor/processor_test.go)

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

Estes comandos são atalho para a execução do script [test.sh](scripts/test.sh) com parâmetro *integration*, que resulta em:

```sh
go test -parallel 10 -timeout 1m30s ./... -tags="integration" -cover -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -n 1
```

Para configurar um arquivo como Integration Test:

- Sufixo - *IT_test.go
- Build Tag - *integration*
    - Ex: // +build test integration
    - Ex: arquivo [apiIT_test.go](internal/pkg/processor/api/apiIT_test.go)

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

Estes comandos são atalho para a execução do script [test.sh](scripts/test.sh) com parâmetro *all*, que resulta em:

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


#### ApacheBench

- Execução de testes de carga.

Executa requisições contra *http://localhost:80/api/transactions*.

```sh
make test-load-ab-container
```

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
