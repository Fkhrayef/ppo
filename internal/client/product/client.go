package product

import "context"

type Client interface {
	RestockItem(ctx context.Context, productID string, quantity int) error
}
