package processor

// AcquirerActorsSender responsável por buscar ator correspondente
// de Acquirer e enviar solicitação de processamento
type AcquirerActorsSender interface {
	Send()
}

// AcquirerActorsResgister responsável por registrar novos atores
type AcquirerActorsResgister interface {
	Resgister()
}
