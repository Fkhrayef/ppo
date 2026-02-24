package lms

import (
	"context"
	"log/slog"
)

// fakeClient returns static responses that match the agreed-upon API contract
// with the LMS team. Use this while the real LMS service is still in development.
type fakeClient struct {
	logger *slog.Logger
}

func NewFake(logger *slog.Logger) Client {
	return &fakeClient{logger: logger}
}

func (f *fakeClient) GetLoan(_ context.Context, loanID string) (*Loan, error) {
	f.logger.Info("[FAKE LMS] GetLoan", "loan_id", loanID)
	return &Loan{
		ID:          loanID,
		UserID:      "user-aaa-bbb-ccc",
		Status:      "active",
		PaidAmount:  15000,
		TotalAmount: 60000,
	}, nil
}

func (f *fakeClient) GetInstallments(_ context.Context, userID string) ([]Installment, error) {
	f.logger.Info("[FAKE LMS] GetInstallments", "user_id", userID)
	return []Installment{
		{ID: "inst-001", LoanID: "loan-001", Amount: 15000, Status: "paid", DueDate: "2026-01-15"},
		{ID: "inst-002", LoanID: "loan-001", Amount: 15000, Status: "upcoming", DueDate: "2026-02-15"},
		{ID: "inst-003", LoanID: "loan-001", Amount: 15000, Status: "upcoming", DueDate: "2026-03-15"},
		{ID: "inst-004", LoanID: "loan-001", Amount: 15000, Status: "upcoming", DueDate: "2026-04-15"},
	}, nil
}

func (f *fakeClient) GetUpcomingInstallments(_ context.Context) ([]Installment, error) {
	f.logger.Info("[FAKE LMS] GetUpcomingInstallments")
	return []Installment{
		{ID: "inst-002", LoanID: "loan-001", Amount: 15000, Status: "upcoming", DueDate: "2026-02-15"},
		{ID: "inst-010", LoanID: "loan-002", Amount: 20000, Status: "upcoming", DueDate: "2026-02-20"},
	}, nil
}

func (f *fakeClient) GetOverdueInstallments(_ context.Context) ([]Installment, error) {
	f.logger.Info("[FAKE LMS] GetOverdueInstallments")
	return []Installment{
		{ID: "inst-007", LoanID: "loan-003", Amount: 12000, Status: "overdue", DueDate: "2026-02-01"},
	}, nil
}

func (f *fakeClient) UpdateLoanStatus(_ context.Context, loanID, status string) error {
	f.logger.Info("[FAKE LMS] UpdateLoanStatus", "loan_id", loanID, "new_status", status)
	return nil
}

func (f *fakeClient) RecordPayment(_ context.Context, req RecordPaymentRequest) error {
	f.logger.Info("[FAKE LMS] RecordPayment",
		"loan_id", req.LoanID,
		"installment_id", req.InstallmentID,
		"amount", req.Amount,
		"transaction_id", req.TransactionID,
	)
	return nil
}
