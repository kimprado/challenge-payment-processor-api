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

// AcquirerActorRegisterExistsError representa erro de AcquirerActor existente
type AcquirerActorRegisterExistsError struct {
	*errors.GenericError
}

// newAcquirerActorRegisterExistsError cria instância de AcquirerActorRegisterExistsError
func newAcquirerActorRegisterExistsError(id AcquirerID) (e *AcquirerActorRegisterExistsError) {
	e = new(AcquirerActorRegisterExistsError)
	e.GenericError = errors.NewGenericError(
		"Falha ao registar ator para Adquirente",
		fmt.Sprintf("Ator %q já existe", id),
	)
	return
}

func (e *AcquirerActorRegisterExistsError) Error() string {
	return e.GenericError.Error()
}

// Is informa se target == e. Verifica se e é do tipo
// AcquirerActorRegisterExistsError, DomainError.
func (e *AcquirerActorRegisterExistsError) Is(target error) bool {
	switch target.(type) {
	case *AcquirerActorRegisterExistsError:
		return true
	case *errors.DomainError:
		return true
	default:
		return false
	}
}
