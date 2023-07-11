package models

import (
	"context"
	"net/http"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
)

type ShipEngineLog struct {
	ID         int64
	EventType  string      // webhook or api_request
	Headers    http.Header `gorm:"type:jsonb"`
	Body       string
	Endpoint   string
	StatusCode int
	RequestURL string
	Method     string
	CreatedAt  time.Time
}

func (log *ShipEngineLog) Create(ctx context.Context) error {
	return util.DBFromContext(ctx).Create(log).Error
}
