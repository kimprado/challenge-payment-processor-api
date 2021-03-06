package processor

// ExternalTransactionDTO representa dados de uma transação
// vinda do ambiente externo.
type ExternalTransactionDTO struct {
	Token string `json:"token"`
	*TransactionDTO
}

// AcquirerTransactionDTO representa dados de uma transação
// enviada para processamento pela Adquirente.
type AcquirerTransactionDTO struct {
	*CardHiddenInfoDTO
	*TransactionDTO
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
	Number string `json:"number"`
	CVV    string `json:"cvv"`
}

// CardOpenInfoDTO representa abertos dados do cartão.
type CardOpenInfoDTO struct {
	Holder   string `json:"holder"`
	Brand    string `json:"brand"`
	Validity string `json:"validity"`
}

// PurchaseDTO representa dados da compra.
type PurchaseDTO struct {
	Total        float64              `json:"total"`
	Installments int                  `json:"installments"`
	Items        []TransactionItemDTO `json:"items"`
}

// TransactionItemDTO representa dados de itens de uma transação.
type TransactionItemDTO struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// MerchantDTO representa dados do vendedor.
type MerchantDTO struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Zipcode string `json:"zipcode"`
}
