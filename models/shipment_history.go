package models

import (
	"time"

	"gorm.io/gorm"
)

type ShipmentHistory struct {
	ID         int
	ShipmentID int
	Note       string
	CreatedBy  int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Shipment      *Shipment
	CreatedByUser *User `gorm:"foreignKey:CreatedBy"`
}

type ShipmentHistoryReturnJSON struct {
	ID         int            `json:"id"`
	ShipmentID int            `json:"shipment_id"`
	Note       string         `json:"note"`
	CreatedBy  UserReturnJSON `json:"created_by"`
	CreatedAt  time.Time      `json:"created_at"`
}

func (sh *ShipmentHistory) ConvertToShipmentHistoryReturnJSON() *ShipmentHistoryReturnJSON {
	return &ShipmentHistoryReturnJSON{
		ID:         sh.ID,
		ShipmentID: sh.ShipmentID,
		Note:       sh.Note,
		CreatedBy:  *sh.CreatedByUser.ConvertToReturnJSON(),
		CreatedAt:  sh.CreatedAt,
	}
}

func (sh *ShipmentHistory) Create() error {
	result := PGDB.Create(&sh)
	return result.Error
}
