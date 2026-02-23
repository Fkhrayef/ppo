package lms

type Loan struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Status      string `json:"status"`
	PaidAmount  int64  `json:"paid_amount"`
	TotalAmount int64  `json:"total_amount"`
}

type Installment struct {
	ID      string `json:"id"`
	LoanID  string `json:"loan_id"`
	Amount  int64  `json:"amount"`
	Status  string `json:"status"`
	DueDate string `json:"due_date"`
}

type RecordPaymentRequest struct {
	LoanID        string `json:"loan_id"`
	InstallmentID string `json:"installment_id"`
	Amount        int64  `json:"amount"`
	TransactionID string `json:"transaction_id"`
}
