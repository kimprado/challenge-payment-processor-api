package processor

// AcquirerActorsSender responsável por buscar ator correspondente
// de Acquirer e enviar solicitação de processamento
type AcquirerActorsSender interface {
	Send()
}
