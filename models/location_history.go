package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type LocationHistory struct {
	ID         int
	LocationID int
	Note       string
	CreatedBy  int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Location      *Location
	CreatedByUser *User `gorm:"foreignKey:CreatedBy"`
}

type LocationHistoryReturnJSON struct {
	ID         int            `json:"id"`
	LocationID int            `json:"location_id"`
	Note       string         `json:"note"`
	CreatedBy  UserReturnJSON `json:"created_by"`
	CreatedAt  time.Time      `json:"created_at"`
}

func (lh *LocationHistory) ConvertToLocationHistoryReturnJSON(ctx context.Context) *LocationHistoryReturnJSON {
	return &LocationHistoryReturnJSON{
		ID:         lh.ID,
		LocationID: lh.LocationID,
		Note:       lh.Note,
		CreatedBy:  *lh.CreatedByUser.ConvertToReturnJSON(ctx),
		CreatedAt:  lh.CreatedAt,
	}
}

func (lh *LocationHistory) Create(ctx context.Context) error {
	result := util.DBFromContext(ctx).Create(&lh)
	return result.Error
}
