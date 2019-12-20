// +build test unit

package processor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateStoneWorker(t *testing.T) {
	t.Parallel()

	var w *StoneAcquirerWorkers
	var r *AcquirerActorsMock
	var p *AcquirerParameter
	var s *HTTPRequestSenderMock
	var ar *AuthorizationRequest
	var repo *CardRepositoryFinderMock

	r = &AcquirerActorsMock{}

	s = newHTTPRequestSenderMock()
	repo = newCardRepositoryFinderMock()
	p = &AcquirerParameter{
		url:        "htttp://localhost/acquirer/stone",
		http:       s,
		cardFinder: repo,
	}

	ar = &AuthorizationRequest{Transaction: &ExternalTransactionDTO{TransactionDTO: &TransactionDTO{CardOpenInfoDTO: &CardOpenInfoDTO{Holder: "João"}}}}

	w = NewStoneAcquirerWorkers(r, p)
	assert.NotNil(t, w)

	assert.True(t, r.Resgistered)

	w.chr <- ar

	select {
	case <-repo.Found:
		t.Log("Consulta de cartão no BD realizada")
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Consulta de cartão no BD não foi realizada")
	}

	select {
	case <-s.Sent:
		t.Log("Requisição http enviada")
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Requisição http não foi enviada")
	}

}

func TestCieloStoneWorker(t *testing.T) {
	t.Parallel()

	var w *CieloAcquirerWorkers
	var r *AcquirerActorsMock
	var p *AcquirerParameter
	var s *HTTPRequestSenderMock
	var ar *AuthorizationRequest
	var repo *CardRepositoryFinderMock

	r = &AcquirerActorsMock{}

	s = newHTTPRequestSenderMock()
	repo = newCardRepositoryFinderMock()
	p = &AcquirerParameter{
		url:        "htttp://localhost/acquirer/cielo",
		http:       s,
		cardFinder: repo,
	}

	ar = &AuthorizationRequest{Transaction: &ExternalTransactionDTO{TransactionDTO: &TransactionDTO{CardOpenInfoDTO: &CardOpenInfoDTO{Holder: "João"}}}}

	w = NewCieloAcquirerWorkers(r, p)
	assert.NotNil(t, w)

	assert.True(t, r.Resgistered)

	w.chr <- ar

	select {
	case <-repo.Found:
		t.Log("Consulta de cartão no BD realizada")
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Consulta de cartão no BD não foi realizada")
	}

	select {
	case <-s.Sent:
		t.Log("Requisição http enviada")
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Requisição http não foi enviada")
	}

}

type HTTPRequestSenderMock struct {
	Sent chan bool
}

func newHTTPRequestSenderMock() (s *HTTPRequestSenderMock) {
	s = new(HTTPRequestSenderMock)
	s.Sent = make(chan bool)
	return
}

func (s *HTTPRequestSenderMock) Send(url string, body interface{}, response interface{}) (err error) {

	s.Sent <- true

	return
}

type AcquirerActorsMock struct {
	Resgistered bool
	chr         chan *AuthorizationRequest
}

func (a *AcquirerActorsMock) Resgister(aid AcquirerID, chr chan *AuthorizationRequest) (err error) {
	a.Resgistered = true
	a.chr = chr
	return
}

type CardRepositoryFinderMock struct {
	Found chan bool
}

func newCardRepositoryFinderMock() (r *CardRepositoryFinderMock) {
	r = new(CardRepositoryFinderMock)
	r.Found = make(chan bool)
	return
}

func (r *CardRepositoryFinderMock) Find(token string) (c *Card, err error) {

	r.Found <- true

	return
}
