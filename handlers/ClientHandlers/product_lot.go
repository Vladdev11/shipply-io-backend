package ClientHandlers

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

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "product_id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if product.ClientID != user.Client.ID {
		util.ErrResponse(w, ErrProductDoesNotBelongToClient, http.StatusBadRequest)
		return
	}

	err = product.GetProductLots(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	productLotReturnJSON := make([]models.ProductLotReturnJSON, 0)
	for _, productLot := range product.ProductLots {
		productLotReturnJSON = append(productLotReturnJSON, *productLot.ConvertToReturnJSON(ctx))
	}

	util.JSONResponse(w, productLotReturnJSON, http.StatusOK)

}
