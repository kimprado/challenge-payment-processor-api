// +build test unit

package processor

import (
	"testing"
	"time"

	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/infra/http"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func TestProcessAuthorizationRequest(t *testing.T) {
	t.Parallel()

	var a *Acquirer
	var p *AcquirerParameter
	var s *HTTPRequestSenderMock
	var repo *CardRepositoryFinderMock
	var ar *AuthorizationRequest

	url := "htttp://localhost/acquirer/stone"
	s = newHTTPRequestSenderMock()
	repo = newCardRepositoryFinderMock()
	p = &AcquirerParameter{
		httpSender: s,
		cardFinder: repo,
	}

	ar = &AuthorizationRequest{ResponseChannel: make(chan *AuthorizationResponse, 1), Transaction: &ExternalTransactionDTO{TransactionDTO: &TransactionDTO{CardOpenInfoDTO: &CardOpenInfoDTO{Holder: "João"}}}}

	a = newAcquirer(url, p)
	assert.NotNil(t, a)

	a.Process(ar)

	select {
	case <-repo.Found:
		t.Log("Consulta de cartão no BD realizada")
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Consulta de cartão no BD não foi realizada")
	}

	select {
	case <-s.Sent:
		t.Log("Requisição http enviada")
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Requisição http não foi enviada")
	}

	select {
	case <-ar.ResponseChannel:
		t.Log("Resposta de processamento enviada")
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Resposta de processamento não foi enviada")
	}

}

func TestProcessAuthorizationRequestCases(t *testing.T) {
	t.Parallel()

	var a *Acquirer
	var p *AcquirerParameter

	var cvva = "cvva"
	var cvvblank = ""

	testCases := []struct {
		//label indica título do Test Case
		label string
		url   string
		ar    *AuthorizationRequest
		repo  *CardRepositoryFinderCaseMock
		s     *HTTPRequestSenderCaseMock
		cvv   string
		err   error
	}{
		{
			"Requisição Válida",
			"htttp://localhost/acquirer/stone",
			newAuthorizationRequest("xpto121a", "João", 1000, 1),
			newCardRepositoryFinderCaseMock(func(token string, chp chan string) (c *Card, err error) {
				chp <- token
				c = &Card{CVV: cvva}
				return
			}),
			newHTTPRequestSenderCaseMock(func(p *requestParam, chp chan *requestParam) (err error) {
				chp <- p
				return
			}),
			cvva,
			nil,
		},
		{
			"Erro Infra Repositório",
			"htttp://localhost/acquirer/stone",
			newAuthorizationRequest("xpto121a", "João", 1000, 1),
			newCardRepositoryFinderCaseMock(func(token string, chp chan string) (c *Card, err error) {
				chp <- token
				err = redis.ErrPoolExhausted
				return
			}),
			nil,
			cvvblank,
			redis.ErrPoolExhausted,
		},
		{
			"Erro Cartão não encontrado",
			"htttp://localhost/acquirer/stone",
			newAuthorizationRequest("xpto121a", "João", 1000, 1),
			newCardRepositoryFinderCaseMock(func(token string, chp chan string) (c *Card, err error) {
				chp <- token
				c = nil
				return
			}),
			nil,
			cvvblank,
			&CardNotFoundError{},
		},
		{
			"Transação negada pelo Adquirente",
			"htttp://localhost/acquirer/stone",
			newAuthorizationRequest("xpto121a", "João", 1000, 1),
			newCardRepositoryFinderCaseMock(func(token string, chp chan string) (c *Card, err error) {
				chp <- token
				c = &Card{CVV: cvva}
				return
			}),
			newHTTPRequestSenderCaseMock(func(p *requestParam, chp chan *requestParam) (err error) {
				chp <- p
				err = &http.StatusBadRequestError{
					Message: "HTTP Bad Request",
					Err:     &http.Error{URL: "htttp://localhost/acquirer/stone", Message: "Valor muito alto", Code: 400},
				}
				return
			}),
			cvva,
			&AcquirerValidationError{},
		},
		{
			"Transação com erro interno no Adquirente",
			"htttp://localhost/acquirer/stone",
			newAuthorizationRequest("xpto121a", "João", 1000, 1),
			newCardRepositoryFinderCaseMock(func(token string, chp chan string) (c *Card, err error) {
				chp <- token
				c = &Card{CVV: cvva}
				return
			}),
			newHTTPRequestSenderCaseMock(func(p *requestParam, chp chan *requestParam) (err error) {
				chp <- p
				err = &http.ServerError{
					Message: "Server Error",
					Err:     &http.Error{URL: "htttp://localhost/acquirer/stone", Message: "Conexão DB", Code: 503},
				}
				return
			}),
			cvva,
			&AcquirerProcessingError{},
		},
		{
			"Requisição com erro HTTP não mapeado",
			"htttp://localhost/acquirer/stone",
			newAuthorizationRequest("xpto121a", "João", 1000, 1),
			newCardRepositoryFinderCaseMock(func(token string, chp chan string) (c *Card, err error) {
				chp <- token
				c = &Card{CVV: cvva}
				return
			}),
			newHTTPRequestSenderCaseMock(func(p *requestParam, chp chan *requestParam) (err error) {
				chp <- p
				err = &http.Error{
					Message: "Redirect",
					Err:     &http.Error{URL: "htttp://localhost/acquirer/stone", Message: "Novo servidor", Code: 300},
				}
				return
			}),
			cvva,
			&AcquirerConnectivityError{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {

			p = &AcquirerParameter{
				httpSender: tc.s,
				cardFinder: tc.repo,
			}

			a = newAcquirer(tc.url, p)
			assert.NotNil(t, a)

			a.Process(tc.ar)

			select {
			case token := <-tc.repo.chParam:
				t.Log("Consulta de cartão no BD realizada")
				assert.Equal(t, tc.ar.Transaction.Token, token)
			case <-time.After(1 * time.Second):
				assert.Fail(t, "Consulta de cartão no BD não foi realizada")
			}

			if tc.s != nil {
				select {
				case httpReq := <-tc.s.chParam:
					t.Log("Requisição http enviada")
					assert.Equal(t, tc.url, httpReq.url)

					//Validar informação CVV sensível
					assert.Equal(t, tc.cvv, httpReq.body.CVV)

					assert.Equal(t, tc.ar.Transaction.Holder, httpReq.body.Holder)
					assert.Equal(t, tc.ar.Transaction.Total, httpReq.body.Total)
					assert.Equal(t, tc.ar.Transaction.Installments, httpReq.body.Installments)
				case <-time.After(1 * time.Second):
					assert.Fail(t, "Requisição http não foi enviada")
				}
			}

			select {
			case resp := <-tc.ar.ResponseChannel:
				t.Log("Resposta de processamento enviada")

				if tc.err != nil {
					assert.IsType(t, tc.err, resp.Err)
					return
				}

				assert.NotNil(t, resp.Authorized)
			case <-time.After(1 * time.Second):
				assert.Fail(t, "Resposta de processamento não foi enviada")
			}

		})
	}

}

func TestCreateStoneWorker(t *testing.T) {
	t.Parallel()

	var w *StoneAcquirerWorkers
	var r *AcquirerActorsMock
	var p *AcquirerParameter
	var c config.Configuration
	var s *HTTPRequestSenderMock
	var ar *AuthorizationRequest
	var repo *CardRepositoryFinderMock

	r = &AcquirerActorsMock{}

	c.StoneAcquirer.URL = "htttp://localhost/acquirer/stone"
	s = newHTTPRequestSenderMock()
	repo = newCardRepositoryFinderMock()
	p = &AcquirerParameter{
		httpSender: s,
		cardFinder: repo,
	}

	ar = &AuthorizationRequest{ResponseChannel: make(chan *AuthorizationResponse), Transaction: &ExternalTransactionDTO{TransactionDTO: &TransactionDTO{CardOpenInfoDTO: &CardOpenInfoDTO{Holder: "João"}}}}

	w = NewStoneAcquirerWorkers(r, p, c)
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

	select {
	case <-ar.ResponseChannel:
		t.Log("Resposta de processamento enviada")
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Resposta de processamento não foi enviada")
	}

}

func TestCieloStoneWorker(t *testing.T) {
	t.Parallel()

	var w *CieloAcquirerWorkers
	var r *AcquirerActorsMock
	var p *AcquirerParameter
	var c config.Configuration
	var s *HTTPRequestSenderMock
	var ar *AuthorizationRequest
	var repo *CardRepositoryFinderMock

	r = &AcquirerActorsMock{}

	c.CieloAcquirer.URL = "htttp://localhost/acquirer/stone"
	s = newHTTPRequestSenderMock()
	repo = newCardRepositoryFinderMock()
	p = &AcquirerParameter{
		httpSender: s,
		cardFinder: repo,
	}

	ar = &AuthorizationRequest{ResponseChannel: make(chan *AuthorizationResponse), Transaction: &ExternalTransactionDTO{TransactionDTO: &TransactionDTO{CardOpenInfoDTO: &CardOpenInfoDTO{Holder: "João"}}}}

	w = NewCieloAcquirerWorkers(r, p, c)
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

	select {
	case <-ar.ResponseChannel:
		t.Log("Resposta de processamento enviada")
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Resposta de processamento não foi enviada")
	}

}

type HTTPRequestSenderMock struct {
	Sent chan bool
}

func newHTTPRequestSenderMock() (s *HTTPRequestSenderMock) {
	s = new(HTTPRequestSenderMock)
	s.Sent = make(chan bool, 1)
	return
}

func (s *HTTPRequestSenderMock) Send(url string, body interface{}, response interface{}) (err error) {

	s.Sent <- true

	return
}

type requestParam struct {
	url  string
	body *AcquirerTransactionDTO
}

type HTTPRequestSenderCaseMock struct {
	chParam chan *requestParam
	f       func(p *requestParam, chp chan *requestParam) (err error)
}

func newHTTPRequestSenderCaseMock(f func(p *requestParam, chp chan *requestParam) (err error)) (s *HTTPRequestSenderCaseMock) {
	s = new(HTTPRequestSenderCaseMock)
	s.chParam = make(chan *requestParam, 1)
	s.f = f
	return
}

func (s *HTTPRequestSenderCaseMock) Send(url string, body interface{}, response interface{}) (err error) {
	return s.f(&requestParam{url, body.(*AcquirerTransactionDTO)}, s.chParam)
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

type CardRepositoryFinderCaseMock struct {
	chParam chan string
	f       func(token string, chParam chan string) (c *Card, err error)
}

func newCardRepositoryFinderCaseMock(f func(token string, chParam chan string) (c *Card, err error)) (r *CardRepositoryFinderCaseMock) {
	r = new(CardRepositoryFinderCaseMock)
	r.chParam = make(chan string, 1)
	r.f = f
	return
}

func (r *CardRepositoryFinderCaseMock) Find(token string) (c *Card, err error) {

	return r.f(token, r.chParam)
}

type CardRepositoryFinderMock struct {
	Found chan bool
}

func newCardRepositoryFinderMock() (r *CardRepositoryFinderMock) {
	r = new(CardRepositoryFinderMock)
	r.Found = make(chan bool, 1)
	return
}

func (r *CardRepositoryFinderMock) Find(token string) (c *Card, err error) {

	r.Found <- true

	c = &Card{}

	return
}

func newAuthorizationRequest(token, holder string, total float64, installments int) (ar *AuthorizationRequest) {
	ar = new(AuthorizationRequest)
	ar.ResponseChannel = make(chan *AuthorizationResponse, 1)
	ar.Transaction = &ExternalTransactionDTO{
		Token: token,
		TransactionDTO: &TransactionDTO{
			CardOpenInfoDTO: &CardOpenInfoDTO{Holder: holder},
			PurchaseDTO:     &PurchaseDTO{Total: total, Installments: installments},
		},
	}
	return
}
