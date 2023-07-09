package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type CarrierConnectionHistory struct {
	ID                  int
	CarrierConnectionID int
	Note                string
	CreatedBy           int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	CarrierConnection *CarrierConnection
	CreatedByUser     *User `gorm:"foreignKey:CreatedBy"`
}

type CarrierConnectionHistoryReturnJSON struct {
	ID                  int            `json:"id"`
	CarrierConnectionID int            `json:"carrier_connection_id"`
	Note                string         `json:"note"`
	CreatedBy           UserReturnJSON `json:"created_by"`
	CreatedAt           time.Time      `json:"created_at"`
}

func (cch *CarrierConnectionHistory) ConvertToCarrierConnectionHistoryReturnJSON(ctx context.Context) *CarrierConnectionHistoryReturnJSON {
	return &CarrierConnectionHistoryReturnJSON{
		ID:                  cch.ID,
		CarrierConnectionID: cch.CarrierConnectionID,
		Note:                cch.Note,
		CreatedBy:           *cch.CreatedByUser.ConvertToReturnJSON(ctx),
		CreatedAt:           cch.CreatedAt,
	}
}

func (cch *CarrierConnectionHistory) Create(ctx context.Context) error {
	result := util.DBFromContext(ctx).Create(&cch)
	return result.Error
}
