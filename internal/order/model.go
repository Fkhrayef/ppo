package order

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusCreated   Status = "created"
	StatusActive    Status = "active"
	StatusCancelled Status = "cancelled"
	StatusRefunded  Status = "refunded"
)

type Order struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID   `gorm:"type:uuid;not null;index"`
	LoanID      string      `gorm:"type:varchar(64);not null"`
	Status      Status      `gorm:"type:varchar(20);not null;default:'created'"`
	TotalAmount int64       `gorm:"not null"`
	Currency    string      `gorm:"type:varchar(3);not null;default:'SAR'"`
	CardToken   string      `gorm:"type:varchar(255);not null"`
	Items       []OrderItem `gorm:"foreignKey:OrderID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index"`
	ProductID string    `gorm:"type:varchar(64);not null"`
	Quantity  int       `gorm:"not null"`
	UnitPrice int64     `gorm:"not null"`
	CreatedAt time.Time
}

func (Order) TableName() string     { return "orders" }
func (OrderItem) TableName() string { return "order_items" }
