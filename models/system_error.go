package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
)

type SystemError struct {
	ID             int64
	Message        string
	ClientID       int
	OrganizationID int
	CreatedAt      time.Time
}

func (SystemError) Create(ctx context.Context) error {
	return util.DBFromContext(ctx).Create(&SystemError{}).Error
}

func CreateSystemError(ctx context.Context, message string) {
	util.DBFromContext(ctx).Create(&SystemError{Message: message})
}
