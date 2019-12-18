package processor

// AcquirerProcessor representa adquirente capaz de processar
// transação com cartão.
type AcquirerProcessor interface {
	Process()
}
