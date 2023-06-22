package ClientHandlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/shipply-io/shipply-io-backend/api"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func PurchaseOrderList(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient()
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	request := models.PurchaseOrderListRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	if user.Client.ID != request.ClientID {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	purchaseOrders, count, total, err := user.Client.GetPurchaseOrders(request)
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase orders", http.StatusInternalServerError)
		return
	}

	searchResults, err := models.ConvertPurchaseOrdersToSearchResults(purchaseOrders, total, count)
	if err != nil {
		util.ErrorResponse(w, "failed to convert purchase orders to search results", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, searchResults, http.StatusOK)

}

func PurchaseOrderCreate(w http.ResponseWriter, r *http.Request) {

	_, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	request := models.PurchaseOrderCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	purchaseOrder := &models.PurchaseOrder{
		ClientID:       request.ClientID,
		PONumber:       request.PONumber,
		Status:         request.Status,
		ExpectedDate:   request.ExpectedDate.Time,
		ShipDate:       request.ShipDate.Time,
		ClosedDate:     request.ClosedDate.Time,
		VendorID:       request.VendorID,
		WarehouseID:    request.WarehouseID,
		TrackingNumber: request.TrackingNumber,
		TrackingURL:    request.TrackingURL,
	}

	err = purchaseOrder.Create()
	if err != nil {
		util.ErrorResponse(w, "failed to create purchase order", http.StatusInternalServerError)
		return
	}

	purchaseOrderJSON := purchaseOrder.ConvertToReturnJSON()
	util.JSONResponse(w, purchaseOrderJSON, http.StatusOK)

}

func PurchaseOrderGet(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	err = purchaseOrder.GetItems()
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase order items", http.StatusBadRequest)
		return
	}

	err = purchaseOrder.GetTags()
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase order tags", http.StatusBadRequest)
		return
	}

	err = purchaseOrder.GetHistory()
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase order notes", http.StatusBadRequest)
		return
	}

	purchaseOrderJSON := purchaseOrder.ConvertToReturnJSON()
	util.JSONResponse(w, purchaseOrderJSON, http.StatusOK)

}

func PurchaseOrderUpdate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient()
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	request := models.PurchaseOrderUpdateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	err = purchaseOrder.UpdateWithRequest(request)
	if err != nil {
		util.ErrorResponse(w, "failed to update purchase order", http.StatusInternalServerError)
		return
	}

	purchaseOrderJSON := purchaseOrder.ConvertToReturnJSON()
	util.JSONResponse(w, purchaseOrderJSON, http.StatusOK)

}

func PurchaseOrderDelete(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	err = purchaseOrder.Delete()
	if err != nil {
		util.ErrorResponse(w, "failed to delete purchase order", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func PurchaseOrderStatusCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	request := models.PurchaseOrderStatusCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	if user.OwnerID != request.ClientID {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	purchaseOrderStatus := &models.PurchaseOrderStatus{
		ClientID:    request.ClientID,
		Name:        request.Name,
		StatusColor: request.StatusColor,
		TextColor:   request.TextColor,
	}

	err = purchaseOrderStatus.Create()
	if err != nil {
		util.ErrorResponse(w, "failed to create purchase order status", http.StatusInternalServerError)
		return
	}

	purchaseOrderStatusJSON := purchaseOrderStatus.ConvertToReturnJSON()
	util.JSONResponse(w, purchaseOrderStatusJSON, http.StatusOK)

}

func PurchaseOrderStatusList(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient()
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	request := models.PurchaseOrderStatusListRequest{}
	request.ClientID = user.Client.ID
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	statuses, total, count, err := user.Client.GetPurchaseOrderStatuses(request)
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase order statuses", http.StatusBadRequest)
		return
	}

	searchResults, err := models.ConvertPurchaseOrderStatusesToSearchResults(statuses, total, count)
	if err != nil {
		util.ErrorResponse(w, "failed to convert purchase order statuses to search results", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, searchResults, http.StatusOK)
}

func PurchaseOrderStatusDelete(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	posID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order status id", http.StatusBadRequest)
		return
	}

	purchaseOrderStatus, err := models.GetPurchaseOrderStatusByID(posID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order status", http.StatusBadRequest)
		return
	}

	if user.OwnerID != purchaseOrderStatus.ClientID {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	if purchaseOrderStatus.IsInUse() {
		util.ErrorResponse(w, "purchase order status is in use and cannot be deleted", http.StatusBadRequest)
		return
	}

	err = purchaseOrderStatus.Delete()
	if err != nil {
		util.ErrorResponse(w, "failed to delete purchase order status", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func PurchaseOrderStatusUpdate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	posID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order status id", http.StatusBadRequest)
		return
	}

	purchaseOrderStatus, err := models.GetPurchaseOrderStatusByID(posID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order status", http.StatusBadRequest)
		return
	}

	if user.OwnerID != purchaseOrderStatus.ClientID {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	request := models.PurchaseOrderStatusUpdateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	err = purchaseOrderStatus.UpdateWithRequest(&request)
	if err != nil {
		util.ErrorResponse(w, "failed to update purchase order status", http.StatusInternalServerError)
		return
	}

	purchaseOrderStatusJSON := purchaseOrderStatus.ConvertToReturnJSON()
	util.JSONResponse(w, purchaseOrderStatusJSON, http.StatusOK)
}

func PurchaseOrderHistoryCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	request := &models.PurchaseOrderHistoryCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(request.PurchaseOrderId)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	request.CreatedBy = user.ID
	PurchaseOrderHistory := request.ConvertToPurchaseOrderHistory()
	err = PurchaseOrderHistory.Create()
	if err != nil {
		util.ErrorResponse(w, "failed to create purchase order note", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func PurchaseOrderAttachmentCreate(w http.ResponseWriter, r *http.Request) {

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

	purchaseOrder, err := models.GetPurchaseOrderByID(purchaseOrderID)
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

	err = api.UploadAttachmentToS3(request.File, attachmentUUID.String(), fileExtension)
	if err != nil {
		util.ErrorResponse(w, "failed to upload attachment to s3", http.StatusInternalServerError)
		return
	}

	attachment, err := models.CreateAttachment(&models.Attachment{
		FileName:  request.FileName,
		Extension: fileExtension,
		UUID:      attachmentUUID.String(),
	})
	if err != nil {
		util.ErrorResponse(w, "failed to create attachment", http.StatusInternalServerError)
		return
	}

	_, err = models.CreatePurchaseOrderAttachment(&models.PurchaseOrderAttachment{
		PurchaseOrderID: purchaseOrderID,
		AttachmentID:    attachment.ID,
	})
	if err != nil {
		util.ErrorResponse(w, "failed to create purchase order attachment", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusCreated)

}

func PurchaseOrderAttachmentList(w http.ResponseWriter, r *http.Request) {

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

	purchaseOrder, err := models.GetPurchaseOrderByID(purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase order", http.StatusInternalServerError)
		return
	}

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "user does not have access to this purchase order", http.StatusForbidden)
		return
	}

	err = purchaseOrder.GetAttachments()
	if err != nil {
		util.ErrorResponse(w, "failed to get attachments", http.StatusInternalServerError)
		return
	}

	// This loops through the purchase order attachments and converts them to the return JSON
	purchaseOrderAttachmentsReturnJSON := []models.PurchaseOrderAttachmentReturnJSON{}
	for _, purchaseOrderAttachment := range purchaseOrder.Attachments {
		purchaseOrderAttachmentReturnJSON, err := purchaseOrderAttachment.ConvertToReturnJSON()
		if err != nil {
			models.CreateSystemError(fmt.Sprintf("failed to convert purchase order attachment to return JSON: %s", err.Error()))
			continue
		}

		purchaseOrderAttachmentsReturnJSON = append(purchaseOrderAttachmentsReturnJSON, purchaseOrderAttachmentReturnJSON)
	}

	util.JSONResponse(w, purchaseOrderAttachmentsReturnJSON, http.StatusOK)

}

func PurchaseOrderAttachmentDelete(w http.ResponseWriter, r *http.Request) {

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

	purchaseOrder, err := models.GetPurchaseOrderByID(purchaseOrderID)
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

	purchaseOrderAttachment, err := models.GetPurchaseOrderAttachmentByID(purchaseOrderAttachmentID)
	if err != nil {
		util.ErrorResponse(w, fmt.Sprintf("purchase order attachment with ID %d does not exist", purchaseOrderAttachmentID), http.StatusBadRequest)
		return
	}

	if purchaseOrderAttachment.PurchaseOrderID != purchaseOrderID {
		util.ErrorResponse(w, "purchase order attachment does not belong to purchase order", http.StatusBadRequest)
		return
	}

	err = purchaseOrderAttachment.Delete()
	if err != nil {
		util.ErrorResponse(w, "failed to delete purchase order attachment", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}
