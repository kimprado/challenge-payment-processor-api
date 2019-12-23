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
	case *errors.FriendlyError:
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
	case *errors.FriendlyError:
		return true
	default:
		return false
	}
}

// AcquirerActorRegisterChannelNilError representa erro channel nul no registro de ator
type AcquirerActorRegisterChannelNilError struct {
	*errors.GenericError
}

// newAcquirerActorRegisterChannelNilError cria instância de AcquirerActorRegisterChannelNilError
func newAcquirerActorRegisterChannelNilError(id AcquirerID) (e *AcquirerActorRegisterChannelNilError) {
	e = new(AcquirerActorRegisterChannelNilError)
	e.GenericError = errors.NewGenericError(
		"Falha ao registar ator para Adquirente",
		fmt.Sprintf("Ator %q não pode ter Channel nulo", id),
	)
	return
}

func (e *AcquirerActorRegisterChannelNilError) Error() string {
	return e.GenericError.Error()
}

// Is informa se target == e. Verifica se e é do tipo
// AcquirerActorRegisterChannelNilError, DomainError.
func (e *AcquirerActorRegisterChannelNilError) Is(target error) bool {
	switch target.(type) {
	case *AcquirerActorRegisterChannelNilError:
		return true
	case *errors.DomainError:
		return true
	case *errors.FriendlyError:
		return true
	default:
		return false
	}
}

// CardNotFoundError representa erro de cartão não encontrado
type CardNotFoundError struct {
	*errors.GenericError
}

// NewCardNotFoundError cria instância de CardNotFoundError
func NewCardNotFoundError() (e *CardNotFoundError) {
	e = new(CardNotFoundError)
	e.GenericError = errors.NewGenericError(
		"Falha ao consultar Cartão",
		"Cartão inexistente",
	)
	return
}

func (e *CardNotFoundError) Error() string {
	return e.GenericError.Error()
}

// Is informa se target == e. Verifica se e é do tipo
// CardNotFoundError, DomainError.
func (e *CardNotFoundError) Is(target error) bool {
	switch target.(type) {
	case *CardNotFoundError:
		return true
	case *errors.DomainError:
		return true
	case *errors.FriendlyError:
		return true
	default:
		return false
	}
}

// AcquirerValidationError representa erro de validação no adquirente
type AcquirerValidationError struct {
	*errors.ParametersError
}

// NewAcquirerValidationError cria instância de AcquirerValidationError
func NewAcquirerValidationError(message, url string) (e *AcquirerValidationError) {
	e = new(AcquirerValidationError)
	e.ParametersError = errors.NewParametersError()
	e.ParametersError.Title = "Falha no Adquirente ao Processar Transação"
	e.ParametersError.Detail = message
	e.ParametersError.Instance = url
	return
}

func (e *AcquirerValidationError) Error() string {
	return e.ParametersError.Error()
}

// Is informa se target == e. Verifica se e é do tipo
// AcquirerValidationError, DomainError.
func (e *AcquirerValidationError) Is(target error) bool {
	switch target.(type) {
	case *AcquirerValidationError:
		return true
	case *errors.DomainError:
		return true
	case *errors.FriendlyError:
		return true
	default:
		return false
	}
}

// AcquirerProcessingError representa erro no processamento do adquirente
type AcquirerProcessingError struct {
	*errors.GenericError
}

// NewAcquirerProcessingError cria instância de AcquirerProcessingError
func NewAcquirerProcessingError() (e *AcquirerProcessingError) {
	e = new(AcquirerProcessingError)
	e.GenericError = &errors.GenericError{}
	e.GenericError.Title = "Falha no Adquirente ao Processar Transação"
	return
}

func (e *AcquirerProcessingError) Error() string {
	return e.GenericError.Error()
}

// Is informa se target == e. Verifica se e é do tipo
// AcquirerProcessingError, ComponentError.
func (e *AcquirerProcessingError) Is(target error) bool {
	switch target.(type) {
	case *AcquirerProcessingError:
		return true
	case *errors.ComponentError:
		return true
	case *errors.FriendlyError:
		return true
	default:
		return false
	}
}

// AcquirerConnectivityError representa erro de conectividade com adquirente
type AcquirerConnectivityError struct {
	*errors.GenericError
}

// NewAcquirerConnectivityError cria instância de AcquirerConnectivityError
func NewAcquirerConnectivityError(message string, err error) (e *AcquirerConnectivityError) {
	e = new(AcquirerConnectivityError)
	e.GenericError = &errors.GenericError{}
	e.GenericError.Title = "Falha no Adquirente ao Processar Transação"
	e.GenericError.Detail = message
	e.GenericError.Err = err
	return
}

func (e *AcquirerConnectivityError) Error() string {
	return e.GenericError.Error()
}

// Is informa se target == e. Verifica se e é do tipo
// AcquirerConnectivityError, ComponentError.
func (e *AcquirerConnectivityError) Is(target error) bool {
	switch target.(type) {
	case *AcquirerConnectivityError:
		return true
	case *errors.ComponentError:
		return true
	case *errors.FriendlyError:
		return true
	default:
		return false
	}
}

// PaymentProcessError representa erro no Processamento da Transação
type PaymentProcessError struct {
	*errors.GenericError
}

// newPaymentProcessError cria instância de PaymentProcessError
func newPaymentProcessError() (e *PaymentProcessError) {
	e = new(PaymentProcessError)
	e.GenericError = errors.NewGenericError(
		"Falha no processamento da transação",
		fmt.Sprintf("Não foi possível Processar a Transação"),
	)
	return
}

func (e *PaymentProcessError) Error() string {
	return e.GenericError.Error()
}
