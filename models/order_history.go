package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type OrderHistory struct {
	ID        int
	OrderID   int
	Note      string
	CreatedBy int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Order         *Order
	CreatedByUser *User `gorm:"foreignKey:CreatedBy"`
}

type OrderHistoryReturnJSON struct {
	ID        int            `json:"id"`
	OrderID   int            `json:"order_id"`
	Note      string         `json:"note"`
	CreatedBy UserReturnJSON `json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
}

func (oh *OrderHistory) ConvertToOrderHistoryReturnJSON(ctx context.Context) *OrderHistoryReturnJSON {
	return &OrderHistoryReturnJSON{
		ID:        oh.ID,
		OrderID:   oh.OrderID,
		Note:      oh.Note,
		CreatedBy: *oh.CreatedByUser.ConvertToReturnJSON(ctx),
		CreatedAt: oh.CreatedAt,
	}
}

func (oh *OrderHistory) Create(ctx context.Context) error {
	result := util.DBFromContext(ctx).Create(&oh)
	return result.Error
}
