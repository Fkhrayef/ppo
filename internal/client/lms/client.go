package lms

import "context"

type Client interface {
	GetLoan(ctx context.Context, loanID string) (*Loan, error)
	GetInstallments(ctx context.Context, userID string) ([]Installment, error)
	GetUpcomingInstallments(ctx context.Context) ([]Installment, error)
	GetOverdueInstallments(ctx context.Context) ([]Installment, error)
	UpdateLoanStatus(ctx context.Context, loanID, status string) error
	RecordPayment(ctx context.Context, req RecordPaymentRequest) error
}
