package processor

// Processor representa ponto de entrada para comportamento
// da aplicação. O controlador do domínio.
type Processor interface {
	Process()
}
