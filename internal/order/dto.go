package order

import "github.com/google/uuid"

type CreateRequest struct {
	UserID      uuid.UUID         `json:"user_id" binding:"required"`
	LoanID      string            `json:"loan_id" binding:"required"`
	CardToken   string            `json:"card_token" binding:"required"`
	Currency    string            `json:"currency" binding:"required,len=3"`
	TotalAmount int64             `json:"total_amount" binding:"required,gt=0"`
	Items       []CreateItemInput `json:"items" binding:"required,min=1,dive"`
}

type CreateItemInput struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
	UnitPrice int64  `json:"unit_price" binding:"required,gt=0"`
}

type Response struct {
	ID          uuid.UUID      `json:"id"`
	UserID      uuid.UUID      `json:"user_id"`
	LoanID      string         `json:"loan_id"`
	Status      Status         `json:"status"`
	TotalAmount int64          `json:"total_amount"`
	Currency    string         `json:"currency"`
	Items       []ItemResponse `json:"items"`
	CreatedAt   string         `json:"created_at"`
}

type ItemResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice int64     `json:"unit_price"`
}

func ToResponse(o *Order) Response {
	items := make([]ItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = ItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}
	return Response{
		ID:          o.ID,
		UserID:      o.UserID,
		LoanID:      o.LoanID,
		Status:      o.Status,
		TotalAmount: o.TotalAmount,
		Currency:    o.Currency,
		Items:       items,
		CreatedAt:   o.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
