package processor

// Card representa dados sensíveis de cartão
type Card struct {
	Number string `json:"number"`
	CVV    string `json:"cvv"`
}

func newCard(number, cvv string) (c *Card) {
	c = new(Card)
	c.Number = number
	c.CVV = cvv
	return
}
