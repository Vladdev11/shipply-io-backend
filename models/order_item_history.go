package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderItemHistory struct {
	ID          int
	OrderItemID int
	Note        string
	CreatedBy   int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	OrderItem     *OrderItem
	CreatedByUser *User `gorm:"foreignKey:CreatedBy"`
}

type OrderItemHistoryReturnJSON struct {
	ID          int            `json:"id"`
	OrderItemID int            `json:"order_item_id"`
	Note        string         `json:"note"`
	CreatedBy   UserReturnJSON `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
}

func (oih *OrderItemHistory) ConvertToOrderItemHistoryReturnJSON() *OrderItemHistoryReturnJSON {
	return &OrderItemHistoryReturnJSON{
		ID:          oih.ID,
		OrderItemID: oih.OrderItemID,
		Note:        oih.Note,
		CreatedBy:   *oih.CreatedByUser.ConvertToReturnJSON(),
		CreatedAt:   oih.CreatedAt,
	}
}

func (oih *OrderItemHistory) Create() error {
	result := PGDB.Create(&oih)
	return result.Error
}
