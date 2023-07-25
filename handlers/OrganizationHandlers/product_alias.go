package OrganizationHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ProductAliasCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
		return
	}

	productAliasCreateReq := models.ProductAliasCreateRequest{}
	err = productAliasCreateReq.ParseAndValidateRequest(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productAliasCreateReq.ProductID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = product.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if product.Client.OrganizationID != user.Organization.ID {
		util.ErrResponse(w, ErrProductDoesNotBelongToOrganization, http.StatusBadRequest)
		return
	}

	if err = models.EnsureProductAliasBarcodeUnique(ctx, product.Client.ID, productAliasCreateReq.Barcode); err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	productAlias := &models.ProductAlias{
		ProductID: productAliasCreateReq.ProductID,
		Barcode:   productAliasCreateReq.Barcode,
		Quantity:  productAliasCreateReq.Quantity,
	}

	if err = productAlias.Create(ctx); err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusCreated)
}

func UpdateProductAliasByBarcode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
		return
	}

	barcode, err := util.GetStringFromPath(r, "barcode")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	request := &models.ProductAliasUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	productAlias, err := models.GetProductAliasByBarcode(ctx, barcode)
	if err != nil {
		util.ErrResponse(w, err, http.StatusNotFound)
		return
	}

	product, err := models.GetProductByID(ctx, productAlias.ProductID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = product.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if product.Client.OrganizationID != user.Organization.ID {
		util.ErrResponse(w, ErrProductDoesNotBelongToOrganization, http.StatusBadRequest)
		return
	}

	if request.ProductID != 0 {
		productNew, err := models.GetProductByID(ctx, request.ProductID)
		if err != nil {
			util.ErrResponse(w, fmt.Errorf("%w (updated product)", err), http.StatusBadRequest)
			return
		}

		err = productNew.GetClient(ctx)
		if err != nil {
			util.ErrResponse(w, fmt.Errorf("%w (updated product's client)", err), http.StatusBadRequest)
			return
		}

		if productNew.Client.OrganizationID != user.Organization.ID {
			util.ErrResponse(w, fmt.Errorf("new %w", ErrProductDoesNotBelongToOrganization), http.StatusBadRequest)
			return
		}
	}

	if err := productAlias.UpdateWithRequest(ctx, request); err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func GetProductAliasGetByBarcode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
		return
	}

	barcode, err := util.GetStringFromPath(r, "barcode")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	productAlias, err := models.GetProductAliasByBarcode(ctx, barcode)
	if err != nil {
		util.ErrResponse(w, err, http.StatusNotFound)
		return
	}

	product, err := models.GetProductByID(ctx, productAlias.ProductID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = product.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if product.Client.OrganizationID != user.Organization.ID {
		util.ErrResponse(w, ErrProductDoesNotBelongToOrganization, http.StatusBadRequest)
		return
	}

	response := responses.GenerateGetProductAliasGetByBarcodeResponse(*productAlias)
	util.JSONResponse(w, response, http.StatusOK)
}

func DeleteProductAliasByBarcode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
		return
	}

	barcode, err := util.GetStringFromPath(r, "barcode")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	productAlias, err := models.GetProductAliasByBarcode(ctx, barcode)
	if err != nil {
		util.ErrResponse(w, err, http.StatusNotFound)
		return
	}

	product, err := models.GetProductByID(ctx, productAlias.ProductID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = product.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if product.Client.OrganizationID != user.Organization.ID {
		util.ErrResponse(w, ErrProductDoesNotBelongToOrganization, http.StatusBadRequest)
		return
	}

	if err := productAlias.Delete(ctx); err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}
