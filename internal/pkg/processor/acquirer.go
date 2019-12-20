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
	*AcquirerWorkers
}

// NewStoneAcquirerWorkers cria instância de StoneAcquirerWorkers.
func NewStoneAcquirerWorkers(a AcquirerActorsResgister) (w *StoneAcquirerWorkers) {
	w = new(StoneAcquirerWorkers)
	w.AcquirerWorkers = newAcquirerWorkers("Stone", a)
	for i := 0; i < 10; i++ {
		w.add(newAcquirer())
	}
	return
}

// AcquirerWorkers reprensenta trabalhadores que delegam
// trabalho para Acquirers
type AcquirerWorkers struct {
	aid AcquirerID
	chr chan *AuthorizationRequest
}

// newAcquirerWorkers cria instância de AcquirerWorkers.
func newAcquirerWorkers(aid AcquirerID, a AcquirerActorsResgister) (w *AcquirerWorkers) {
	w = new(AcquirerWorkers)
	w.aid = aid
	w.chr = make(chan *AuthorizationRequest)

	a.Resgister(w.aid, w.chr)

	return
}

func (w *AcquirerWorkers) add(acquirer AcquirerProcessor) {
	var processTransactions = func(acquirer AcquirerProcessor) {
		for r := range w.chr {
			acquirer.Process(r)
		}
	}
	go processTransactions(acquirer)
}
