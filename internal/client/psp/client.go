package psp

import "context"

type Client interface {
	Charge(ctx context.Context, req ChargeRequest) (*ChargeResponse, error)
	Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error)
}
