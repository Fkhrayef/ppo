package postpurchase

type PayInstallmentRequest struct {
	LoanID        string `json:"loan_id" binding:"required"`
	InstallmentID string `json:"installment_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,len=3"`
	CardToken     string `json:"card_token" binding:"required"`
}

type PayInstallmentResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}
