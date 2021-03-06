# Self-Documenting Makefiles https://swcarpentry.github.io/make-novice/08-self-doc/index.html
## help				: Exibe comandos make disponíveis.
.PHONY : help
help : Makefile
	@sed -n 's/^##//p' $<

# TODO: verificar alternativa ao Alpine Linux. Indisponibilidades no Repositório apk
## run				: Executa aplicação empacotada em imagem Alpine Linux.
run: 
	@./scripts/deploy.sh start
	@echo ""
	@echo "Acesse nginx:"
	@echo "http://localhost:80/"
	@echo "Acesse API:"
	@echo "http://`docker-compose port api 3000`/"
	@echo "Acesse swagger:"
	@echo "http://localhost:80/docs"
	@echo "Acesse grafana:"
	@echo "http://localhost:3001/d/kKd-m3qiz/requisicoes-http-payment-processor-api?orgId=1&kiosk=tv"
	@echo "Acesse prometheus:"
	@echo "http://localhost:9090"
	@echo ""

# Alternativa criada devido a algumas indisponibilidades percebidas no 
# repositório apk durante desenvolvimento.
## run-safe			: Executa aplicação empacotada com imagem Golang Official(pesada).
run-safe: 
	@./scripts/deploy.sh start-safe
	@echo ""
	@echo "Acesse nginx:"
	@echo "http://localhost:80/"
	@echo "Acesse API:"
	@echo "http://`docker-compose port api-safe 3000`/"
	@echo "Acesse swagger:"
	@echo "http://localhost:80/docs"
	@echo "Acesse grafana:"
	@echo "http://localhost:3001/d/kKd-m3qiz/requisicoes-http-payment-processor-api?orgId=1&kiosk=tv"
	@echo "Acesse prometheus:"
	@echo "http://localhost:9090"
	@echo ""

## stop				: Pára aplicação.
stop:
	@./scripts/deploy.sh stop

## build				: Compila aplicação. Gera arquivo './payment-processor-api.bin'.
build:
	@./scripts/compile.sh build

## build-static			: Compila aplicação com lincagem estática. Ex: 'make build-static path=./'.
build-static:
	@./scripts/compile.sh build-static $(path)

## wire				: Gera/Atualiza códigos(wire_gen*.go) do framework de Injeção de Dependências.
wire:
	@./scripts/compile.sh wire
	@./scripts/compile.sh wire-testes

## test-unit			: Testes de unidade.
test-unit:
	@./scripts/test.sh unit

## test-integration		: Testes de integração.
test-integration:
	@./scripts/test.sh integration

## test-all			: Executa testes de unidade e integração.
test-all:
	@./scripts/test.sh all

## test-unit-container		: Executa testes de unidade em ambiente containerizado.
test-unit-container:
	@docker-compose up --build test-unit

## test-integration-container	: Executa testes de integração em ambiente containerizado.
test-integration-container:
	@docker-compose up --build test-integration

## test-all-container		: Executa testes de unidade e integração em ambiente containerizado.
test-all-container:
	@docker-compose up --build test-all

## test-load-ab-container		: Executa testes de carga com ApacheBench e API containerizada.
test-load-ab-container: infra-monitoring-start
	@docker-compose up -d --build  acquirers api redisdb test-load-ab
	@docker-compose logs --tail="100" -f test-load-ab &

## test-load-ab-local		: Executa testes de carga com ApacheBench localmente. Alternativa a test-load-ab-container.
test-load-ab-local:
	@third_party/wait-for-it.sh localhost:3000 -s -t 30 -- ./test/load-test-apachebench.sh

## test-load-ab-stop		: Interrompe containers de testes ApacheBench.
test-load-ab-stop:
	@docker-compose rm -fsv test-load-ab api-load

## test-load-jmeter-container	: Executa testes de carga com Jmeter e API containerizada.
test-load-jmeter-container: infra-monitoring-start
	@docker-compose up -d --build  acquirers api redisdb test-load-jmeter
	@docker-compose logs --tail="100" -f test-load-jmeter &

## test-load-jmeter-local		: Executa testes de carga localmente com Jmeter. Ex: modo gráfico 'make test-load-jmeter-local gui=s'
test-load-jmeter-local:
	@./test/load-test-jmeter.sh $(gui)

## infra-start			: Inicia serviços de dependência containerizados.
infra-start:
	@docker-compose up -d --build redisdb nginx acquirers

## infra-stop			: Interrompe serviços de dependência containerizados.
infra-stop:
	@docker-compose rm -fsv redisdb nginx swagger acquirers

## infra-test-start		: Inicia serviços de dependência de testes containerizados.
infra-test-start:
	@docker-compose up -d --build redis-test acquirers

## infra-test-stop		: Interrompe containers de testes.
infra-test-stop:
	@docker-compose rm -fsv test-unit test-integration test-all test-load-ab test-load-jmeter redis-test acquirers

## infra-monitoring-start		: Inicia serviços de dependência de monitoramento.
infra-monitoring-start:
	@docker-compose up -d --build prometheus grafana

## infra-monitoring-stop		: Interrompe containers de monitoramento.
infra-monitoring-stop:
	@docker-compose rm -fsv prometheus grafana

## package			: Empacota API na imagem challenge/payment-processor-api:latest - Alpine Linux.
package: 
	@./scripts/package.sh package

## package-safe			: Empacota API na imagem challenge/payment-processor-api:latest - Golang Official(pesada).
package-safe: 
	@./scripts/package.sh package-safe
