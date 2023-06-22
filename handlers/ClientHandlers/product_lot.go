package ClientHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ProductLotListByProduct(w http.ResponseWriter, r *http.Request) {

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

	productID, err := util.GetIntFromPath(r, "product_id")
	if err != nil {
		util.ErrorResponse(w, "failed to get shipping method id", http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(productID)
	if err != nil {
		util.ErrorResponse(w, "failed to get product", http.StatusBadRequest)
		return
	}

	if product.ClientID != user.Client.ID {
		util.ErrorResponse(w, "product does not belong to client", http.StatusBadRequest)
		return
	}

	err = product.GetProductLots()
	if err != nil {
		util.ErrorResponse(w, "failed to get product lots", http.StatusBadRequest)
		return
	}

	productLotReturnJSON := make([]models.ProductLotReturnJSON, 0)
	for _, productLot := range product.ProductLots {
		productLotReturnJSON = append(productLotReturnJSON, *productLot.ConvertToReturnJSON())
	}

	util.JSONResponse(w, productLotReturnJSON, http.StatusOK)

}
