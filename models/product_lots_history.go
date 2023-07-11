package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type ProductLotsHistory struct {
	ID            int
	ProductLotsID int
	Note          string
	CreatedBy     int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	ProductLot    *ProductLot
	CreatedByUser *User `gorm:"foreignKey:CreatedBy"`
}

type ProductLotsHistoryReturnJSON struct {
	ID            int            `json:"id"`
	ProductLotsID int            `json:"product_lots_id"`
	Note          string         `json:"note"`
	CreatedBy     UserReturnJSON `json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
}

func (plh *ProductLotsHistory) ConvertToProductLotsHistoryReturnJSON(ctx context.Context) *ProductLotsHistoryReturnJSON {
	return &ProductLotsHistoryReturnJSON{
		ID:            plh.ID,
		ProductLotsID: plh.ProductLotsID,
		Note:          plh.Note,
		CreatedBy:     *plh.CreatedByUser.ConvertToReturnJSON(ctx),
		CreatedAt:     plh.CreatedAt,
	}
}

func (plh *ProductLotsHistory) Create(ctx context.Context) error {
	result := util.DBFromContext(ctx).Create(&plh)
	return result.Error
}
