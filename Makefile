# Self-Documenting Makefiles https://swcarpentry.github.io/make-novice/08-self-doc/index.html
## help				: Exibe comandos make disponíveis.
.PHONY : help
help : Makefile
	@sed -n 's/^##//p' $<

# TODO: verificar alternativa ao Alpine Linux. Indisponibilidades no Repositório apk
## run				: Executa aplicação empacotada em imagem Alpine Linux.
run: 
	@./scripts/deploy.sh start
	@echo "Acesse nginx:"
	@echo "http://localhost:80/"
	@echo "Acesse API:"
	@echo "http://`docker-compose port api 3000`/"
	@echo "Acesse swagger:"
	@echo "http://localhost:80/docs"

# Alternativa criada devido a algumas indisponibilidades percebidas no 
# repositório apk durante desenvolvimento.
## run-safe			: Executa aplicação empacotada com imagem Golang Official(pesada).
run-safe: 
	@./scripts/deploy.sh start-safe
	@echo "Acesse nginx:"
	@echo "http://localhost:80/"
	@echo "Acesse API:"
	@echo "http://`docker-compose port api-safe 3000`/"
	@echo "Acesse swagger:"
	@echo "http://localhost:80/docs"

## stop				: Pára aplicação.
stop:
	@./scripts/deploy.sh stop

## build				: Compila aplicação. Gera arquivo './payment-processor-api.bin'.
build:
	@./scripts/compile.sh build

## build-static			: Compila aplicação com lincagem estática. Ex: 'make build-static path=./'.
build-static:
	@./scripts/compile.sh build-static $(path)

## wire				: Gera/Atualiza códigos(wire_gen.go) do framework de Injeção de Dependências.
wire:
	@./scripts/compile.sh wire
	@./scripts/compile.sh wire-testes

## generate			: Atualiza códigos(wire_gen.go) do framework de Injeção de Dependências.
generate:
	@./scripts/compile.sh generate
	@./scripts/compile.sh generate-testes

## test-unit			: Testes de unidade.
test-unit:
	@./scripts/test.sh unit

## test-integration		: Testes de integração.
test-integration:
	@./scripts/test.sh integration

## test-all			: Executa testes de unidade e integração.
test-all:
	@./scripts/test.sh all

test-envvars:
	@./scripts/test.sh envvars

## test-unit-container		: Executa testes de unidade em ambiente containerizado.
test-unit-container:
	@docker-compose up --build test-unit

## test-integration-container	: Executa testes de integração em ambiente containerizado.
test-integration-container:
	@docker-compose up --build test-integration

## test-all-container		: Executa testes de unidade e integração em ambiente containerizado.
test-all-container:
	@docker-compose up --build test-all

test-envvars-container:
	@docker-compose up --build test-envvars

## test-load-ab			: Executa testes de carga com ApacheBench.
test-load-ab:
	@docker-compose up -d --build redisdb test-load-ab
	@docker-compose logs --tail="10" -f test-load-ab &

## test-load-ab-container		: Executa testes de carga com ApacheBench e API containerizada.
test-load-ab-container:
	@docker-compose up -d --build  api-load redisdb test-load-ab
	@docker-compose logs --tail="10" -f test-load-ab &

## infra-start			: Inicia serviços de dependência containerizados.
infra-start:
	@docker-compose up -d --build redisdb nginx

## infra-stop			: Interrompe serviços de dependência containerizados.
infra-stop:
	@docker-compose rm -fsv redisdb nginx swagger

## infra-test-start		: Inicia serviços de dependência de testes containerizados.
infra-test-start:
	@docker-compose up -d --build redis-test

## infra-test-stop		: Interrompe serviços de dependência de testes containerizados.
infra-test-stop:
	@docker-compose rm -fsv redis-test

## package			: Empacota API na imagem challenge/payment-processor-api:latest - Alpine Linux.
package: 
	@./scripts/package.sh package

## package-safe			: Empacota API na imagem challenge/payment-processor-api:latest - Golang Official(pesada).
package-safe: 
	@./scripts/package.sh package-safe