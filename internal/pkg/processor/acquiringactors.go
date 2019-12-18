package processor

//AcquiringActorsSender responsável por buscar ator correspondente
// de Acquirer e enviar solicitação de processamento
type AcquiringActorsSender interface {
	Send()
}
