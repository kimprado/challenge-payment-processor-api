package processor

// Processor representa ponto de entrada para comportamento
// da aplicação. O controlador do domínio.
type Processor interface {
	Process(t *TransactionDTO) (ar *AuthorizationResponse)
}

// PaymentProcessorService implementa Processor e é ponto de entrada
// para domínio da aplicação.
type PaymentProcessorService struct {
}

// NewService cria instância de Sevice.
func NewService() (s *PaymentProcessorService) {
	s = new(PaymentProcessorService)
	return
}

// Process delega processamento da transação para Acquirer.
func (s *PaymentProcessorService) Process(t *TransactionDTO) (ar *AuthorizationResponse) {
	return
}

// AuthorizationRequest representa dados de transação.
type AuthorizationRequest struct {
	Transaction     *TransactionDTO
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
