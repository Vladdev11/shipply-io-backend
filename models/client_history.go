package models

import (
	"time"

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

func (ch *ClientHistory) ConvertToClientHistoryReturnJSON() *ClientHistoryReturnJSON {
	return &ClientHistoryReturnJSON{
		ID:        ch.ID,
		ClientID:  ch.ClientID,
		Note:      ch.Note,
		CreatedBy: *ch.CreatedByUser.ConvertToReturnJSON(),
		CreatedAt: ch.CreatedAt,
	}
}

func (ch *ClientHistory) Create() error {
	result := PGDB.Create(&ch)
	return result.Error
}
