// +build test integration

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/challenge/payment-processor/internal/pkg/processor"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationProcessTransaction(t *testing.T) {
	t.Parallel()

	var err error

	err = setUpCards(t)
	if err != nil {
		t.Errorf("Erro ao preparar teste: %+v\n", err)
		return
	}

	c, err := initializeConfigTest()
	if err != nil {
		t.Errorf("Erro ao criar Configuração: %+v\n", err)
		return
	}

	c.RedisDB.Prefix = t.Name()

	startWorkers, err := initializeControllerWithDependenciesTest(c)
	if err != nil {
		t.Errorf("Conexão banco de dados %v\n", err)
		return
	}

	ctrl := startWorkers.ctrl
	assert.NotNil(t, ctrl)

	testCases := []struct {
		//label indica título do Test Case
		label        string
		a            *processor.AcquirerID
		t            *processor.ExternalTransactionDTO
		statusCode   int
		responseBody string
	}{
		{
			"Stone - Transação Válida",
			newAcquirerID_IT("Stone"),
			newExternalTransactionDTO_IT("xpto121a", "João", 1000, 1),
			http.StatusOK,
			`{"message":"Transação autorizada"}`,
		},
		{
			"Cielo - Transação Válida",
			newAcquirerID_IT("Cielo"),
			newExternalTransactionDTO_IT("xpto121a", "João", 500, 1),
			http.StatusOK,
			`{"message":"Transação autorizada"}`,
		},
		{
			"Stone - Transação Inválida - Valor inválido > 1000",
			newAcquirerID_IT("Stone"),
			newExternalTransactionDTO_IT("xpto121a", "João", 1001, 1),
			http.StatusBadRequest,
			`{"title":"Validação do Adquirente ao Processar Transação","detail":"Valor inválido"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {

			var err error

			router := httprouter.New()
			router.POST("/transactions", ctrl.Process)

			var body io.Reader

			if tc.t != nil {
				bb, _ := json.Marshal(tc.t)
				body = bytes.NewReader(bb)
			}

			req, err := http.NewRequest("POST", "/transactions", body)
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

func newAcquirerID_IT(id string) (a *processor.AcquirerID) {
	aux := (processor.AcquirerID)(id)
	a = &aux
	return
}

func newExternalTransactionDTO_IT(token, holder string, total float64, installments int) (t *processor.ExternalTransactionDTO) {
	t = &processor.ExternalTransactionDTO{
		Token: token,
		TransactionDTO: &processor.TransactionDTO{
			CardOpenInfoDTO: &processor.CardOpenInfoDTO{Holder: holder},
			PurchaseDTO:     &processor.PurchaseDTO{Total: total, Installments: installments},
		},
	}
	return
}

// setUpCards cria carga de dados para teste.
// Popula Redis com valores em nova chave para o teste.
// Nome da chave se baseia no nome do teste, o que permite
// executar testes de integração em paralelo :).
func setUpCards(t *testing.T) (err error) {

	loadCases := []struct {
		token string
		card  processor.Card
	}{
		{"xpto121a", processor.Card{Number: "121", CVV: "a"}},
		{"xpto122b", processor.Card{Number: "122", CVV: "b"}},
		{"xpto123c", processor.Card{Number: "123", CVV: "c"}},
	}

	c, err := initializeConfigTest()
	if err != nil {
		return
	}

	redis, err := initializeRedisTest(c)
	if err != nil {
		return
	}

	con := redis.Get()
	defer con.Close()

	err = con.Send("MULTI")
	if err != nil {
		return
	}

	for _, lc := range loadCases {
		var cardJSON []byte
		cardJSON, err = json.Marshal(lc.card)
		con.Send("SET", fmt.Sprintf("%v:card:%v", t.Name(), lc.token), string(cardJSON))
		if err != nil {
			return
		}
	}

	_, err = con.Do("EXEC")
	if err != nil {
		return
	}

	return
}
