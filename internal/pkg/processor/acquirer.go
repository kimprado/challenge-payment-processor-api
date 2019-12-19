package processor

// AcquirerProcessor representa adquirente capaz de processar
// transação com cartão.
type AcquirerProcessor interface {
	Process()
}

// AcquirerID representa identificação de um Acquirer
type AcquirerID string
