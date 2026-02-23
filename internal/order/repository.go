package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/example/ppo/pkg/apperror"
)

type Repository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error
	FindByLoanID(ctx context.Context, loanID string) (*Order, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, order *Order) error {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		return apperror.NewInternal("creating order", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	var o Order
	err := r.db.WithContext(ctx).Preload("Items").First(&o, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.NewNotFound(fmt.Sprintf("order %s not found", id))
	}
	if err != nil {
		return nil, apperror.NewInternal("fetching order", err)
	}
	return &o, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error {
	res := r.db.WithContext(ctx).Model(&Order{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return apperror.NewInternal("updating order status", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.NewNotFound(fmt.Sprintf("order %s not found", id))
	}
	return nil
}

func (r *repository) FindByLoanID(ctx context.Context, loanID string) (*Order, error) {
	var o Order
	err := r.db.WithContext(ctx).Preload("Items").First(&o, "loan_id = ?", loanID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.NewNotFound(fmt.Sprintf("order with loan %s not found", loanID))
	}
	if err != nil {
		return nil, apperror.NewInternal("fetching order by loan", err)
	}
	return &o, nil
}
