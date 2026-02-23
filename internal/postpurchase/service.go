package postpurchase

import (
	"context"
	"log/slog"

	"github.com/example/ppo/internal/client/lms"
	"github.com/example/ppo/internal/client/psp"
	"github.com/example/ppo/pkg/apperror"
)

type Service interface {
	GetInstallments(ctx context.Context, userID string) ([]lms.Installment, error)
	PayInstallment(ctx context.Context, req PayInstallmentRequest) (*PayInstallmentResponse, error)
}

type service struct {
	lmsClient lms.Client
	pspClient psp.Client
	logger    *slog.Logger
}

func NewService(lmsClient lms.Client, pspClient psp.Client, logger *slog.Logger) Service {
	return &service{
		lmsClient: lmsClient,
		pspClient: pspClient,
		logger:    logger,
	}
}

func (s *service) GetInstallments(ctx context.Context, userID string) ([]lms.Installment, error) {
	installments, err := s.lmsClient.GetInstallments(ctx, userID)
	if err != nil {
		return nil, apperror.NewUpstream("fetching installments from LMS", err)
	}
	return installments, nil
}

// PayInstallment charges the user's card then records the payment in LMS.
// If the PSP charge succeeds but LMS fails, we log the orphaned transaction
// for manual reconciliation (a real system would use an outbox/saga).
func (s *service) PayInstallment(ctx context.Context, req PayInstallmentRequest) (*PayInstallmentResponse, error) {
	chargeResp, err := s.pspClient.Charge(ctx, psp.ChargeRequest{
		Amount:    req.Amount,
		Currency:  req.Currency,
		CardToken: req.CardToken,
	})
	if err != nil {
		return nil, apperror.NewUpstream("charging via PSP", err)
	}

	if err := s.lmsClient.RecordPayment(ctx, lms.RecordPaymentRequest{
		LoanID:        req.LoanID,
		InstallmentID: req.InstallmentID,
		Amount:        req.Amount,
		TransactionID: chargeResp.TransactionID,
	}); err != nil {
		s.logger.Error("PSP charge succeeded but LMS recording failed — needs reconciliation",
			"loan_id", req.LoanID,
			"installment_id", req.InstallmentID,
			"transaction_id", chargeResp.TransactionID,
			"error", err,
		)
		return nil, apperror.NewUpstream("recording payment in LMS", err)
	}

	return &PayInstallmentResponse{
		TransactionID: chargeResp.TransactionID,
		Status:        "paid",
	}, nil
}
