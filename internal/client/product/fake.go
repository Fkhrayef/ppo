package product

import (
	"context"
	"log/slog"
)

// fakeClient returns static responses that match the agreed-upon API contract
// with the Product team. Use this while the real Product service is still in development.
type fakeClient struct {
	logger *slog.Logger
}

func NewFake(logger *slog.Logger) Client {
	return &fakeClient{logger: logger}
}

func (f *fakeClient) RestockItem(_ context.Context, productID string, quantity int) error {
	f.logger.Info("[FAKE PRODUCT] RestockItem",
		"product_id", productID,
		"quantity", quantity,
	)
	return nil
}
