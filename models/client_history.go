package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type ClientHistory struct {
	ID        int
	ClientID  int
	Note      string
	CreatedBy int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Client        *Client
	CreatedByUser *User `gorm:"foreignKey:CreatedBy"`
}

type ClientHistoryReturnJSON struct {
	ID        int            `json:"id"`
	ClientID  int            `json:"client_id"`
	Note      string         `json:"note"`
	CreatedBy UserReturnJSON `json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
}

func (ch *ClientHistory) ConvertToClientHistoryReturnJSON(ctx context.Context) *ClientHistoryReturnJSON {
	return &ClientHistoryReturnJSON{
		ID:        ch.ID,
		ClientID:  ch.ClientID,
		Note:      ch.Note,
		CreatedBy: *ch.CreatedByUser.ConvertToReturnJSON(ctx),
		CreatedAt: ch.CreatedAt,
	}
}

func (ch *ClientHistory) Create(ctx context.Context) error {
	result := util.DBFromContext(ctx).Create(&ch)
	return result.Error
}
