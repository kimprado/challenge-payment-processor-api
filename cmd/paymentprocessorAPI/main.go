package main

import (
	"fmt"
	"log"
	"os"

	app "github.com/challenge/payment-processor/internal/app"
	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
)

func main() {
	fmt.Println("iniciando serviço")
	defer fmt.Println("serviço encerrado")

	var err error
	var config config.Configuration
	var appender logging.FileAppender
	var app *app.PaymentProcessorApp

	config, err = initializeConfig("")
	if err != nil {
		fmt.Printf("Erro ao carregar configurações %v\n", err)
		return
	}

	log.Println("Configurações carregadas!")

	if config.Logging.File != "" {
		log.Printf("Arquivo de logging %q\n", config.Logging.File)
		appender, _ = initializeAppender("")
		appender.Configure()
	} else {
		log.SetOutput(os.Stdout)
	}

	app, err = initializeApp("")
	if err != nil {
		fmt.Printf("Erro ao iniciar aplicação %v\n", err)
		return
	}

	app.Bootstrap()
}
