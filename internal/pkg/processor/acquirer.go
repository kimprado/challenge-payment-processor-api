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

// StoneAcquirer representa Adquirente Stone
type StoneAcquirer struct {
	Acquirer
}

// CieloAcquirer representa Adquirente Cielo
type CieloAcquirer struct {
	Acquirer
}

// Acquirer implementa funcionalidades de de-para e envio
// da transação para Adquirente.
type Acquirer struct {
}

// AcquirerWorkers reprensenta trabalhadores que delegam
// trabalho para Acquirers
type AcquirerWorkers struct {
}
