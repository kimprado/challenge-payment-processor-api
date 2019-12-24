package app

import (
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/challenge/payment-processor/internal/pkg/processor"
	"github.com/challenge/payment-processor/internal/pkg/webserver"
)

// PaymentProcessorApp representa instância da aplicação
type PaymentProcessorApp struct {
	cardsLoader *DefaultCardsLoader
	webServer   *webserver.WebServer
	logger      logging.Logger
}

// NewPaymentProcessorApp cria app
func NewPaymentProcessorApp(cl *DefaultCardsLoader, ws *webserver.WebServer, sw *processor.StoneAcquirerWorkers, cw *processor.CieloAcquirerWorkers, l logging.Logger) (a *PaymentProcessorApp) {
	a = new(PaymentProcessorApp)
	a.cardsLoader = cl
	a.webServer = ws
	a.logger = l
	return
}

// Bootstrap é responsável por iniciar a aplicação
func (a *PaymentProcessorApp) Bootstrap() {
	a.logger.Infof("Iniciando serviços da aplicação...\n")

	a.cardsLoader.load()

	a.webServer.Start()
}
