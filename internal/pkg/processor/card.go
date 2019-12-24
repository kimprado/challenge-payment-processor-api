package processor

// Card representa dados sensíveis de cartão
type Card struct {
	Number string `json:"number"`
	CVV    string `json:"cvv"`
}
