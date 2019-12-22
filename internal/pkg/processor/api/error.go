package api

import (
	"fmt"
	"net/http"

	"github.com/challenge/payment-processor/internal/pkg/processor"

	"github.com/challenge/payment-processor/internal/pkg/commom/errors"
)

func statusCode(e error) (s int) {

	switch v := e.(type) {
	case *processor.CardNotFoundError:
		s = http.StatusBadRequest // 404
	case *errors.ParametersError:
		s = http.StatusBadRequest // 400
	default:
		// TODO: Fazer com que este código não seja executado.
		// Adotar alguma das seguintes medidas:
		// 	- Aplicar testes/verificações.
		// 	- Eventualmente substituir por log.err, caso não exista.
		// 	- cobertura suficiente que garanta que não será executado.
		panic(fmt.Sprintf("Tipo de erro não definido %T", v))
	}
	return
}
