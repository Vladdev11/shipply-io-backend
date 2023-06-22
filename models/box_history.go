package models

import (
	"time"

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

func (bh *BoxHistory) ConvertToBoxHistoryReturnJSON() *BoxHistoryReturnJSON {
	return &BoxHistoryReturnJSON{
		ID:        bh.ID,
		BoxID:     bh.BoxID,
		Note:      bh.Note,
		CreatedBy: *bh.CreatedByUser.ConvertToReturnJSON(),
		CreatedAt: bh.CreatedAt,
	}
}

func (bh *BoxHistory) Create() error {
	result := PGDB.Create(&bh)
	return result.Error
}
