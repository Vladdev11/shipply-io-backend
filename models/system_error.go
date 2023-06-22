package models

import "time"

type SystemError struct {
	ID             int64
	Message        string
	ClientID       int
	OrganizationID int
	CreatedAt      time.Time
}

func (SystemError) Create() error {
	return PGDB.Create(&SystemError{}).Error
}

func CreateSystemError(message string) {
	PGDB.Create(&SystemError{Message: message})
}
