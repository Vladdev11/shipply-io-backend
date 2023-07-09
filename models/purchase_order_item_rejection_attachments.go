package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/api"
	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type PurchaseOrderItemRejectionAttachment struct {
	ID                           int
	PurchaseOrderItemRejectionID int
	AttachmentID                 int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Attachment                 *Attachment
	PurchaseOrderItemRejection *PurchaseOrderItemRejection
}

type PurchaseOrderItemRejectionAttachmentReturnJSON struct {
	URL string `json:"url"`
}

func (poira *PurchaseOrderItemRejectionAttachment) Create(ctx context.Context) error {
	result := util.DBFromContext(ctx).Create(&poira)
	return result.Error
}

func (poira *PurchaseOrderItemRejectionAttachment) ConvertToReturnJSON(ctx context.Context) PurchaseOrderItemRejectionAttachmentReturnJSON {

	attachment, err := GetAttachmentByID(ctx, poira.AttachmentID)
	if err != nil {
		return PurchaseOrderItemRejectionAttachmentReturnJSON{}
	}

	URL, err := api.S3FromContext(ctx).GetAttachmentURL(attachment.UUID, attachment.Extension, attachment.FileName)
	if err != nil {
		return PurchaseOrderItemRejectionAttachmentReturnJSON{}
	}

	return PurchaseOrderItemRejectionAttachmentReturnJSON{
		URL: URL,
	}
}
