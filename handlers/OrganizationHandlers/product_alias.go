package OrganizationHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ProductAliasCreate(w http.ResponseWriter, r *http.Request) {
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

	productAliasCreateReq := models.ProductAliasCreateRequest{}
	errors := productAliasCreateReq.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productAliasCreateReq.ProductID)
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

	if !models.IsProductAliasBarcodeUnique(ctx, product.Client.ID, productAliasCreateReq.Barcode) {
		util.ErrorResponse(w, "barcode already exists", http.StatusBadRequest)
		return
	}

	productAlias := &models.ProductAlias{
		ProductID: productAliasCreateReq.ProductID,
		Barcode:   productAliasCreateReq.Barcode,
		Quantity:  productAliasCreateReq.Quantity,
	}

	if err = productAlias.Create(ctx); err != nil {
		util.ErrorResponse(w, "failed to create product alias", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, productAlias, http.StatusOK)
}

func ProductAliasUpdateByBarcode(w http.ResponseWriter, r *http.Request) {
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

	barcode, err := util.GetStringFromPath(r, "barcode")
	if err != nil {
		util.ErrorResponse(w, "invalid barcode", http.StatusBadRequest)
		return
	}

	request := &models.ProductAliasUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		util.ErrorResponse(w, "invalid json", http.StatusBadRequest)
		return
	}

	productAlias, err := models.GetProductAliasByBarcode(ctx, barcode)
	if err != nil {
		util.ErrorResponse(w, "failed to find product alias", http.StatusNotFound)
		return
	}

	product, err := models.GetProductByID(ctx, productAlias.ProductID)
	if err != nil {
		util.ErrorResponse(w, "failed to get associated product", http.StatusBadRequest)
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

	if request.ProductID != 0 {
		productNew, err := models.GetProductByID(ctx, request.ProductID)
		if err != nil {
			util.ErrorResponse(w, "failed to get new associated product", http.StatusBadRequest)
			return
		}

		err = productNew.GetClient(ctx)
		if err != nil {
			util.ErrorResponse(w, "failed to get client", http.StatusBadRequest)
			return
		}

		if productNew.Client.OrganizationID != user.Organization.ID {
			util.ErrorResponse(w, "new associated product does not belong to organization", http.StatusBadRequest)
			return
		}
	}

	if err := productAlias.UpdateWithRequest(ctx, request); err != nil {
		util.ErrorResponse(w, "failed to update product alias", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, productAlias, http.StatusOK)
}

func ProductAliasGetByBarcode(w http.ResponseWriter, r *http.Request) {
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

	barcode, err := util.GetStringFromPath(r, "barcode")
	if err != nil {
		util.ErrorResponse(w, "invalid barcode", http.StatusBadRequest)
		return
	}

	productAlias, err := models.GetProductAliasByBarcode(ctx, barcode)
	if err != nil {
		util.ErrorResponse(w, "failed to find product alias", http.StatusNotFound)
		return
	}

	product, err := models.GetProductByID(ctx, productAlias.ProductID)
	if err != nil {
		util.ErrorResponse(w, "failed to get associated product", http.StatusBadRequest)
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

	util.JSONResponse(w, productAlias, http.StatusOK)
}

func ProductAliasDeleteByBarcode(w http.ResponseWriter, r *http.Request) {
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

	barcode, err := util.GetStringFromPath(r, "barcode")
	if err != nil {
		util.ErrorResponse(w, "invalid barcode", http.StatusBadRequest)
		return
	}

	productAlias, err := models.GetProductAliasByBarcode(ctx, barcode)
	if err != nil {
		util.ErrorResponse(w, "failed to find product alias", http.StatusNotFound)
		return
	}

	product, err := models.GetProductByID(ctx, productAlias.ProductID)
	if err != nil {
		util.ErrorResponse(w, "failed to get associated product", http.StatusBadRequest)
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

	if err := productAlias.Delete(ctx); err != nil {
		util.ErrorResponse(w, "failed to delete product alias", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
