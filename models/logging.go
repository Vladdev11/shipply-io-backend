package models

import (
	"net/http"
	"time"
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

func (log *ShipEngineLog) Create() error {
	return PGDB.Create(log).Error
}
