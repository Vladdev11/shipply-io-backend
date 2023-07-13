package ClientHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

func UserSavedFilterCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	table, err := util.GetStringFromPath(r, "table")
	if err != nil {
		util.ErrorResponse(w, "failed to get table from path", http.StatusBadRequest)
		return
	}

	if _, ok := models.AllowedFilterTables[table]; ok {
		request := models.UserSavedFilterCreateRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			util.ErrorResponse(w, "failed to decode request", http.StatusBadRequest)
			return
		}
		usf, err := request.Execute(ctx, user.ID, table)
		if err != nil {
			util.ErrorResponse(w, "failed to create user saved filter", http.StatusBadRequest)
			return
		}
		usfj, err := usf.AsJSON()
		if err != nil {
			util.ErrorResponse(w, "failed to convert user saved filter to json", http.StatusBadRequest)
			return
		}
		util.JSONResponse(w, usfj, http.StatusOK)
	} else {
		util.ErrorResponse(w, "invalid table", http.StatusBadRequest)
	}
}

func UserSavedFilterList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	table, err := util.GetStringFromPath(r, "table")
	if err != nil {
		fmt.Println(mux.Vars(r))
		util.ErrorResponse(w, "failed to get table from path", http.StatusBadRequest)
		return
	}

	if _, ok := models.AllowedFilterTables[table]; ok {
		filters, err := models.GetUserSavedFiltersForTable(ctx, user.ID, table)
		if err != nil {
			util.ErrorResponse(w, "failed to get user filters", http.StatusBadRequest)
			return
		}
		util.JSONResponse(w, filters, http.StatusOK)
		return
	} else {
		util.ErrorResponse(w, "invalid table", http.StatusBadRequest)
	}
}

func UserSavedFilterUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	table, err := util.GetStringFromPath(r, "table")
	if err != nil {
		util.ErrorResponse(w, "failed to get table from path", http.StatusBadRequest)
		return
	}

	filterID, err := util.GetIntFromPath(r, "filter_id")
	if err != nil {
		util.ErrorResponse(w, "invalid filter id", http.StatusBadRequest)
		return
	}

	if _, ok := models.AllowedFilterTables[table]; ok {
		request := models.UserSavedFilterUpdateRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			util.ErrorResponse(w, "failed to decode request", http.StatusBadRequest)
			return
		}
		usf, err := request.Execute(ctx, user.ID, table, filterID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				util.ErrorResponse(w, "filter not found", http.StatusNotFound)
				return
			}
			util.ErrorResponse(w, "failed to update user saved filter", http.StatusBadRequest)
			return
		}
		usfj, err := usf.AsJSON()
		if err != nil {
			util.ErrorResponse(w, "failed to convert user saved filter to json", http.StatusBadRequest)
			return
		}
		util.JSONResponse(w, usfj, http.StatusOK)
		return
	} else {
		util.ErrorResponse(w, "invalid table", http.StatusBadRequest)
	}
}

func UserSavedFilterDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	table, err := util.GetStringFromPath(r, "table")
	if err != nil {
		util.ErrorResponse(w, "failed to get table from path", http.StatusBadRequest)
		return
	}

	filterID, err := util.GetIntFromPath(r, "filter_id")
	if err != nil {
		util.ErrorResponse(w, "invalid filter id", http.StatusBadRequest)
		return
	}

	if _, ok := models.AllowedFilterTables[table]; ok {
		if err := models.DeleteUserSavedFilter(ctx, user.ID, table, filterID); err != nil {
			if err == gorm.ErrRecordNotFound {
				util.ErrorResponse(w, "filter not found", http.StatusNotFound)
				return
			}
			util.ErrorResponse(w, "failed to delete user filter", http.StatusBadRequest)
			return
		}
		util.SuccessResponse(w, http.StatusOK)
		return
	} else {
		util.ErrorResponse(w, "invalid table", http.StatusBadRequest)
	}
}
