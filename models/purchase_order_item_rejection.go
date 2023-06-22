package models

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PurchaseOrderItemRejection struct {
	ID                  int
	PurchaseOrderItemID int
	RejectedReaseon     string
	Note                string
	CreatedBy           int
	Quantity            int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	CreatedByUser     *User `gorm:"foreignKey:CreatedBy"`
	PurchaseOrderItem *PurchaseOrderItem

	Attachments []PurchaseOrderItemRejectionAttachment `gorm:"-"`
}

func (poir *PurchaseOrderItemRejection) Create() error {
	result := PGDB.Create(&poir)
	return result.Error
}

type PurchaseOrderItemRejectRequestData struct {
	Quantity     int    `json:"quantity"`
	LocationID   int    `json:"location_id"`
	RejectReason string `json:"reject_reason"`
	Notes        string `json:"notes"`
}

type PurchaseOrderItemRejectImage struct {
	ImageData multipart.File `json:"image_data"`
	FileType  string         `json:"file_type"`
	FileName  string         `json:"file_name"`
}

type PurchaseOrderItemRejectRequest struct {
	Data   PurchaseOrderItemRejectRequestData `json:"data"`
	File1  multipart.File                     `json:"file_1"`
	File2  multipart.File                     `json:"file_2"`
	File3  multipart.File                     `json:"file_3"`
	File4  multipart.File                     `json:"file_4"`
	File5  multipart.File                     `json:"file_5"`
	File6  multipart.File                     `json:"file_6"`
	File7  multipart.File                     `json:"file_7"`
	File8  multipart.File                     `json:"file_8"`
	File9  multipart.File                     `json:"file_9"`
	File10 multipart.File                     `json:"file_10"`
	Images []PurchaseOrderItemRejectImage
}

func (poir *PurchaseOrderItemRejectRequest) ParseAndValidateRequest(r *http.Request) []string {

	errors := []string{}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errors = append(errors, err.Error())
		return errors
	}

	var images []PurchaseOrderItemRejectImage
	for i := 1; i <= 10; i++ {
		imageKey := fmt.Sprintf("file%d", i)
		file, fileHeader, err := r.FormFile(imageKey)
		if err == nil {
			defer file.Close()
			if fileHeader.Header.Get("Content-Type") != "image/jpeg" && fileHeader.Header.Get("Content-Type") != "image/png" {
				errors = append(errors, "file must be a valid image")
				continue
			}

			images = append(images, PurchaseOrderItemRejectImage{
				ImageData: file,
				FileType:  strings.Split(fileHeader.Header.Get("Content-Type"), "/")[1],
				FileName:  fileHeader.Filename,
			})
		}
	}

	poir.Images = images

	// Retrieve the JSON data
	dataField := r.FormValue("data")

	// Parse the JSON data
	var request PurchaseOrderItemRejectRequest
	err := json.Unmarshal([]byte(dataField), &request.Data)
	if err != nil {
		return []string{err.Error()}
	}

	if request.Data.Quantity <= 0 {
		errors = append(errors, "quantity must be greater than 0")
	}
	poir.Data.Quantity = request.Data.Quantity

	if request.Data.LocationID <= 0 {
		errors = append(errors, "location_id must be greater than 0")
	}
	poir.Data.LocationID = request.Data.LocationID

	if request.Data.RejectReason == "" {
		errors = append(errors, "reject_reason is required")
	}
	poir.Data.RejectReason = request.Data.RejectReason

	poir.Data.Notes = request.Data.Notes

	if len(errors) > 0 {
		return errors
	}

	return nil

}

type PurchaseOrderItemRejectionReturnJSON struct {
	ID                  int            `json:"id"`
	PurchaseOrderItemID int            `json:"purchase_order_item_id"`
	RejectedReaseon     string         `json:"rejected_reason"`
	Note                string         `json:"note"`
	CreatedBy           UserReturnJSON `json:"created_by"`
	CreatedAt           time.Time      `json:"created_at"`
	Images              []string       `json:"images"`
	Quantity            int            `json:"quantity"`
}

func (poir *PurchaseOrderItemRejection) ConvertToReturnJSON() *PurchaseOrderItemRejectionReturnJSON {

	var createdByUser User
	var err error

	if poir.CreatedByUser == nil {
		createdByUser, err = GetUserByID(poir.CreatedBy)
		if err != nil {
			return nil
		}
	}

	attachments := []string{}
	if poir.Attachments != nil {
		for _, attachment := range poir.Attachments {
			returnJSON := attachment.ConvertToReturnJSON()
			attachments = append(attachments, returnJSON.URL)
		}
	}

	return &PurchaseOrderItemRejectionReturnJSON{
		ID:                  poir.ID,
		PurchaseOrderItemID: poir.PurchaseOrderItemID,
		RejectedReaseon:     poir.RejectedReaseon,
		Note:                poir.Note,
		CreatedBy:           *createdByUser.ConvertToReturnJSON(),
		CreatedAt:           poir.CreatedAt,
		Images:              attachments,
		Quantity:            poir.Quantity,
	}
}

func (poir *PurchaseOrderItemRejection) GetImages() error {

	var attachments []PurchaseOrderItemRejectionAttachment
	result := PGDB.Where("purchase_order_item_rejection_id = ?", poir.ID).Find(&attachments)
	if result.Error != nil {
		return result.Error
	}

	poir.Attachments = attachments

	return nil

}
