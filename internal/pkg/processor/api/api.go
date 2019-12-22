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
func NewController(p processor.Processor, l logging.LoggerAPI) (ctrl *Controller) {
	ctrl = new(Controller)
	ctrl.processor = p
	ctrl.logger = l
	return
}

// Exchange calcula câmbio
func (c *Controller) Exchange(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var err error

	paramErr := errors.NewParametersError()

	var dto processor.ExternalTransactionDTO
	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		c.logger.Errorf("Erro ao converter transação %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)

		paramErr.Add(
			errors.ParameterError{
				Name:   "body",
				Value:  "",
				Reason: fmt.Sprintf("Não foi possivel converter paraâmetro JSON"),
			},
		)
		return
	}

	if paramErr.ContainsError() {
		c.logger.Warnf("Consulta Exchange : %v\n", paramErr)

		web.NewHTTPResponse(
			w,
			statusCode(paramErr),
			nil,
			paramErr,
		).WriteJSON()

		return
	}

	var a processor.AcquirerID
	var ar *processor.AuthorizationResponse
	a = (processor.AcquirerID)(r.Header.Get("X-ACQUIRER-ID"))

	ar = c.processor.Process(a, &dto)

	if err != nil {
		c.logger.Warnf("Erro ao realizar transação: %v\n", err)

		web.NewHTTPResponse(
			w,
			statusCode(err),
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
