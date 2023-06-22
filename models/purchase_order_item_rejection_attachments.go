package models

import (
	"time"

	"github.com/shipply-io/shipply-io-backend/api"
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

func (poira *PurchaseOrderItemRejectionAttachment) Create() error {
	result := PGDB.Create(&poira)
	return result.Error
}

func (poira *PurchaseOrderItemRejectionAttachment) ConvertToReturnJSON() PurchaseOrderItemRejectionAttachmentReturnJSON {

	attachment, err := GetAttachmentByID(poira.AttachmentID)
	if err != nil {
		return PurchaseOrderItemRejectionAttachmentReturnJSON{}
	}

	URL, err := api.GetAttachmentURL(attachment.UUID, attachment.Extension, attachment.FileName)
	if err != nil {
		return PurchaseOrderItemRejectionAttachmentReturnJSON{}
	}

	return PurchaseOrderItemRejectionAttachmentReturnJSON{
		URL: URL,
	}
}
