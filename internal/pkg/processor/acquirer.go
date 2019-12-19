package processor

// AcquirerProcessor representa adquirente capaz de processar
// transação com cartão.
type AcquirerProcessor interface {
	Process()
}

// AcquirerTransactionMapper representa comportamento capaz de fazer de-para,
// e transformar dados para formato que adquirente espera.
type AcquirerTransactionMapper interface {
	MapTransaction()
}

// AcquirerID representa identificação de um Acquirer
type AcquirerID string
