// +build test unit

package processor

import (
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
