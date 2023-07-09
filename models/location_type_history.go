package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
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

func (lth *LocationTypeHistory) ConvertToLocationTypeHistoryReturnJSON(ctx context.Context) *LocationTypeHistoryReturnJSON {
	return &LocationTypeHistoryReturnJSON{
		ID:             lth.ID,
		LocationTypeID: lth.LocationTypeID,
		Note:           lth.Note,
		CreatedBy:      *lth.CreatedByUser.ConvertToReturnJSON(ctx),
		CreatedAt:      lth.CreatedAt,
	}
}

func (lth *LocationTypeHistory) Create(ctx context.Context) error {
	result := util.DBFromContext(ctx).Create(&lth)
	return result.Error
}
