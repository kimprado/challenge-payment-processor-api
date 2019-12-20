package processor

import "github.com/challenge/payment-processor/internal/pkg/infra/http"

// AcquirerProcessor representa adquirente capaz de processar
// transação com cartão.
type AcquirerProcessor interface {
	Process(r *AuthorizationRequest)
}

// AcquirerTransactionMapper representa comportamento capaz de fazer de-para,
// e transformar dados para formato que adquirente espera.
type AcquirerTransactionMapper interface {
	MapTransaction()
}

// AcquirerID representa identificação de um Acquirer
type AcquirerID string

// Acquirer implementa funcionalidades de de-para e envio
// da transação para Adquirente.
type Acquirer struct {
	url  string
	http http.RequestSender
}

// Process implementa AcquirerProcessor
func (a *Acquirer) Process(r *AuthorizationRequest) {
	a.http.Send()
}

func newAcquirer() (a *Acquirer) {

	return
}

// StoneAcquirerWorkers reprensenta trabalhadores de
// Stone Acquirer
type StoneAcquirerWorkers struct {
	AcquirerWorkers
}

// NewStoneAcquirerWorkers cria instância de StoneAcquirerWorkers.
func NewStoneAcquirerWorkers(a AcquirerActorsResgister) (w *StoneAcquirerWorkers) {
	w = new(StoneAcquirerWorkers)
	w.aid = "Stone"
	w.chr = make(chan *AuthorizationRequest)
	w.acquirers = []AcquirerProcessor{}

	for i := 0; i < 10; i++ {
		w.acquirers = append(w.acquirers, newAcquirer())
	}

	a.Resgister(w.aid, w.chr)

	w.consume()

	return
}

// AcquirerWorkers reprensenta trabalhadores que delegam
// trabalho para Acquirers
type AcquirerWorkers struct {
	aid       AcquirerID
	chr       chan *AuthorizationRequest
	acquirers []AcquirerProcessor
}

func (w *AcquirerWorkers) consume() {

	var processTransactions = func(acquirer AcquirerProcessor) {
		for r := range w.chr {
			acquirer.Process(r)
		}
	}

	for _, acquirer := range w.acquirers {
		go processTransactions(acquirer)
	}

}
