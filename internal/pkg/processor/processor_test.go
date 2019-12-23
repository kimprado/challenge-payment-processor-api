// +build test unit

package processor

import (
	"net/url"
	"testing"

	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/stretchr/testify/assert"
)

func TestProcessTransaction(t *testing.T) {
	t.Parallel()

	var p Processor
	var a AcquirerActorsSender

	l := logging.NewLoggerProcessor(map[string]string{
		"ROOT": "INFO",
	})

	a = &AcquirerActorsSenderMock{}
	p = NewPaymentProcessorService(a, l)

	assert.NotNil(t, p)

	result := p.Process("Stone", &ExternalTransactionDTO{})

	assert.NotNil(t, result)

}

func TestHandleDomainErrorProcessTransaction(t *testing.T) {
	t.Parallel()

	var p Processor
	var a AcquirerActorsSender

	l := logging.NewLoggerProcessor(map[string]string{
		"ROOT": "INFO",
	})

	a = &AcquirerActorsSenderDomainErrorMock{}
	p = NewPaymentProcessorService(a, l)

	assert.NotNil(t, p)

	result := p.Process("Stone", &ExternalTransactionDTO{})

	assert.NotNil(t, result)
	assert.Nil(t, result.Authorized)
	assert.NotNil(t, result.Err)
	assert.IsType(t, &CardNotFoundError{}, result.Err)

}

func TestHandleComponentErrorProcessTransaction(t *testing.T) {
	t.Parallel()

	var p Processor
	var a AcquirerActorsSender

	l := logging.NewLoggerProcessor(map[string]string{
		"ROOT": "INFO",
	})

	a = &AcquirerActorsSenderComponentErrorMock{}
	p = NewPaymentProcessorService(a, l)

	assert.NotNil(t, p)

	result := p.Process("Stone", &ExternalTransactionDTO{})

	assert.NotNil(t, result)
	assert.Nil(t, result.Authorized)
	assert.NotNil(t, result.Err)
	assert.IsType(t, &AcquirerConnectivityError{}, result.Err)

}

// Process deve converter erros genéricos como, net/url.Error, em PaymentProcessError.
func TestHandleGenericErrorProcessTransaction(t *testing.T) {
	t.Parallel()

	var p Processor
	var a AcquirerActorsSender

	l := logging.NewLoggerProcessor(map[string]string{
		"ROOT": "INFO",
	})

	a = &AcquirerActorsSenderGenericErrorMock{}
	p = NewPaymentProcessorService(a, l)

	assert.NotNil(t, p)

	result := p.Process("Stone", &ExternalTransactionDTO{})

	assert.NotNil(t, result)
	assert.Nil(t, result.Authorized)
	assert.NotNil(t, result.Err)

	// Process deve converter net/url.Error em PaymentProcessError.
	assert.IsType(t, &PaymentProcessError{}, result.Err)

}

type AcquirerActorsSenderMock struct {
}

func newAcquirerActorsSenderMock() (a *AcquirerActorsSenderMock) {
	a = new(AcquirerActorsSenderMock)
	return
}

func (a *AcquirerActorsSenderMock) Send(aid AcquirerID, r *AuthorizationRequest) (err error) {
	r.ResponseChannel <- &AuthorizationResponse{}
	return
}

type AcquirerActorsSenderDomainErrorMock struct {
}

func newAcquirerActorsSenderDomainErrorMock() (a *AcquirerActorsSenderDomainErrorMock) {
	a = new(AcquirerActorsSenderDomainErrorMock)
	return
}

func (a *AcquirerActorsSenderDomainErrorMock) Send(aid AcquirerID, r *AuthorizationRequest) (err error) {
	r.ResponseChannel <- &AuthorizationResponse{Err: &CardNotFoundError{}}
	return
}

type AcquirerActorsSenderComponentErrorMock struct {
}

func newAcquirerActorsSenderComponentErrorMock() (a *AcquirerActorsSenderComponentErrorMock) {
	a = new(AcquirerActorsSenderComponentErrorMock)
	return
}

func (a *AcquirerActorsSenderComponentErrorMock) Send(aid AcquirerID, r *AuthorizationRequest) (err error) {
	r.ResponseChannel <- &AuthorizationResponse{Err: &AcquirerConnectivityError{}}
	return
}

type AcquirerActorsSenderGenericErrorMock struct {
}

func newAcquirerActorsSenderGenericErrorMock() (a *AcquirerActorsSenderGenericErrorMock) {
	a = new(AcquirerActorsSenderGenericErrorMock)
	return
}

func (a *AcquirerActorsSenderGenericErrorMock) Send(aid AcquirerID, r *AuthorizationRequest) (err error) {
	r.ResponseChannel <- &AuthorizationResponse{Err: &url.Error{}}
	return
}
