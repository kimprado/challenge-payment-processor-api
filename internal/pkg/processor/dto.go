package processor

// TransactionDTO representa dados de uma transação
// para processamento.
type TransactionDTO struct {
	token string
	total float64
	itens []TransactionItemDTO
}

// TransactionItemDTO representa itens de uma transação.
type TransactionItemDTO struct {
	description  string
	installments int
}
