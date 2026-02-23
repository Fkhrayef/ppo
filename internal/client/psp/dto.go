package psp

type ChargeRequest struct {
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
	CardToken string `json:"card_token"`
}

type ChargeResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

type RefundRequest struct {
	OrderID   string `json:"order_id"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
	CardToken string `json:"card_token"`
}

type RefundResponse struct {
	RefundID string `json:"refund_id"`
	Status   string `json:"status"`
}
