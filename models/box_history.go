package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type BoxHistory struct {
	ID        int
	BoxID     int
	Note      string
	CreatedBy int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Box           *Box
	CreatedByUser *User `gorm:":CreatedBy"`
}

type BoxHistoryReturnJSON struct {
	ID        int            `json:"id"`
	BoxID     int            `json:"box_id"`
	Note      string         `json:"note"`
	CreatedBy UserReturnJSON `json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
}

func (bh *BoxHistory) ConvertToBoxHistoryReturnJSON(ctx context.Context) *BoxHistoryReturnJSON {
	return &BoxHistoryReturnJSON{
		ID:        bh.ID,
		BoxID:     bh.BoxID,
		Note:      bh.Note,
		CreatedBy: *bh.CreatedByUser.ConvertToReturnJSON(ctx),
		CreatedAt: bh.CreatedAt,
	}
}

func (bh *BoxHistory) Create(ctx context.Context) error {
	result := util.DBFromContext(ctx).Create(&bh)
	return result.Error
}
