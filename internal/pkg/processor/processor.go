package processor

// Processor representa ponto de entrada para comportamento
// da aplicação. O controlador do domínio.
type Processor interface {
	Process(t *TransactionDTO) (ar *AuthorizationResponse)
}

// Service implementa Processor é ponto de entrada
// para domínio da aplicação.
type Service struct {
}

// NewService cria instância de Sevice.
func NewService() (s *Service) {
	s = new(Service)
	return
}

// Process delega processamento da transação para Acquirer.
func (s *Service) Process(t *TransactionDTO) (ar *AuthorizationResponse) {
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
