package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/challenge/payment-processor/internal/pkg/commom/errors"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/challenge/payment-processor/internal/pkg/commom/web"
	"github.com/challenge/payment-processor/internal/pkg/processor"
	"github.com/julienschmidt/httprouter"
)

// Controller trata requisições http de paredão
type Controller struct {
	processor processor.Processor
	logger    logging.LoggerAPI
}

// NewController é responsável por instanciar Controller
func NewController(p processor.Processor, l logging.LoggerAPI) (c *Controller) {
	c = new(Controller)
	c.processor = p
	c.logger = l
	return
}

// Process realiza processamento da transação
func (c *Controller) Process(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var errDecode error

	var ar *processor.AuthorizationResponse
	var a processor.AcquirerID
	var dto processor.ExternalTransactionDTO

	paramErr := errors.NewParametersError()

	a = (processor.AcquirerID)(r.Header.Get("X-ACQUIRER-ID"))

	if a == "" {
		paramErr.Add(
			errors.ParameterError{
				Name:   "X-ACQUIRER-ID",
				Value:  "",
				Reason: "Header 'X-ACQUIRER-ID' não pode ser vazio",
			},
		)
	}

	if r.Body == nil {
		paramErr.Add(
			errors.ParameterError{
				Name:   "body",
				Value:  "",
				Reason: "'body' não pode ser vazio",
			},
		)
	} else {
		errDecode = json.NewDecoder(r.Body).Decode(&dto)
	}

	if errDecode != nil {
		paramErr.Add(
			errors.ParameterError{
				Name:   "body",
				Value:  "",
				Reason: fmt.Sprintf("Não foi possivel converter parâmetro JSON"),
			},
		)
	}

	if paramErr.ContainsError() {
		c.logger.Warnf("Processar transação: %v\n", paramErr)

		web.NewHTTPResponse(
			w,
			statusCode(paramErr),
			nil,
			paramErr,
		).WriteJSON()

		return
	}

	ar = c.processor.Process(a, &dto)

	if ar.Err != nil {

		web.NewHTTPResponse(
			w,
			statusCode(ar.Err),
			nil,
			ar.Err,
		).WriteJSON()

		return
	}

	web.NewHTTPResponse(
		w,
		http.StatusOK,
		ar.Authorized,
		nil,
	).WriteJSON()

}
