// +build test unit

package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/challenge/payment-processor/internal/pkg/processor"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestProcessTransaction(t *testing.T) {
	t.Parallel()

	l := logging.NewLoggerAPI(map[string]string{
		"ROOT": "INFO",
	})

	testCases := []struct {
		//label indica título do Test Case
		label              string
		a                  *processor.AcquirerID
		t                  *processor.ExternalTransactionDTO
		invalidTransaction *string
		p                  processor.Processor
		statusCode         int
		responseBody       string
		err                error
	}{
		{
			"Transação Válida",
			newAcquirerID("Stone"),
			newExternalTransactionDTO("xpto121a", "João", 1000, 1),
			nil,
			newProcessorCaseMock(func(a processor.AcquirerID, t *processor.ExternalTransactionDTO) (ar *processor.AuthorizationResponse) {
				ar = &processor.AuthorizationResponse{Authorized: &processor.AuthorizationMessage{Message: "Autorizada"}}
				return
			}),
			200,
			`{"message":"Autorizada"}`,
			nil,
		},
		{
			"AcquirerID Header nulo",
			nil,
			newExternalTransactionDTO("xpto121a", "João", 1000, 1),
			nil,
			nil,
			400,
			`{"title":"Um ou Mais parâmetros não são válidos","invalid-parameters":[{"name":"X-ACQUIRER-ID","value":"","reason":"'X-ACQUIRER-ID' não pode ser vazio"}]}`,
			nil,
		},
		{
			"AcquirerID Header Vazio",
			newAcquirerID(""),
			newExternalTransactionDTO("xpto121a", "João", 1000, 1),
			nil,
			nil,
			400,
			`{"title":"Um ou Mais parâmetros não são válidos","invalid-parameters":[{"name":"X-ACQUIRER-ID","value":"","reason":"'X-ACQUIRER-ID' não pode ser vazio"}]}`,
			nil,
		},
		{
			"Body vazio",
			newAcquirerID("Stone"),
			nil,
			nil,
			nil,
			400,
			`{"title":"Um ou Mais parâmetros não são válidos","invalid-parameters":[{"name":"body","value":"","reason":"'body' não pode ser vazio"}]}`,
			nil,
		},
		{
			"Body JSON inválido",
			newAcquirerID("Stone"),
			nil,
			newInvalidTransaction(`{"token": "abc":]}}`),
			nil,
			400,
			`{"title":"Um ou Mais parâmetros não são válidos","invalid-parameters":[{"name":"body","value":"","reason":"Não foi possivel converter parâmetro JSON"}]}`,
			nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {

			var err error

			handler := NewController(tc.p, l)

			router := httprouter.New()
			router.POST("/transactions/", handler.Process)

			var body io.Reader

			if tc.t != nil {
				bb, _ := json.Marshal(tc.t)
				body = bytes.NewReader(bb)
			}

			if tc.invalidTransaction != nil {
				body = bytes.NewReader([]byte(*tc.invalidTransaction))
			}

			req, err := http.NewRequest("POST", "/transactions/", body)
			if err != nil {
				t.Fatal(err)
			}

			if tc.a != nil {
				req.Header.Set("X-ACQUIRER-ID", tc.a.String())
			}

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.statusCode, rr.Code)
			if notNil := assert.NotNil(t, rr.Body); !notNil {
				return
			}
			assert.Equal(t, tc.responseBody, rr.Body.String()[:len(rr.Body.String())-1])
		})
	}

}

func newAcquirerID(id string) (a *processor.AcquirerID) {
	aux := (processor.AcquirerID)(id)
	a = &aux
	return
}

func newExternalTransactionDTO(token, holder string, total float64, installments int) (t *processor.ExternalTransactionDTO) {
	t = &processor.ExternalTransactionDTO{
		Token: token,
		TransactionDTO: &processor.TransactionDTO{
			CardOpenInfoDTO: &processor.CardOpenInfoDTO{Holder: holder},
			PurchaseDTO:     &processor.PurchaseDTO{Total: total, Installments: installments},
		},
	}
	return
}

type ProcessorCaseMock struct {
	f func(a processor.AcquirerID, t *processor.ExternalTransactionDTO) (ar *processor.AuthorizationResponse)
}

func newProcessorCaseMock(f func(a processor.AcquirerID, t *processor.ExternalTransactionDTO) (ar *processor.AuthorizationResponse)) (p *ProcessorCaseMock) {
	p = new(ProcessorCaseMock)
	p.f = f
	return
}

func (r *ProcessorCaseMock) Process(a processor.AcquirerID, t *processor.ExternalTransactionDTO) (ar *processor.AuthorizationResponse) {

	return r.f(a, t)
}

func newInvalidTransaction(t string) *string {
	return &t
}
