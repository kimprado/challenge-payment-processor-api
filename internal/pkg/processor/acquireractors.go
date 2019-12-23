package processor

// AcquirerActorsSender responsável por buscar ator correspondente
// de Acquirer e enviar solicitação de processamento
type AcquirerActorsSender interface {
	Send(aid AcquirerID, r *AuthorizationRequest) (err error)
}

// AcquirerActorsResgister responsável por registrar novos atores
type AcquirerActorsResgister interface {
	Resgister(aid AcquirerID, chr chan *AuthorizationRequest) (err error)
}

// ActorsMap representa mapa de atores.
// Chave do mapa tipada com AcquirerID garante maior segurança
type ActorsMap map[AcquirerID]chan *AuthorizationRequest

// NewActorsMap cria instância de ActorsMap.
func NewActorsMap() (m ActorsMap) {
	m = make(map[AcquirerID]chan *AuthorizationRequest)
	return
}

// AcquirerActors representa lista de atores de adquirentes disponíveis
type AcquirerActors struct {
	actors ActorsMap
}

// NewAcquirerActors cria instância de AcquirerActors.
func NewAcquirerActors(m ActorsMap) (a *AcquirerActors) {
	a = new(AcquirerActors)
	a.actors = m
	return
}

// Send implementa AcquirerActorsSender.
// Envia requisição para canal do ator do Adquirente.
func (a *AcquirerActors) Send(aid AcquirerID, r *AuthorizationRequest) (err error) {
	actor, ok := a.actors[aid]
	if !ok {
		err = newAcquirerActorSendNotFoundError(aid)
		return
	}
	actor <- r
	return
}

// Resgister implementa AcquirerActorsResgister.
// Registra ator para Adquirente com identificação AcquirerID.
func (a *AcquirerActors) Resgister(aid AcquirerID, chr chan *AuthorizationRequest) (err error) {
	if _, ok := a.actors[aid]; ok {
		err = newAcquirerActorRegisterExistsError(aid)
		return
	}

	if chr == nil {
		err = newAcquirerActorRegisterChannelNilError(aid)
		return
	}

	a.actors[aid] = chr
	return
}
