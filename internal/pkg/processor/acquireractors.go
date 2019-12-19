package processor

// AcquirerActorsSender responsável por buscar ator correspondente
// de Acquirer e enviar solicitação de processamento
type AcquirerActorsSender interface {
	Send(aid AcquirerID, a *AuthorizationRequest) (err error)
}

// AcquirerActorsResgister responsável por registrar novos atores
type AcquirerActorsResgister interface {
	Resgister()
}

// AcquirerActors representa lista de atores de adquirentes disponíveis
type AcquirerActors struct {
	// Chave do mapa tipada com AcquirerID garante maior segurança
	actors map[AcquirerID]chan *AuthorizationRequest
}

// NewAcquirerActors cria instância de AcquirerActors.
func NewAcquirerActors() (a *AcquirerActors) {
	a = new(AcquirerActors)
	a.actors = make(map[AcquirerID]chan *AuthorizationRequest)
	return
}

// Send implementa AcquirerActorsSender.
// Envia requisição para canal do ator do Adquirente.
func (a *AcquirerActors) Send(aid AcquirerID, ar *AuthorizationRequest) (err error) {
	actor, ok := a.actors[aid]
	if !ok {
		return
	}
	actor <- ar
	return
}
