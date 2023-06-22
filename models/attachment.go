package models

import (
	"time"

	"gorm.io/gorm"
)

type Attachment struct {
	ID        int
	UUID      string
	FileName  string
	Extension string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (a *Attachment) Delete() error {

	err := PGDB.Delete(a).Error
	if err != nil {
		return err
	}

	return nil
}

func CreateAttachment(attachment *Attachment) (*Attachment, error) {

	err := PGDB.Create(attachment).Error
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func GetAttachmentByID(id int) (*Attachment, error) {

	var attachment Attachment
	err := PGDB.Where("id = ?", id).First(&attachment).Error
	if err != nil {
		return nil, err
	}

	return &attachment, nil
}
