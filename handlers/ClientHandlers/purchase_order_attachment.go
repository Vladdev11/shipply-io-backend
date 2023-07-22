package ClientHandlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/shipply-io/shipply-io-backend/api"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ListPurchaseOrderAttachments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid purchase order ID. Must be an integer.", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase order", http.StatusInternalServerError)
		return
	}

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "user does not have access to this purchase order", http.StatusForbidden)
		return
	}

	err = purchaseOrder.GetAttachments(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get attachments", http.StatusInternalServerError)
		return
	}

	for i := range purchaseOrder.Attachments {
		purchaseOrder.Attachments[i].GetAttachment(ctx)
		purchaseOrder.Attachments[i].URL, err = api.S3FromContext(ctx).GetAttachmentURL(purchaseOrder.Attachments[i].Attachment.UUID, purchaseOrder.Attachments[i].Attachment.Extension, purchaseOrder.Attachments[i].Attachment.FileName)
		if err != nil {
			util.ErrorResponse(w, "failed to get attachment url", http.StatusInternalServerError)
			return
		}
	}

	response := responses.GenerateListPurchaseOrderAttachmentsResponse(purchaseOrder.Attachments)
	util.JSONResponse(w, response, http.StatusOK)

}

func PurchaseOrderAttachmentCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid purchase order ID. Must be an integer.", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase order", http.StatusBadRequest)
		return
	}

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "user does not have access to this purchase order", http.StatusForbidden)
		return
	}

	request := models.PurchaseOrderAttachmentCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if len(errors) > 0 {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	attachmentUUID := uuid.New()
	fileExtension := util.GetFileExtension(request.FileName)

	err = api.S3FromContext(ctx).UploadAttachment(request.File, attachmentUUID.String(), fileExtension)
	if err != nil {
		util.ErrorResponse(w, "failed to upload attachment to s3", http.StatusInternalServerError)
		return
	}

	attachment, err := models.CreateAttachment(r.Context(), &models.Attachment{
		FileName:  request.FileName,
		Extension: fileExtension,
		UUID:      attachmentUUID.String(),
	})
	if err != nil {
		util.ErrorResponse(w, "failed to create attachment", http.StatusInternalServerError)
		return
	}

	_, err = models.CreatePurchaseOrderAttachment(ctx, &models.PurchaseOrderAttachment{
		PurchaseOrderID: purchaseOrderID,
		AttachmentID:    attachment.ID,
	})
	if err != nil {
		util.ErrorResponse(w, "failed to create purchase order attachment", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusCreated)

}

func PurchaseOrderAttachmentDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid purchase order ID. Must be an integer.", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase order", http.StatusInternalServerError)
		return
	}

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "user does not have access to this purchase order", http.StatusForbidden)
		return
	}

	purchaseOrderAttachmentID, err := util.GetIntFromPath(r, "purchase_order_attachment_id")
	if err != nil {
		util.ErrorResponse(w, "Invalid purchase order attachment ID. Must be an integer.", http.StatusBadRequest)
		return
	}

	purchaseOrderAttachment, err := models.GetPurchaseOrderAttachmentByID(ctx, purchaseOrderAttachmentID)
	if err != nil {
		util.ErrorResponse(w, fmt.Sprintf("purchase order attachment with ID %d does not exist", purchaseOrderAttachmentID), http.StatusBadRequest)
		return
	}

	if purchaseOrderAttachment.PurchaseOrderID != purchaseOrderID {
		util.ErrorResponse(w, "purchase order attachment does not belong to purchase order", http.StatusBadRequest)
		return
	}

	err = purchaseOrderAttachment.Delete(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to delete purchase order attachment", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}
