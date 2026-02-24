package psp

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// fakeClient returns static responses that match the agreed-upon API contract
// with the PSP team. Use this while the real PSP service is still in development.
type fakeClient struct {
	logger *slog.Logger
}

func NewFake(logger *slog.Logger) Client {
	return &fakeClient{logger: logger}
}

func (f *fakeClient) Charge(_ context.Context, req ChargeRequest) (*ChargeResponse, error) {
	txnID := fmt.Sprintf("fake-txn-%d", time.Now().UnixMilli())
	f.logger.Info("[FAKE PSP] Charge",
		"amount", req.Amount,
		"currency", req.Currency,
		"card_token", req.CardToken,
		"transaction_id", txnID,
	)
	return &ChargeResponse{
		TransactionID: txnID,
		Status:        "captured",
	}, nil
}

func (f *fakeClient) Refund(_ context.Context, req RefundRequest) (*RefundResponse, error) {
	refundID := fmt.Sprintf("fake-ref-%d", time.Now().UnixMilli())
	f.logger.Info("[FAKE PSP] Refund",
		"order_id", req.OrderID,
		"amount", req.Amount,
		"currency", req.Currency,
		"refund_id", refundID,
	)
	return &RefundResponse{
		RefundID: refundID,
		Status:   "refunded",
	}, nil
}
