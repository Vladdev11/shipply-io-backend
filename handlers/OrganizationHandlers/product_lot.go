package OrganizationHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ProductLotListByProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "product_id")
	if err != nil {
		util.ErrorResponse(w, "failed to get shipping method id", http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productID)
	if err != nil {
		util.ErrorResponse(w, "failed to get product", http.StatusBadRequest)
		return
	}

	err = product.GetClient(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusBadRequest)
		return
	}

	if product.Client.OrganizationID != user.Organization.ID {
		util.ErrorResponse(w, "product does not belong to organization", http.StatusBadRequest)
		return
	}

	err = product.GetProductLots(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get product lots", http.StatusBadRequest)
		return
	}

	productLotReturnJSON := make([]models.ProductLotReturnJSON, 0)
	for _, productLot := range product.ProductLots {
		productLotReturnJSON = append(productLotReturnJSON, *productLot.ConvertToReturnJSON(ctx))
	}

	util.JSONResponse(w, productLotReturnJSON, http.StatusOK)

}

func ProductLotCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	var request models.ProductLotCreateRequest
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, request.ProductID)
	if err != nil {
		util.ErrorResponse(w, "failed to get product", http.StatusBadRequest)
		return
	}

	if !product.NeedsLotNumber {
		util.ErrorResponse(w, "product does not have lot tracking enabled", http.StatusBadRequest)
		return
	}

	err = product.GetClient(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusBadRequest)
		return
	}

	if product.Client.OrganizationID != user.Organization.ID {
		util.ErrorResponse(w, "product does not belong to organization", http.StatusBadRequest)
		return
	}

	productLot := models.ProductLot{
		CreatedBy:  user.ID,
		LotNumber:  request.LotNumber,
		ProductID:  request.ProductID,
		ExpiryDate: request.ExpiryDate.Time,
	}

	productLot.DetermineActive()

	err = productLot.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create product lot", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, productLot.ConvertToReturnJSON(ctx), http.StatusOK)

}
