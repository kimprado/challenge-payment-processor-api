package processor

import (
	"errors"

	commomerrors "github.com/challenge/payment-processor/internal/pkg/commom/errors"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
)

// Processor representa ponto de entrada para comportamento
// da aplicação. O controlador do domínio.
type Processor interface {
	Process(a AcquirerID, t *ExternalTransactionDTO) (ar *AuthorizationResponse)
}

// PaymentProcessorService implementa Processor e é ponto de entrada
// para domínio da aplicação.
type PaymentProcessorService struct {
	actors AcquirerActorsSender
	logger logging.LoggerProcessor
}

// NewPaymentProcessorService cria instância de Sevice.
func NewPaymentProcessorService(a AcquirerActorsSender, l logging.LoggerProcessor) (p *PaymentProcessorService) {
	p = new(PaymentProcessorService)
	p.actors = a
	p.logger = l
	return
}

// Process delega processamento da transação para Acquirer.
func (p *PaymentProcessorService) Process(a AcquirerID, t *ExternalTransactionDTO) (ar *AuthorizationResponse) {
	r := &AuthorizationRequest{
		Transaction:     t,
		ResponseChannel: make(chan *AuthorizationResponse, 1),
	}
	p.actors.Send(a, r)
	ar = <-r.ResponseChannel
	if ar.Err != nil {
		if errors.Is(ar.Err, &commomerrors.DomainError{}) {
			p.logger.Warnf("Erro ao realizar transação: %v\n", ar.Err)
		} else {
			p.logger.Errorf("Erro ao realizar transação: %v\n", ar.Err)
		}

		// Garante que erro retornado seja tratado, amigável
		_, ar.Err = commomerrors.GetFriendlyErrorOr(ar.Err, NewPaymentProcessError())
	}
	return
}

// AuthorizationRequest representa dados de transação.
type AuthorizationRequest struct {
	Transaction     *ExternalTransactionDTO
	ResponseChannel chan *AuthorizationResponse
}

// AuthorizationResponse representa dados de transação.
type AuthorizationResponse struct {
	Authorized *AuthorizationMessage
	Err        error
}

// AuthorizationMessage representa mensagem de sucesso na autorização
type AuthorizationMessage struct {
	Message string `json:"message,omitempty"`
}
