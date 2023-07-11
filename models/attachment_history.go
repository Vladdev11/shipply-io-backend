package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type AttachmentHistory struct {
	ID           int
	AttachmentID int
	Note         string
	CreatedBy    int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Attachment    *Attachment
	CreatedByUser *User `gorm:":CreatedBy"`
}

type AttachmentHistoryReturnJSON struct {
	ID           int            `json:"id"`
	AttachmentID int            `json:"attachment_id"`
	Note         string         `json:"note"`
	CreatedBy    UserReturnJSON `json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
}

func (ah *AttachmentHistory) ConvertToAttachmentHistoryReturnJSON(ctx context.Context) *AttachmentHistoryReturnJSON {
	return &AttachmentHistoryReturnJSON{
		ID:           ah.ID,
		AttachmentID: ah.AttachmentID,
		Note:         ah.Note,
		CreatedBy:    *ah.CreatedByUser.ConvertToReturnJSON(ctx),
		CreatedAt:    ah.CreatedAt,
	}
}

func (ah *AttachmentHistory) Create(ctx context.Context) error {
	result := util.DBFromContext(ctx).Create(&ah)
	return result.Error
}
