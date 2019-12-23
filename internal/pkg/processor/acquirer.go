package processor

import (
	"errors"

	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/infra/http"
)

// AcquirerProcessor representa adquirente capaz de processar
// transação com cartão.
type AcquirerProcessor interface {
	Process(r *AuthorizationRequest)
}

// AcquirerTransactionMapper representa comportamento capaz de fazer de-para,
// e transformar dados para formato que adquirente espera.
type AcquirerTransactionMapper interface {
	mapTransaction(et *ExternalTransactionDTO) (t *AcquirerTransactionDTO, err error)
}

// AcquirerID representa identificação de um Acquirer
type AcquirerID string

func (a AcquirerID) String() string {
	return (string)(a)
}

// AcquirerParameter encapsula parâmetros para criação de Acquirer
type AcquirerParameter struct {
	cardFinder CardRepositoryFinder
	httpSender http.RequestSender
}

// NewAcquirerParameter cria instância de AcquirerParameter
func NewAcquirerParameter(r CardRepositoryFinder, rs http.RequestSender) (c *AcquirerParameter) {
	c = new(AcquirerParameter)
	c.cardFinder = r
	c.httpSender = rs
	return
}

// Acquirer implementa funcionalidades de de-para e envio
// da transação para Adquirente.
type Acquirer struct {
	url string
	*AcquirerParameter
}

func newAcquirer(url string, p *AcquirerParameter) (a *Acquirer) {
	a = new(Acquirer)
	a.url = url
	a.AcquirerParameter = p
	return
}

// Process implementa AcquirerProcessor
func (a *Acquirer) Process(r *AuthorizationRequest) {
	var err error
	var t *AcquirerTransactionDTO
	t, err = a.mapTransaction(r.Transaction)
	if err != nil {
		r.ResponseChannel <- &AuthorizationResponse{Err: err}
		return
	}

	var response AuthorizationMessage
	err = a.httpSender.Send(a.url, t, &response)

	var httpError *http.Error
	if errors.Is(err, &http.StatusBadRequestError{}) && errors.As(err, &httpError) {
		r.ResponseChannel <- &AuthorizationResponse{Err: newAcquirerValidationError(httpError.Message, httpError.URL)}
		return
	}

	if errors.Is(err, &http.ServerError{}) {
		r.ResponseChannel <- &AuthorizationResponse{Err: NewAcquirerProcessingError()}
		return
	}

	if err != nil {
		r.ResponseChannel <- &AuthorizationResponse{Err: newAcquirerConnectivityError("Erro de conexão com Adquirente. Tente novamente mais tarde.", err)}
		return
	}

	r.ResponseChannel <- &AuthorizationResponse{Authorized: &response}
}

// mapTransaction implementa AcquirerTransactionMapper
func (a *Acquirer) mapTransaction(et *ExternalTransactionDTO) (t *AcquirerTransactionDTO, err error) {
	var c *Card
	c, err = a.cardFinder.Find(et.Token)

	if err != nil {
		return
	}

	if c == nil {
		err = NewCardNotFoundError()
		return
	}

	t = new(AcquirerTransactionDTO)
	t.TransactionDTO = et.TransactionDTO
	t.CardHiddenInfoDTO = &CardHiddenInfoDTO{Number: c.Number, CVV: c.CVV}
	return
}

// StoneAcquirerWorkers reprensenta trabalhadores de
// Stone Acquirer
type StoneAcquirerWorkers struct {
	*AcquirerWorkers
}

// NewStoneAcquirerWorkers cria instância de StoneAcquirerWorkers.
func NewStoneAcquirerWorkers(a AcquirerActorsResgister, p *AcquirerParameter, c config.Configuration) (w *StoneAcquirerWorkers) {
	w = new(StoneAcquirerWorkers)
	w.AcquirerWorkers = newAcquirerWorkers("Stone", a)
	for i := 0; i < 10; i++ {
		w.add(newAcquirer(c.StoneAcquirer.URL, p))
	}
	return
}

// CieloAcquirerWorkers reprensenta trabalhadores de
// Cielo Acquirer
type CieloAcquirerWorkers struct {
	*AcquirerWorkers
}

// NewCieloAcquirerWorkers cria instância de CieloAcquirerWorkers.
func NewCieloAcquirerWorkers(a AcquirerActorsResgister, p *AcquirerParameter, c config.Configuration) (w *CieloAcquirerWorkers) {
	w = new(CieloAcquirerWorkers)
	w.AcquirerWorkers = newAcquirerWorkers("Cielo", a)
	for i := 0; i < 10; i++ {
		w.add(newAcquirer(c.CieloAcquirer.URL, p))
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
