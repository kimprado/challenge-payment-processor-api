package processor

import "time"

// ExternalTransactionDTO representa dados de uma transação
// vinda do ambiente externo.
type ExternalTransactionDTO struct {
	*TransactionDTO
	Token string
}

// AcquirerTransactionDTO representa dados de uma transação
// enviada para processamento pelo Adquirente.
type AcquirerTransactionDTO struct {
	*TransactionDTO
	*CardHiddenInfoDTO
}

// TransactionDTO representa dados de uma transação
// para processamento.
type TransactionDTO struct {
	*CardOpenInfoDTO
	*PurchaseDTO
	*MerchantDTO
}

// CardHiddenInfoDTO representa dados sensíveis do cartão.
type CardHiddenInfoDTO struct {
	Number string
	CVV    string
}

// CardOpenInfoDTO representa abertos dados do cartão.
type CardOpenInfoDTO struct {
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
