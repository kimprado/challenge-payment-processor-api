package processor

import "time"

// TransactionDTO representa dados de uma transação
// para processamento.
type TransactionDTO struct {
	CardDTO
	PurchaseDTO
	MerchantDTO
}

// CardDTO representa dados do cartão.
type CardDTO struct {
	Token    string
	Holder   string
	Brand    string
	Validity time.Time
}

// PurchaseDTO representa dados da compra.
type PurchaseDTO struct {
	Total        float64
	Installments int
	Items        []TransactionItemDTO
}

// TransactionItemDTO representa dados de itens de uma transação.
type TransactionItemDTO struct {
	Description string
	Price       float64
}

// MerchantDTO representa dados do vendedor.
type MerchantDTO struct {
	ID      string
	Address string
	Zipcode string
}
