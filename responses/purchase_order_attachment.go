package responses

import "github.com/shipply-io/shipply-io-backend/models"

/* ---------------------- ListPurchaseOrderAttachments ---------------------- */

// ListPurchaseOrderAttachmentsResponse represents the expected response body for the ListPurchaseOrderAttachments endpoint
type ListPurchaseOrderAttachmentsResponse struct {
	PurchaseOrderAttachments []PurchaseOrderAttachmentResponseForListPurchaseOrderAttachments `json:"purchase_order_attachments"`
}

// PurchaseOrderAttachmentResponseForListPurchaseOrderAttachments represents the expected response body for the PurchaseOrderAttachment used by the ListPurchaseOrderAttachments endpoint
type PurchaseOrderAttachmentResponseForListPurchaseOrderAttachments struct {
	ID       int    `json:"id"`
	FileName string `json:"file_name"`
	URL      string `json:"url"`
}

// GenerateListPurchaseOrderAttachmentsResponse generates the ListPurchaseOrderAttachmentsResponse from the provided PurchaseOrderAttachments
func GenerateListPurchaseOrderAttachmentsResponse(purchaseOrderAttachments []models.PurchaseOrderAttachment) ListPurchaseOrderAttachmentsResponse {

	var purchaseOrderAttachmentResponses []PurchaseOrderAttachmentResponseForListPurchaseOrderAttachments

	for _, purchaseOrderAttachment := range purchaseOrderAttachments {
		purchaseOrderAttachmentResponses = append(purchaseOrderAttachmentResponses, PurchaseOrderAttachmentResponseForListPurchaseOrderAttachments{
			ID:       purchaseOrderAttachment.ID,
			FileName: purchaseOrderAttachment.Attachment.FileName,
			URL:      purchaseOrderAttachment.URL,
		})
	}

	return ListPurchaseOrderAttachmentsResponse{
		PurchaseOrderAttachments: purchaseOrderAttachmentResponses,
	}

}
