package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
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

func (a *Attachment) Delete(ctx context.Context) error {

	err := util.DBFromContext(ctx).Delete(a).Error
	if err != nil {
		return err
	}

	return nil
}

func CreateAttachment(ctx context.Context, attachment *Attachment) (*Attachment, error) {

	err := util.DBFromContext(ctx).Create(attachment).Error
	if err != nil {
		return nil, ErrCreateFailed{Object: "attachment", Err: err}
	}

	return attachment, nil
}

func GetAttachmentByID(ctx context.Context, id int) (*Attachment, error) {

	var attachment Attachment
	err := util.DBFromContext(ctx).Where("id = ?", id).First(&attachment).Error
	if err != nil {
		return nil, err
	}

	return &attachment, nil
}
