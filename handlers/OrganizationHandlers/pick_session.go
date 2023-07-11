package OrganizationHandlers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

func CreatePickSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	_, err = models.GetCurrentPickSessionByUserID(ctx, user.ID)
	if err == nil {
		util.ErrorResponse(w, "user already has an active pick session", http.StatusBadRequest)
		return
	}

	request := models.PickSessionCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	pickSession := models.PickSession{
		UserID:      user.ID,
		WarehouseID: request.WarehouseID,
	}
	err = pickSession.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create pick session", http.StatusInternalServerError)
		return
	}

	ordersInPickSessionCount := 0
	for ordersInPickSessionCount < request.ToteCount {

		order, err := models.GetNextOrderReadyForPicking(ctx, request.WarehouseID)
		if err != nil {
			if err.Error() == "no orders ready for picking" {
				break
			}
			util.ErrorResponse(w, "failed to get next order ready for picking", http.StatusInternalServerError)
			return
		}

		pickSessionOrder := models.PickSessionOrder{
			PickSessionID: pickSession.ID,
			OrderID:       order.ID,
		}
		err = pickSessionOrder.Create(ctx)
		if err != nil {
			util.ErrorResponse(w, "failed to create pick session order", http.StatusInternalServerError)
			return
		}

		err = pickSessionOrder.CreatePickSessionOrderItems(ctx, request.WarehouseID)
		if err != nil {
			util.ErrorResponse(w, "failed to create pick session order items", http.StatusInternalServerError)
			return
		}

		ordersInPickSessionCount++
	}

	pickSessionResponse, err := pickSession.ConvertToPickSessionResponse(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to convert pick session to response", http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, pickSessionResponse, http.StatusCreated)

}

func GetActivePickSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	pickSession, err := models.GetCurrentPickSessionByUserID(ctx, user.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session", http.StatusBadRequest)
		return
	}

	pickSessionResponse, err := pickSession.ConvertToPickSessionResponse(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to convert pick session to response", http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, pickSessionResponse, http.StatusOK)

}

func PickSessionSelectItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	pickSession, err := models.GetCurrentPickSessionByUserID(ctx, user.ID)
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			util.ErrorResponse(w, "no active pick session", http.StatusBadRequest)
			return
		}

		util.ErrorResponse(w, "failed to get pick session", http.StatusBadRequest)
		return
	}

	request := models.PickSessionSelectItemRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	// This big chunk to check if the barcode scanned has a product in the pick session that has not been picked yet
	pickSessionOrderItems, err := models.GetPickSessionOrderItemsByPickSessionID(ctx, pickSession.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session order items", http.StatusInternalServerError)
		return
	}

	// Assume all products are fully picked, then remove the ones that are not and put them in the remainingProductIDs slice
	fullyPickedProductIDs := []int{}
	for _, pickSessionOrderItem := range pickSessionOrderItems {
		if !util.SliceContainsInt(fullyPickedProductIDs, pickSessionOrderItem.ProductID) {
			fullyPickedProductIDs = append(fullyPickedProductIDs, pickSessionOrderItem.ProductID)
		}
	}

	remainingProductIDs := []int{}
	for _, pickSessionOrderItem := range pickSessionOrderItems {
		if pickSessionOrderItem.QuantityPicked < pickSessionOrderItem.QuantityToPick {

			// Add the product to the remaining slice if it is not already in there
			if !util.SliceContainsInt(remainingProductIDs, pickSessionOrderItem.ProductID) {
				remainingProductIDs = append(remainingProductIDs, pickSessionOrderItem.ProductID)
			}

			// Remove the product from the fully picked slice if it is not fully picked
			for i, productID := range fullyPickedProductIDs {
				if productID == pickSessionOrderItem.ProductID {
					fullyPickedProductIDs = append(fullyPickedProductIDs[:i], fullyPickedProductIDs[i+1:]...)
					break
				}
			}
		}
	}

	var scannedProduct *models.Product
	for _, productID := range remainingProductIDs {

		product, err := models.GetProductByID(ctx, productID)
		if err != nil {
			util.ErrorResponse(w, "failed to get product", http.StatusInternalServerError)
			return
		}

		if product.Barcode == request.Barcode {
			scannedProduct = product
			break
		}

	}

	if scannedProduct == nil {
		for _, productID := range fullyPickedProductIDs {

			product, err := models.GetProductByID(ctx, productID)
			if err != nil {
				util.ErrorResponse(w, "failed to get product", http.StatusInternalServerError)
				return
			}

			if product.Barcode == request.Barcode {
				util.ErrorResponse(w, "product already fully picked", http.StatusBadRequest)
				return
			}

		}

		util.ErrorResponse(w, "product not found in pick session", http.StatusBadRequest)
		return
	}

	var pickSessionOrderItemsWithScannedProductID []models.PickSessionOrderItem
	for _, pickSessionOrderItem := range pickSessionOrderItems {
		if pickSessionOrderItem.ProductID == scannedProduct.ID {
			err = pickSessionOrderItem.GetPickSessionOrder(ctx)
			if err != nil {
				util.ErrorResponse(w, "failed to get pick session order", http.StatusInternalServerError)
				return
			}

			pickSessionOrderItemsWithScannedProductID = append(pickSessionOrderItemsWithScannedProductID, pickSessionOrderItem)
		}
	}

	// sort pick session order items by whether they have a tote assigned (unassigned totes first)
	sort.Slice(pickSessionOrderItems, func(i, j int) bool {

		if pickSessionOrderItems[i].PickSessionOrder.LocationID == nil && pickSessionOrderItems[j].PickSessionOrder.LocationID != nil {
			return true
		}

		if pickSessionOrderItems[i].PickSessionOrder.LocationID != nil && pickSessionOrderItems[j].PickSessionOrder.LocationID == nil {
			return false
		}

		return false
	})

	// Get the next pick session order item to pick (items with unassigned totes will come first)
	nextPickSessionOrderItemToPick := pickSessionOrderItemsWithScannedProductID[0]

	if nextPickSessionOrderItemToPick.PickSessionOrder.LocationID == nil {
		util.JSONResponse(w, models.PickSessionSelectItemResponse{
			PickSessionOrderID:      nextPickSessionOrderItemToPick.PickSessionOrder.ID,
			QuantityRemainingToPick: nextPickSessionOrderItemToPick.QuantityToPick - nextPickSessionOrderItemToPick.QuantityPicked,
		}, http.StatusOK)
		return
	}

	tote, err := models.GetLocationByID(ctx, *nextPickSessionOrderItemToPick.PickSessionOrder.LocationID)
	if err != nil {
		util.ErrorResponse(w, "failed to get next picking tote", http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, models.PickSessionSelectItemResponse{
		ToteName:                tote.Name,
		PickSessionOrderID:      nextPickSessionOrderItemToPick.PickSessionOrder.ID,
		QuantityRemainingToPick: nextPickSessionOrderItemToPick.QuantityToPick - nextPickSessionOrderItemToPick.QuantityPicked,
	}, http.StatusOK)

}

func PickSessionAssignTote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	pickSession, err := models.GetCurrentPickSessionByUserID(ctx, user.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session", http.StatusBadRequest)
		return
	}

	request := models.PickSessionAssignToteRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	toteLocationID, err := strconv.Atoi(request.ToteBarcode)
	if err != nil {
		util.ErrorResponse(w, "invalid tote barcode", http.StatusBadRequest)
		return
	}

	tote, err := models.GetLocationByID(ctx, toteLocationID)
	if err != nil {
		util.ErrorResponse(w, "failed to get tote", http.StatusBadRequest)
		return
	}

	if !tote.IsTote {
		util.ErrorResponse(w, "location is not a tote", http.StatusBadRequest)
		return
	}

	if tote.WarehouseID != pickSession.WarehouseID {
		util.ErrorResponse(w, "tote not assigned to this warehouse", http.StatusBadRequest)
		return
	}

	if tote.HasActivePickSessionOrder(ctx) {
		util.ErrorResponse(w, "tote already assigned to a pick session order", http.StatusBadRequest)
		return
	}

	pickSessionOrder, err := models.GetPickSessionOrderByID(ctx, request.PickSessionOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session order", http.StatusBadRequest)
		return
	}

	if pickSessionOrder.PickSessionID != pickSession.ID {
		util.ErrorResponse(w, "order not assigned to this pick session", http.StatusBadRequest)
		return
	}

	if pickSessionOrder.LocationID != nil {
		util.ErrorResponse(w, "order already assigned to a tote", http.StatusBadRequest)
		return
	}

	// Assign Tote to Order
	pickSessionOrder.LocationID = &tote.ID
	err = pickSessionOrder.Update(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to update pick session order", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func PickSessionConfirmTote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	pickSession, err := models.GetCurrentPickSessionByUserID(ctx, user.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session", http.StatusBadRequest)
		return
	}

	request := models.PickSessionConfirmToteRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	pickSessionOrder, err := models.GetPickSessionOrderByID(ctx, request.PickSessionOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session order", http.StatusBadRequest)
		return
	}

	if pickSessionOrder.PickSessionID != pickSession.ID {
		util.ErrorResponse(w, "order not assigned to this pick session", http.StatusBadRequest)
		return
	}

	if pickSessionOrder.LocationID == nil {
		util.ErrorResponse(w, "order not assigned to a tote", http.StatusBadRequest)
		return
	}

	toteLocationID, err := strconv.Atoi(request.ToteBarcode)
	if err != nil {
		util.ErrorResponse(w, "invalid tote barcode", http.StatusBadRequest)
		return
	}

	tote, err := models.GetLocationByID(ctx, toteLocationID)
	if err != nil {
		util.ErrorResponse(w, "failed to get tote", http.StatusBadRequest)
		return
	}

	if !tote.IsTote {
		util.ErrorResponse(w, "location is not a tote", http.StatusBadRequest)
		return
	}

	if tote.WarehouseID != pickSession.WarehouseID {
		util.ErrorResponse(w, "tote not assigned to this warehouse", http.StatusBadRequest)
		return
	}

	if *pickSessionOrder.LocationID != tote.ID {
		util.ErrorResponse(w, "tote not assigned to this order", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func PickSessionPick(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	pickSession, err := models.GetCurrentPickSessionByUserID(ctx, user.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session", http.StatusBadRequest)
		return
	}

	request := models.PickSessionPickRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	toteLocationID, err := strconv.Atoi(request.ToteBarcode)
	if err != nil {
		util.ErrorResponse(w, "invalid tote barcode", http.StatusBadRequest)
		return
	}

	tote, err := models.GetLocationByID(ctx, toteLocationID)
	if err != nil {
		util.ErrorResponse(w, "failed to get tote", http.StatusBadRequest)
		return
	}

	if tote.WarehouseID != pickSession.WarehouseID {
		util.ErrorResponse(w, "tote not assigned to this warehouse", http.StatusBadRequest)
		return
	}

	if !tote.HasActivePickSessionOrder(ctx) {
		util.ErrorResponse(w, "tote not assigned to a pick session order", http.StatusBadRequest)
		return
	}

	pickSessionOrder, err := tote.GetActivePickSessionOrder(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session order", http.StatusBadRequest)
		return
	}

	if pickSessionOrder.PickSessionID != pickSession.ID {
		util.ErrorResponse(w, "tote is assigned to a different pick session", http.StatusBadRequest)
		return
	}

	err = pickSessionOrder.GetPickSessionOrderItems(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session order items", http.StatusBadRequest)
		return
	}

	var pickSessionOrderItem *models.PickSessionOrderItem
	for _, item := range pickSessionOrder.PickSessionOrderItems {

		product, err := models.GetProductByID(ctx, item.ProductID)
		if err != nil {
			util.ErrorResponse(w, "failed to get product", http.StatusBadRequest)
			return
		}

		if product.Barcode == request.ProductBarcode {
			pickSessionOrderItem = &item
			break
		}

	}

	if pickSessionOrderItem == nil {
		util.ErrorResponse(w, "product not found in order", http.StatusBadRequest)
		return
	}

	if pickSessionOrderItem.QuantityPicked >= pickSessionOrderItem.QuantityToPick {
		util.ErrorResponse(w, "product already picked for this order", http.StatusBadRequest)
		return
	}

	// picks quantity 1
	err = pickSessionOrderItem.Pick(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to pick product", http.StatusInternalServerError)
		return
	}

	// TODO move inventory to tote location

	// refresh pickSessionOrder.PickSessionOrderItems
	err = pickSessionOrder.GetPickSessionOrderItems(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session order items", http.StatusInternalServerError)
		return
	}

	// Determine if pick_session_order is fully picked, and if so, set in database
	isPickSessionOrderFullyPicked := true
	for _, item := range pickSessionOrder.PickSessionOrderItems {
		if item.QuantityPicked < item.QuantityToPick {
			isPickSessionOrderFullyPicked = false
			break
		}
	}

	if isPickSessionOrderFullyPicked {
		pickSessionOrder.Picked = true
		err = pickSessionOrder.Update(ctx)
		if err != nil {
			util.ErrorResponse(w, "failed to update pick session order", http.StatusInternalServerError)
			return
		}
	}

	// New tote for next pick (pick session order item is fully picked)
	if pickSessionOrderItem.QuantityPicked >= pickSessionOrderItem.QuantityToPick {

		pickSessionOrderItems, err := models.GetPickSessionOrderItemsByPickSessionID(ctx, pickSession.ID)
		if err != nil {
			util.ErrorResponse(w, "failed to get pick session order items", http.StatusInternalServerError)
			return
		}

		var pickSessionOrderItemsWithSameProductIDAndQuantityRemaining []models.PickSessionOrderItem
		for _, item := range pickSessionOrderItems {
			if item.QuantityPicked < item.QuantityToPick {
				if pickSessionOrderItem.ProductID == item.ProductID {

					err = item.GetPickSessionOrder(ctx)
					if err != nil {
						util.ErrorResponse(w, "failed to get pick session order", http.StatusInternalServerError)
						return
					}
					pickSessionOrderItemsWithSameProductIDAndQuantityRemaining = append(pickSessionOrderItemsWithSameProductIDAndQuantityRemaining, item)

				}
			}
		}

		if len(pickSessionOrderItemsWithSameProductIDAndQuantityRemaining) == 0 {
			// No more items to pick for this product)
			util.JSONResponse(w, models.PickSessionPickResponse{IsProductPickingComplete: true}, http.StatusOK)
			return
		}

		// sort pick session order items by whether they have a tote assigned (unassigned totes first)
		sort.Slice(pickSessionOrderItemsWithSameProductIDAndQuantityRemaining, func(i, j int) bool {

			if pickSessionOrderItems[i].PickSessionOrder.LocationID == nil && pickSessionOrderItems[j].PickSessionOrder.LocationID != nil {
				return true
			}

			if pickSessionOrderItems[i].PickSessionOrder.LocationID != nil && pickSessionOrderItems[j].PickSessionOrder.LocationID == nil {
				return false
			}

			return false
		})

		// Get the next pick session order item to pick (items with unassigned totes will come first)
		nextPickSessionOrderItemToPick := pickSessionOrderItemsWithSameProductIDAndQuantityRemaining[0]

		response := models.PickSessionPickResponse{
			IsProductPickingComplete: false,
			QuantityRemainingToPick:  nextPickSessionOrderItemToPick.QuantityToPick - nextPickSessionOrderItemToPick.QuantityPicked,
			IsNewTote:                true,
			PickSessionOrderID:       nextPickSessionOrderItemToPick.PickSessionOrder.ID,
		}

		if nextPickSessionOrderItemToPick.PickSessionOrder.LocationID == nil {
			util.JSONResponse(w, response, http.StatusOK)
			return
		}

		tote, err := models.GetLocationByID(ctx, *nextPickSessionOrderItemToPick.PickSessionOrder.LocationID)
		if err != nil {
			util.ErrorResponse(w, "failed to get tote", http.StatusInternalServerError)
			return
		}

		response.ToteName = tote.Name

		util.JSONResponse(w, response, http.StatusOK)
		return
	}

	response := models.PickSessionPickResponse{
		IsProductPickingComplete: false,
		IsNewTote:                false,
		ToteName:                 tote.Name,
		QuantityRemainingToPick:  pickSessionOrderItem.QuantityToPick - pickSessionOrderItem.QuantityPicked,
		PickSessionOrderID:       pickSessionOrderItem.PickSessionOrder.ID,
	}

	util.JSONResponse(w, response, http.StatusOK)

}

func PickSessionComplete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	pickSession, err := models.GetCurrentPickSessionByUserID(ctx, user.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session", http.StatusBadRequest)
		return
	}

	// TODO maybe allow user to complete pick session even if not all items picked. This would require all orders assigned to a tote to be fully picked

	pickSessionOrderItems, err := models.GetPickSessionOrderItemsByPickSessionID(ctx, pickSession.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session order items", http.StatusInternalServerError)
		return
	}

	for _, pickSessionOrderItem := range pickSessionOrderItems {
		if pickSessionOrderItem.QuantityPicked < pickSessionOrderItem.QuantityToPick {
			util.ErrorResponse(w, "pick session not complete", http.StatusBadRequest)
			return
		}
	}

	err = pickSession.Complete(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to complete pick session", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}
