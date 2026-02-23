package order

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/example/ppo/internal/client/lms"
	"github.com/example/ppo/internal/client/product"
	"github.com/example/ppo/internal/client/psp"
	"github.com/example/ppo/pkg/apperror"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest) (*Order, error)
	Cancel(ctx context.Context, orderID uuid.UUID) error
}

type service struct {
	repo       Repository
	lmsClient  lms.Client
	pspClient  psp.Client
	prodClient product.Client
	logger     *slog.Logger
}

func NewService(
	repo Repository,
	lmsClient lms.Client,
	pspClient psp.Client,
	prodClient product.Client,
	logger *slog.Logger,
) Service {
	return &service{
		repo:       repo,
		lmsClient:  lmsClient,
		pspClient:  pspClient,
		prodClient: prodClient,
		logger:     logger,
	}
}

func (s *service) Create(ctx context.Context, req CreateRequest) (*Order, error) {
	items := make([]OrderItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = OrderItem{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice,
		}
	}

	o := &Order{
		UserID:      req.UserID,
		LoanID:      req.LoanID,
		Status:      StatusCreated,
		TotalAmount: req.TotalAmount,
		Currency:    req.Currency,
		CardToken:   req.CardToken,
		Items:       items,
	}

	if err := s.repo.Create(ctx, o); err != nil {
		return nil, err
	}

	return o, nil
}

// Cancel orchestrates a full cancellation:
// 1. Fetch order from DB
// 2. Call LMS to see how much the user actually paid
// 3. Refund via PSP if anything was paid
// 4. Mark loan as refunded in LMS
// 5. Restock every item via Product Service
// 6. Update local order status
func (s *service) Cancel(ctx context.Context, orderID uuid.UUID) error {
	o, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if o.Status == StatusCancelled || o.Status == StatusRefunded {
		return apperror.NewConflict(fmt.Sprintf("order %s is already %s", orderID, o.Status))
	}

	loan, err := s.lmsClient.GetLoan(ctx, o.LoanID)
	if err != nil {
		return apperror.NewUpstream("fetching loan from LMS", err)
	}

	if loan.PaidAmount > 0 {
		_, err = s.pspClient.Refund(ctx, psp.RefundRequest{
			OrderID:   orderID.String(),
			Amount:    loan.PaidAmount,
			Currency:  o.Currency,
			CardToken: o.CardToken,
		})
		if err != nil {
			return apperror.NewUpstream("refunding via PSP", err)
		}
	}

	if err := s.lmsClient.UpdateLoanStatus(ctx, o.LoanID, "refunded"); err != nil {
		s.logger.Error("failed to update loan status after refund", "loan_id", o.LoanID, "error", err)
		return apperror.NewUpstream("updating loan status in LMS", err)
	}

	for _, item := range o.Items {
		if err := s.prodClient.RestockItem(ctx, item.ProductID, item.Quantity); err != nil {
			s.logger.Error("failed to restock item", "product_id", item.ProductID, "error", err)
			return apperror.NewUpstream("restocking inventory", err)
		}
	}

	if err := s.repo.UpdateStatus(ctx, orderID, StatusRefunded); err != nil {
		return err
	}

	return nil
}
