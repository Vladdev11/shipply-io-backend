package ClientHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ProductAliasCreate(w http.ResponseWriter, r *http.Request) {
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

	productAliasCreateReq := models.ProductAliasCreateRequest{}
	errors := productAliasCreateReq.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(productAliasCreateReq.ProductID)
	if err != nil {
		util.ErrorResponse(w, "failed to get product", http.StatusBadRequest)
		return
	}

	if product.ClientID != user.Client.ID {
		util.ErrorResponse(w, "product does not belong to client", http.StatusBadRequest)
		return
	}

	if !models.IsProductAliasBarcodeUnique(user.Client.ID, productAliasCreateReq.Barcode) {
		util.ErrorResponse(w, "barcode already exists", http.StatusBadRequest)
		return
	}

	productAlias := &models.ProductAlias{
		ProductID: productAliasCreateReq.ProductID,
		Barcode:   productAliasCreateReq.Barcode,
		Quantity:  productAliasCreateReq.Quantity,
	}

	if err = productAlias.Create(); err != nil {
		util.ErrorResponse(w, "failed to create product alias", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, productAlias, http.StatusOK)

}

func ProductAliasUpdateByBarcode(w http.ResponseWriter, r *http.Request) {
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

	productAlias, err := models.GetProductAliasByBarcode(barcode)
	if err != nil {
		util.ErrorResponse(w, "failed to find product alias", http.StatusNotFound)
		return
	}

	product, err := models.GetProductByID(productAlias.ProductID)
	if err != nil {
		util.ErrorResponse(w, "failed to get associated product", http.StatusBadRequest)
		return
	}

	if product.ClientID != user.Client.ID {
		util.ErrorResponse(w, "product does not belong to client", http.StatusBadRequest)
		return
	}

	if request.ProductID != 0 {
		productNew, err := models.GetProductByID(request.ProductID)
		if err != nil {
			util.ErrorResponse(w, "failed to get new associated product", http.StatusBadRequest)
			return
		}

		if productNew.ClientID != user.Client.ID {
			util.ErrorResponse(w, "product does not belong to client", http.StatusBadRequest)
			return
		}
	}

	if err := productAlias.UpdateWithRequest(request); err != nil {
		util.ErrorResponse(w, "failed to update product alias", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, productAlias, http.StatusOK)
}

func ProductAliasGetByBarcode(w http.ResponseWriter, r *http.Request) {
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

	barcode, err := util.GetStringFromPath(r, "barcode")
	if err != nil {
		util.ErrorResponse(w, "invalid barcode", http.StatusBadRequest)
		return
	}

	productAlias, err := models.GetProductAliasByBarcode(barcode)
	if err != nil {
		util.ErrorResponse(w, "failed to find product alias", http.StatusNotFound)
		return
	}

	product, err := models.GetProductByID(productAlias.ProductID)
	if err != nil {
		util.ErrorResponse(w, "failed to get associated product", http.StatusBadRequest)
		return
	}

	if product.ClientID != user.Client.ID {
		util.ErrorResponse(w, "product does not belong to client", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, productAlias, http.StatusOK)
}

func ProductAliasDeleteByBarcode(w http.ResponseWriter, r *http.Request) {
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

	barcode, err := util.GetStringFromPath(r, "barcode")
	if err != nil {
		util.ErrorResponse(w, "invalid barcode", http.StatusBadRequest)
		return
	}

	productAlias, err := models.GetProductAliasByBarcode(barcode)
	if err != nil {
		util.ErrorResponse(w, "failed to find product alias", http.StatusNotFound)
		return
	}

	product, err := models.GetProductByID(productAlias.ProductID)
	if err != nil {
		util.ErrorResponse(w, "failed to get associated product", http.StatusBadRequest)
		return
	}

	if product.ClientID != user.Client.ID {
		util.ErrorResponse(w, "product does not belong to client", http.StatusBadRequest)
		return
	}

	if err := productAlias.Delete(); err != nil {
		util.ErrorResponse(w, "failed to delete product alias", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
