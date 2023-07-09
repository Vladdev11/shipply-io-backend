package ClientHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ProductSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusBadRequest)
		return
	}

	productSearch := models.ProductSearchRequest{}
	err = productSearch.ParseAndValidateRequest(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	productSearch.ClientID = user.OwnerID

	products, err := user.Client.SearchProducts(ctx, productSearch)
	if err != nil {
		util.ErrorResponse(w, "failed to search products", http.StatusBadRequest)
		return
	}

	productJSON := []models.ProductReturnJSON{}
	for i := range products {
		productJSON = append(productJSON, *products[i].ConvertToReturnJSON())
	}

	util.JSONResponse(w, productJSON, http.StatusOK)

}
