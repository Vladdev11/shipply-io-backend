package OrganizationHandlers

import (
	"fmt"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/api"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/tasks"
	"github.com/shipply-io/shipply-io-backend/util"
)

func PurchaseOrderItemCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to find organization", http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, purchaseOrder.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	request := models.PurchaseOrderItemCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, request.ProductID)
	if err != nil {
		util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
		return
	}

	purchaseOrderItem := &models.PurchaseOrderItem{
		PurchaseOrderID: purchaseOrder.ID,
		ProductID:       product.ID,
		Product:         product,
	}

	if product.ClientID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "product does not belong to this client", http.StatusBadRequest)
		return
	}

	err = purchaseOrderItem.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create purchase order item", http.StatusBadRequest)
		return
	}

	purchaseOrderItemJson := purchaseOrderItem.ConvertToReturnJSON(ctx)
	util.JSONResponse(w, purchaseOrderItemJson, http.StatusCreated)
}

func PurchaseOrderItemGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to find organization", http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, purchaseOrder.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	purchaseOrderItemID, err := util.GetIntFromPath(r, "item_id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order item id", http.StatusBadRequest)
		return
	}

	purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, purchaseOrderItemID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order item", http.StatusBadRequest)
		return
	}

	if purchaseOrderItem.PurchaseOrderID != purchaseOrder.ID {
		util.ErrorResponse(w, "purchase order item does not belong to this purchase order", http.StatusBadRequest)
		return
	}

	err = purchaseOrderItem.GetProduct(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
		return
	}

	purchaseOrderItemJson := purchaseOrderItem.ConvertToReturnJSON(ctx)
	util.JSONResponse(w, purchaseOrderItemJson, http.StatusOK)

}

func PurchaseOrderItemBulkUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to find organization", http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, purchaseOrder.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	purchaseOrderItemUpdateRequests, errs := models.ParseAndValidatePurchaseOrderItemUpdateRequests(r)
	if errs != nil {
		util.ErrorsResponse(w, errs, http.StatusBadRequest)
		return
	}

	itemsMap := make(map[int]struct {
		POItem  *models.PurchaseOrderItem
		Request models.PurchaseOrderItemUpdateRequest
	})

	for _, purchaseOrderItemUpdateRequest := range purchaseOrderItemUpdateRequests {

		request := purchaseOrderItemUpdateRequest

		purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, purchaseOrderItemUpdateRequest.ID)
		if err != nil {
			util.ErrorResponse(w, "failed to find purchase order item", http.StatusBadRequest)
			return
		}

		if purchaseOrderItem.PurchaseOrderID != purchaseOrder.ID {
			util.ErrorResponse(w, "purchase order item does not belong to this purchase order", http.StatusBadRequest)
			return
		}

		itemsMap[purchaseOrderItem.ID] = struct {
			POItem  *models.PurchaseOrderItem
			Request models.PurchaseOrderItemUpdateRequest
		}{
			POItem:  purchaseOrderItem,
			Request: request,
		}

		go tasks.AllocateInventoryByProduct(ctx, purchaseOrderItem.ProductID)
	}

	for _, item := range itemsMap {
		err := item.POItem.UpdateWithRequest(ctx, &item.Request)
		if err != nil {
			util.ErrorResponse(w, "failed to update purchase order item", http.StatusInternalServerError)
			return
		}
	}

	util.JSONResponse(w, util.JSONSuccess(), http.StatusOK)
}

func PurchaseOrderItemUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to find organization", http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, purchaseOrder.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	purchaseOrderItemID, err := util.GetIntFromPath(r, "item_id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order item id", http.StatusBadRequest)
		return
	}

	purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, purchaseOrderItemID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order item", http.StatusBadRequest)
		return
	}

	if purchaseOrderItem.PurchaseOrderID != purchaseOrder.ID {
		util.ErrorResponse(w, "purchase order item does not belong to this purchase order", http.StatusBadRequest)
		return
	}

	request := &models.PurchaseOrderItemUpdateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if len(errors) > 0 {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	err = purchaseOrderItem.UpdateWithRequest(ctx, request)
	if err != nil {
		util.ErrorResponse(w, "failed to update purchase order item", http.StatusInternalServerError)
		return
	}

	err = purchaseOrderItem.GetProduct(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
		return
	}

	go tasks.AllocateInventoryByProduct(ctx, purchaseOrderItem.ProductID)

	purchaseOrderItemJson := purchaseOrderItem.ConvertToReturnJSON(ctx)
	util.JSONResponse(w, purchaseOrderItemJson, http.StatusOK)
}

func PurchaseOrderItemDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to find organization", http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order ID", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, purchaseOrder.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	purchaseOrderItemID, err := util.GetIntFromPath(r, "item_id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order item id", http.StatusBadRequest)
		return
	}

	purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, purchaseOrderItemID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order item", http.StatusBadRequest)
		return
	}

	if purchaseOrderItem.PurchaseOrderID != purchaseOrder.ID {
		util.ErrorResponse(w, "purchase order item does not belong to this purchase order", http.StatusBadRequest)
		return
	}

	err = purchaseOrderItem.Delete(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to delete purchase order item", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func PurchaseOrderItemReceive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to find organization", http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "purchase_order_id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, purchaseOrder.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	batchReceiveRequests, errors := models.ValidateBatchPurchaseOrderItemReceiveRequest(r)
	if len(errors) > 0 {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	var purchaseOrderJSON []models.PurchaseOrderItemReturnJSON

	//loop over receive requests
	for _, batchReceiveRequest := range batchReceiveRequests {
		//get purchase order items
		purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, batchReceiveRequest.PurchaseOrderItemID)
		if err != nil {
			util.ErrorResponse(w, "failed to find purchase order item", http.StatusBadRequest)
			return
		}

		//get product
		product, err := models.GetProductByID(ctx, purchaseOrderItem.ProductID)
		if err != nil {
			util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
			return
		}

		if product.NeedsLotNumber && batchReceiveRequest.ProductLotID == 0 {
			util.ErrorResponse(w, "lot number is required for this product", http.StatusBadRequest)
			return
		}

		purchaseOrderItem.Received += batchReceiveRequest.Quantity
		if purchaseOrderItem.Received < 0 {
			util.ErrorResponse(w, "received quantity cannot be less than 0", http.StatusBadRequest)
			return
		}

		if batchReceiveRequest.Quantity > 0 {
			err = models.CreateInventory(ctx, purchaseOrderItem.ProductID, batchReceiveRequest.Quantity, batchReceiveRequest.LocationID, false, 0)
			if err != nil {
				util.ErrorResponse(w, "failed to create inventory", http.StatusBadRequest)
				return
			}

			err = product.UpdateProductInventoryLevels(ctx)
			if err != nil {
				util.ErrorResponse(w, "failed to update product inventory levels", http.StatusBadRequest)
				return
			}

			//create inventory audit log
			inventoryAuditLog := models.InventoryAuditLog{
				ProductID:  purchaseOrderItem.ProductID,
				LocationID: batchReceiveRequest.LocationID,
				Delta:      batchReceiveRequest.Quantity,
				Note:       "received inventory",
				ChangedBy:  user.ID,
			}

			err = inventoryAuditLog.Create(ctx)
			if err != nil {
				util.ErrorResponse(w, "failed to create inventory audit log", http.StatusBadRequest)
				return
			}

		} else if batchReceiveRequest.Quantity < 0 {

			//use negative batchReceiveRequest.Quantity to make postive, which is what RemoveInventory expects
			err = models.RemoveInventory(ctx, purchaseOrderItem.ProductID, -batchReceiveRequest.Quantity, batchReceiveRequest.LocationID, false)
			if err != nil {
				util.ErrorResponse(w, "failed to remove inventory", http.StatusBadRequest)
				return
			}

			err = product.UpdateProductInventoryLevels(ctx)
			if err != nil {
				util.ErrorResponse(w, "failed to update product inventory levels", http.StatusBadRequest)
				return
			}

			//create inventory audit log
			inventoryAuditLog := models.InventoryAuditLog{
				ProductID:  purchaseOrderItem.ProductID,
				LocationID: batchReceiveRequest.LocationID,
				Delta:      batchReceiveRequest.Quantity,
				Note:       "received inventory",
				ChangedBy:  user.ID,
			}

			err = inventoryAuditLog.Create(ctx)
			if err != nil {
				util.ErrorResponse(w, "failed to create inventory audit log", http.StatusBadRequest)
				return
			}
		}

		//update allocated inventory
		go tasks.AllocateInventoryByProduct(ctx, purchaseOrderItem.ProductID)

		err = purchaseOrderItem.Update(ctx)
		if err != nil {
			util.ErrorResponse(w, "failed to update purchase order item", http.StatusInternalServerError)
			return
		}

		err = purchaseOrderItem.GetProduct(ctx)
		if err != nil {
			util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
			return
		}
		purchaseOrderItem.ProductID = 0

		purchaseOrderJSON = append(purchaseOrderJSON, purchaseOrderItem.ConvertToReturnJSON(ctx))

	}

	util.JSONResponse(w, purchaseOrderJSON, http.StatusOK)

}

func PurchaseOrderItemReject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to find organization", http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "purchase_order_id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, purchaseOrder.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	purchaseOrderItemID, err := util.GetIntFromPath(r, "item_id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order item id", http.StatusBadRequest)
		return
	}

	purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, purchaseOrderItemID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order item", http.StatusBadRequest)
		return
	}

	request := models.PurchaseOrderItemRejectRequest{}
	errors := request.ParseAndValidateRequest(r)
	if len(errors) > 0 {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, purchaseOrderItem.ProductID)
	if err != nil {
		util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
		return
	}

	location, err := models.GetLocationByID(ctx, request.Data.LocationID)
	if err != nil {
		util.ErrorResponse(w, "failed to find location", http.StatusBadRequest)
		return
	}

	err = location.GetWarehouse(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to find warehouse", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsWarehouseOwner(ctx, location.WarehouseID) {
		util.ErrorResponse(w, "user does not have access to this warehouse", http.StatusForbidden)
		return
	}

	poir := models.PurchaseOrderItemRejection{
		PurchaseOrderItemID: purchaseOrderItemID,
		RejectedReaseon:     request.Data.RejectReason,
		Note:                request.Data.Notes,
		CreatedBy:           user.ID,
		Quantity:            request.Data.Quantity,
	}

	err = poir.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create purchase order item rejection", http.StatusInternalServerError)
		return
	}

	for _, image := range request.Images {

		uuid, err := util.GenerateUUID()
		if err != nil {
			util.ErrorResponse(w, "failed to generate uuid", http.StatusBadRequest)
			return
		}

		err = api.S3FromContext(ctx).UploadAttachment(image.ImageData, uuid, fmt.Sprintf(".%s", image.FileType))
		if err != nil {
			util.ErrorResponse(w, "failed to upload image", http.StatusBadRequest)
			return
		}

		attachment, err := models.CreateAttachment(ctx, &models.Attachment{
			Extension: image.FileType,
			UUID:      uuid,
			FileName:  image.FileName,
		})

		if err != nil {
			util.ErrorResponse(w, "failed to create attachment", http.StatusInternalServerError)
			return
		}

		poira := models.PurchaseOrderItemRejectionAttachment{
			PurchaseOrderItemRejectionID: poir.ID,
			AttachmentID:                 attachment.ID,
		}

		err = poira.Create(ctx)
		if err != nil {
			util.ErrorResponse(w, "failed to create inbound shipment item reject attachement", http.StatusInternalServerError)
			return
		}
	}

	err = models.CreateInventory(ctx, purchaseOrderItem.ProductID, request.Data.Quantity, request.Data.LocationID, true, poir.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to create inventory", http.StatusBadRequest)
		return
	}

	err = product.UpdateProductInventoryLevels(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to update product inventory levels", http.StatusBadRequest)
		return
	}

	//create inventory audit log
	inventoryAuditLog := models.InventoryAuditLog{
		ProductID:  purchaseOrderItem.ProductID,
		LocationID: request.Data.LocationID,
		Delta:      request.Data.Quantity,
		Note:       "received damaged inventory",
		ChangedBy:  user.ID,
	}

	err = inventoryAuditLog.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create inventory audit log", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func PurchaseOrderItemUpdateIPAInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to find organization", http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "purchase_order_id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, purchaseOrder.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	purchaseOrderItemID, err := util.GetIntFromPath(r, "item_id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order item id", http.StatusBadRequest)
		return
	}

	purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, purchaseOrderItemID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order item", http.StatusBadRequest)
		return
	}

	request := models.PurchaseOrderItemUpdateIPARequest{}
	errors := request.ParseAndValidateRequest(r)
	if len(errors) > 0 {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, purchaseOrderItem.ProductID)
	if err != nil {
		util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
		return
	}

	product.Length = request.Length
	product.Width = request.Width
	product.Height = request.Height
	product.Weight = request.Weight
	product.WeightUnit = "oz"

	err = product.UpdateIPA(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to update product ipa", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func PurchaseOrderItemScanInput(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to find organization", http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "purchase_order_id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, purchaseOrder.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	request := models.PurchaseOrderItemScanInputRequest{}
	errors := request.ParseAndValidateRequest(r)
	if len(errors) > 0 {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByBarcodeAndClientID(ctx, request.Value, purchaseOrder.ClientID)
	if err != nil {
		util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
		return
	}

	err = purchaseOrder.GetItems(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order items", http.StatusBadRequest)
		return
	}

	//TODO add in barcode aliases

	for _, item := range purchaseOrder.Items {
		if item.ProductID == product.ID {
			//return purchase order item scan input response
			response := models.PurchaseOrderItemScanInputResponse{
				PurchaseOrderItemID: item.ID,
				Quantity:            1,
			}

			util.JSONResponse(w, response, http.StatusOK)
			return
		}
	}

	util.ErrorResponse(w, "product not found in purchase order", http.StatusBadRequest)
}
