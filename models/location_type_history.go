package models

import (
	"time"

	"gorm.io/gorm"
)

type LocationTypeHistory struct {
	ID             int
	LocationTypeID int
	Note           string
	CreatedBy      int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	LocationType  *LocationType
	CreatedByUser *User `gorm:"foreignKey:CreatedBy"`
}

type LocationTypeHistoryReturnJSON struct {
	ID             int            `json:"id"`
	LocationTypeID int            `json:"location_type_id"`
	Note           string         `json:"note"`
	CreatedBy      UserReturnJSON `json:"created_by"`
	CreatedAt      time.Time      `json:"created_at"`
}

func (lth *LocationTypeHistory) ConvertToLocationTypeHistoryReturnJSON() *LocationTypeHistoryReturnJSON {
	return &LocationTypeHistoryReturnJSON{
		ID:             lth.ID,
		LocationTypeID: lth.LocationTypeID,
		Note:           lth.Note,
		CreatedBy:      *lth.CreatedByUser.ConvertToReturnJSON(),
		CreatedAt:      lth.CreatedAt,
	}
}

func (lth *LocationTypeHistory) Create() error {
	result := PGDB.Create(&lth)
	return result.Error
}
