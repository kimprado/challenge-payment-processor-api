package errors

import (
	"errors"
	"fmt"
)

// GenericError representa informações sobre erro
type GenericError struct {
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}

// NewGenericError cria instância de GenericError
func NewGenericError(title, detail string) (e *GenericError) {
	e = new(GenericError)
	e.Title = title
	e.Detail = detail
	return
}

func (e *GenericError) Error() string {
	return fmt.Sprintf("%s[%s]", e.Title, e.Detail)
}

// Is informa se target == e. Verifica se e é do tipo
// GenericError.
func (e *GenericError) Is(target error) bool {
	switch target.(type) {
	case *GenericError:
		return true
	default:
		return false
	}
}

// DomainError representa erros do domínio da aplicação
type DomainError struct {
	GenericError
}

// NewDomainError cria instância de DomainError
func NewDomainError(title, detail string) (e *DomainError) {
	e = new(DomainError)
	e.Title = title
	e.Detail = detail
	return
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s", e)
}

// ParametersError representa informações sobre erro de parâmetros
type ParametersError struct {
	Title             string           `json:"title"`
	Detail            string           `json:"detail,omitempty"`
	Instance          string           `json:"instance,omitempty"`
	InvalidParameters []ParameterError `json:"invalid-parameters,omitempty"`
}

// NewParametersError cria instância de ParametersError
func NewParametersError() (e *ParametersError) {
	e = new(ParametersError)
	e.Title = "Um ou Mais parâmetros não são válidos"
	e.InvalidParameters = []ParameterError{}
	return
}

func (e *ParametersError) Error() string {
	var instance string
	if e.Instance != "" {
		instance = fmt.Sprintf("(%v)", e.Instance)
	}

	return fmt.Sprintf("%s %v %v %v", e.Title, e.Detail, e.InvalidParameters, instance)
}

// Add adiciona novo ParameterError
func (e *ParametersError) Add(p ParameterError) {
	e.InvalidParameters = append(e.InvalidParameters, p)
}

// ContainsError informa se existe erros
func (e *ParametersError) ContainsError() bool {
	return len(e.InvalidParameters) > 0
}

// ParameterError representa informações sobre erro de parâmetros
type ParameterError struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Reason string `json:"reason"`
}

// GetDomainErrorOr retorna fromErr caso seja do tipo DomainError,
// caso contrário defaultErr.
func GetDomainErrorOr(fromErr error, defaultErr error) (e error) {
	if got := errors.Is(fromErr, &DomainError{}); got {
		e = fromErr
		return
	}
	e = defaultErr
	return
}
