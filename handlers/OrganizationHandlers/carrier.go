package OrganizationHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ListCarriers(w http.ResponseWriter, r *http.Request) {
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

	carriers, err := models.GetCarriers(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get carriers", http.StatusInternalServerError)
		return
	}

	carriers = models.GetCarrierRequiredFields(carriers)

	response := responses.GenerateListCarriersResponse(carriers)
	util.JSONResponse(w, response, http.StatusOK)

}
