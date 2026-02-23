package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/example/ppo/internal/client/lms"
	"github.com/example/ppo/internal/client/psp"
	"github.com/example/ppo/internal/order"
)

type Scheduler struct {
	cron      *cron.Cron
	lmsClient lms.Client
	pspClient psp.Client
	orderRepo order.Repository
	logger    *slog.Logger
}

func New(
	lmsClient lms.Client,
	pspClient psp.Client,
	orderRepo order.Repository,
	logger *slog.Logger,
) *Scheduler {
	return &Scheduler{
		cron:      cron.New(cron.WithSeconds()),
		lmsClient: lmsClient,
		pspClient: pspClient,
		orderRepo: orderRepo,
		logger:    logger,
	}
}

func (s *Scheduler) Start() error {
	if _, err := s.cron.AddFunc("0 0 9 * * *", s.sendReminders); err != nil {
		return err
	}

	if _, err := s.cron.AddFunc("0 0 2 * * *", s.autoChargeOverdue); err != nil {
		return err
	}

	s.cron.Start()
	return nil
}

func (s *Scheduler) Stop() { s.cron.Stop() }

// sendReminders fetches installments due soon and logs them.
// In production this would call a notification service (push / SMS / email).
func (s *Scheduler) sendReminders() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	upcoming, err := s.lmsClient.GetUpcomingInstallments(ctx)
	if err != nil {
		s.logger.Error("failed to fetch upcoming installments", "error", err)
		return
	}

	for _, inst := range upcoming {
		s.logger.Info("payment reminder",
			"loan_id", inst.LoanID,
			"installment_id", inst.ID,
			"due_date", inst.DueDate,
			"amount", inst.Amount,
		)
	}
}

// autoChargeOverdue fetches overdue installments from LMS and
// charges the card we have stored on the corresponding order.
func (s *Scheduler) autoChargeOverdue() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	overdue, err := s.lmsClient.GetOverdueInstallments(ctx)
	if err != nil {
		s.logger.Error("failed to fetch overdue installments", "error", err)
		return
	}

	for _, inst := range overdue {
		s.processOverdueInstallment(ctx, inst)
	}
}

func (s *Scheduler) processOverdueInstallment(ctx context.Context, inst lms.Installment) {
	log := s.logger.With("loan_id", inst.LoanID, "installment_id", inst.ID)

	matched, err := s.orderRepo.FindByLoanID(ctx, inst.LoanID)
	if err != nil {
		log.Error("no order found for overdue installment", "error", err)
		return
	}

	chargeResp, err := s.pspClient.Charge(ctx, psp.ChargeRequest{
		Amount:    inst.Amount,
		Currency:  matched.Currency,
		CardToken: matched.CardToken,
	})
	if err != nil {
		log.Error("auto-charge failed", "error", err)
		return
	}

	if err := s.lmsClient.RecordPayment(ctx, lms.RecordPaymentRequest{
		LoanID:        inst.LoanID,
		InstallmentID: inst.ID,
		Amount:        inst.Amount,
		TransactionID: chargeResp.TransactionID,
	}); err != nil {
		log.Error("auto-charge succeeded but LMS update failed",
			"transaction_id", chargeResp.TransactionID,
			"error", err,
		)
		return
	}

	log.Info("auto-charge completed", "transaction_id", chargeResp.TransactionID)
}
