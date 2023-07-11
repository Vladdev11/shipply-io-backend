package models

import (
	"context"
	"mime/multipart"
	"net/http"
	"reflect"
	"time"

	"github.com/shipply-io/shipply-io-backend/api"
	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type PurchaseOrderAttachment struct {
	ID              int
	PurchaseOrderID int
	AttachmentID    int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Attachment    *Attachment
	PurchaseOrder *PurchaseOrder
}

type PurchaseOrderAttachmentCreateRequest struct {
	File     multipart.File `json:"file"`
	FileName string         `json:"file_name"`
}

type PurchaseOrderAttachmentReturnJSON struct {
	ID       int    `json:"id"`
	FileName string `json:"file_name"`
	URL      string `json:"url"`
}

func (poacr *PurchaseOrderAttachmentCreateRequest) ParseAndValidateRequest(r *http.Request) []string {
	errors := []string{}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errors = append(errors, err.Error())
		return errors
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		errors = append(errors, "file is required in a required multipart file format field")
		return errors
	}
	defer file.Close()
	poacr.File = file

	poacr.FileName = r.FormValue("file_name")
	if poacr.FileName == "" {
		poacr.FileName = fileHeader.Filename
	} else if reflect.TypeOf(poacr.FileName) != reflect.TypeOf("") {
		errors = append(errors, "file_name must be a string")
	} else if len(poacr.FileName) > 255 {
		errors = append(errors, "file_name must be less than 255 characters")
	}

	return errors
}

func (poa *PurchaseOrderAttachment) ConvertToReturnJSON(ctx context.Context) (PurchaseOrderAttachmentReturnJSON, error) {

	poa.GetAttachment(ctx)

	URL, err := api.S3FromContext(ctx).GetAttachmentURL(poa.Attachment.UUID, poa.Attachment.Extension, poa.Attachment.FileName)
	if err != nil {
		return PurchaseOrderAttachmentReturnJSON{}, err
	}

	return PurchaseOrderAttachmentReturnJSON{
		ID:       poa.ID,
		FileName: poa.Attachment.FileName,
		URL:      URL,
	}, nil
}

func (poa *PurchaseOrderAttachment) Delete(ctx context.Context) error {

	err := util.DBFromContext(ctx).Delete(poa).Error
	if err != nil {
		return err
	}

	return nil

}

func (poa *PurchaseOrderAttachment) GetAttachment(ctx context.Context) error {

	err := util.DBFromContext(ctx).Model(poa).Association("Attachment").Find(&poa.Attachment)
	if err != nil {
		return err
	}

	return nil

}

func CreatePurchaseOrderAttachment(ctx context.Context, purchaseOrderAttachment *PurchaseOrderAttachment) (*PurchaseOrderAttachment, error) {

	err := util.DBFromContext(ctx).Create(purchaseOrderAttachment).Error
	if err != nil {
		return nil, err
	}

	return purchaseOrderAttachment, nil

}

func GetPurchaseOrderAttachmentByID(ctx context.Context, id int) (*PurchaseOrderAttachment, error) {

	purchaseOrderAttachment := &PurchaseOrderAttachment{}

	err := util.DBFromContext(ctx).First(purchaseOrderAttachment, id).Error
	if err != nil {
		return nil, err
	}

	return purchaseOrderAttachment, nil

}
