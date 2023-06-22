package middlewares

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/api/shopify"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func VerifyShopifyWebhook(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !shopify.VerifyWebhookRequest(r) {
			util.ErrorResponse(w, "invalid request", http.StatusBadRequest)
			return
		}

		err := models.LogShopifyWebhookRequest(r)
		if err != nil {
			models.CreateSystemError(err.Error())
		}

		next.ServeHTTP(w, r)
	})
}
