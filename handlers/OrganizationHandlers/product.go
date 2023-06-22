package OrganizationHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ProductSearch(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusBadRequest)
		return
	}

	productSearch := models.ProductSearchRequest{}
	err = productSearch.ParseAndValidateRequest(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if productSearch.ClientID != 0 {
		if !user.Organization.IsClientOwner(productSearch.ClientID) {
			util.ErrorResponse(w, "client does not belong to organization", http.StatusForbidden)
			return
		}
	}

	productSearch.OrganizationID = user.OwnerID
	products, err := user.Organization.SearchProducts(productSearch)
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
