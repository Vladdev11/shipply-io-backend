package models

import (
	"time"

	"gorm.io/gorm"
)

type PickSessionHistory struct {
	ID                 int `json:"id" gorm:"primary_key"`
	PickSessionID      int `json:"pick_session_id"`
	PickSessionOrderID int `json:"pick_session_order_id"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	Order            Order            `json:"order"`
	PickSession      PickSession      `json:"pick_session"`
	PickSessionOrder PickSessionOrder `json:"pick_session_order"`
}

func (psh *PickSessionHistory) Create() error {
	return PGDB.Create(psh).Error
}
