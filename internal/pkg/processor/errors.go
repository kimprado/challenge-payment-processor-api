package processor

import (
	"fmt"

	"github.com/challenge/payment-processor/internal/pkg/commom/errors"
)

// AcquirerActorSendNotFoundError representa erro de AcquirerActor não encontrado para envio
type AcquirerActorSendNotFoundError struct {
	*errors.GenericError
}

// newAcquirerActorSendNotFoundError cria instância de AcquirerActorSendNotFoundError
func newAcquirerActorSendNotFoundError(id AcquirerID) (e *AcquirerActorSendNotFoundError) {
	e = new(AcquirerActorSendNotFoundError)
	e.GenericError = errors.NewGenericError(
		"Falha ao enviar requisição para ator",
		fmt.Sprintf("Ator %q inexistente", id),
	)
	return
}

func (e *AcquirerActorSendNotFoundError) Error() string {
	return e.GenericError.Error()
}

// Is informa se target == e. Verifica se e é do tipo
// AcquirerActorSendNotFoundError, DomainError.
func (e *AcquirerActorSendNotFoundError) Is(target error) bool {
	switch target.(type) {
	case *AcquirerActorSendNotFoundError:
		return true
	case *errors.DomainError:
		return true
	default:
		return false
	}
}
