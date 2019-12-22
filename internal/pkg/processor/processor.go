package processor

// Processor representa ponto de entrada para comportamento
// da aplicação. O controlador do domínio.
type Processor interface {
	Process(a AcquirerID, t *ExternalTransactionDTO) (ar *AuthorizationResponse)
}

// PaymentProcessorService implementa Processor e é ponto de entrada
// para domínio da aplicação.
type PaymentProcessorService struct {
	actors AcquirerActorsSender
}

// NewPaymentProcessorService cria instância de Sevice.
func NewPaymentProcessorService(a AcquirerActorsSender) (p *PaymentProcessorService) {
	p = new(PaymentProcessorService)
	p.actors = a
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
	Message string `json:"message"`
}
