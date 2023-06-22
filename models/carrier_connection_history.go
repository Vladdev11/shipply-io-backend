package models

import (
	"time"

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

func (cch *CarrierConnectionHistory) ConvertToCarrierConnectionHistoryReturnJSON() *CarrierConnectionHistoryReturnJSON {
	return &CarrierConnectionHistoryReturnJSON{
		ID:                  cch.ID,
		CarrierConnectionID: cch.CarrierConnectionID,
		Note:                cch.Note,
		CreatedBy:           *cch.CreatedByUser.ConvertToReturnJSON(),
		CreatedAt:           cch.CreatedAt,
	}
}

func (cch *CarrierConnectionHistory) Create() error {
	result := PGDB.Create(&cch)
	return result.Error
}
