package api

import (
	"net/http"

	"github.com/challenge/payment-processor/internal/pkg/processor"

	"github.com/challenge/payment-processor/internal/pkg/commom/errors"
)

func statusCode(e error) (s int) {

	switch e.(type) {
	case *processor.CardNotFoundError:
		s = http.StatusNotFound // 404
	case *errors.ParametersError:
		s = http.StatusBadRequest // 400
	case *processor.AcquirerProcessingError:
		s = http.StatusServiceUnavailable // 503
	case *processor.AcquirerConnectivityError:
		s = http.StatusServiceUnavailable // 503
	default:
		s = http.StatusInternalServerError // 500
	}
	return
}
