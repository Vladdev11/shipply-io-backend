package models

import (
	"time"

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

func (plh *ProductLotsHistory) ConvertToProductLotsHistoryReturnJSON() *ProductLotsHistoryReturnJSON {
	return &ProductLotsHistoryReturnJSON{
		ID:            plh.ID,
		ProductLotsID: plh.ProductLotsID,
		Note:          plh.Note,
		CreatedBy:     *plh.CreatedByUser.ConvertToReturnJSON(),
		CreatedAt:     plh.CreatedAt,
	}
}

func (plh *ProductLotsHistory) Create() error {
	result := PGDB.Create(&plh)
	return result.Error
}
