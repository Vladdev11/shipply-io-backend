package OrganizationHandlers

import (
	"fmt"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"

	shipengineHandlers "github.com/shipply-io/shipply-io-backend/api/shipengine/handlers"
	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
)

func ListBoxes(w http.ResponseWriter, r *http.Request) {
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

	// TODO convert to typical search query
	boxes, err := user.Organization.GetBoxes(ctx)
	if err != nil && err != gorm.ErrRecordNotFound {
		util.ErrorResponse(w, "failed to get boxes", http.StatusBadRequest)
		return
	}

	response := responses.GenerateListBoxesResponse(boxes)
	util.JSONResponse(w, response, http.StatusOK)
}

func GetBox(w http.ResponseWriter, r *http.Request) {
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

	boxID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "box_id is required", http.StatusBadRequest)
		return
	}

	box, err := models.GetBoxByID(ctx, boxID)
	if err != nil {
		util.ErrorResponse(w, "failed to get box", http.StatusBadRequest)
		return
	}

	if box.OrganizationID != user.Organization.ID {
		util.ErrorResponse(w, "box does not belong to organization", http.StatusUnauthorized)
		return
	}

	response := responses.GenerateGetBoxResponse(*box)
	util.JSONResponse(w, response, http.StatusOK)
}

func CreateBox(w http.ResponseWriter, r *http.Request) {
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

	request := models.BoxCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	box := models.Box{
		Name:           request.Name,
		OrganizationID: user.Organization.ID,
		Barcode:        request.Barcode,
		Length:         request.Length,
		Width:          request.Width,
		Height:         request.Height,
		Weight:         request.Weight,
		Type:           request.Type,
		Active:         request.Active,
		Cost:           request.Cost,
		CreatedBy:      user.ID,
	}

	uuid, err := util.GenerateUUID()
	if err != nil {
		util.ErrorResponse(w, "failed to generate uuid", http.StatusBadRequest)
		return
	}

	shipenginePackage, err := shipengineHandlers.CreatePackage(ctx, ShipengineModels.Package{
		PackageCode: fmt.Sprintf("custom_%s", uuid),
		Name:        box.Name,
		Dimensions: &ShipengineModels.Dimensions{
			Length: box.Length,
			Width:  box.Width,
			Height: box.Height,
			Unit:   "inch",
		},
	})
	if err != nil {
		fmt.Println(err)
		util.ErrorResponse(w, "failed to create shipengine package", http.StatusBadRequest)
		return
	}

	box.ShipenginePackageID = shipenginePackage.PackageID
	box.ShipenginePackageCode = shipenginePackage.PackageCode

	err = box.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create box", http.StatusBadRequest)
		return
	}

	response := responses.GenerateCreateBoxResponse(box)
	util.JSONResponse(w, response, http.StatusOK)
}

func UpdateBox(w http.ResponseWriter, r *http.Request) {
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

	boxID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "box_id is required", http.StatusBadRequest)
		return
	}

	box, err := models.GetBoxByID(ctx, boxID)
	if err != nil {
		util.ErrorResponse(w, "failed to get box", http.StatusBadRequest)
		return
	}

	if box.OrganizationID != user.Organization.ID {
		util.ErrorResponse(w, "box does not belong to organization", http.StatusUnauthorized)
		return
	}

	request := models.BoxUpdateRequest{}
	request.ID = boxID
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	box.UpdateWithRequest(request)

	err = box.Update(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to update box", http.StatusBadRequest)
		return
	}

	// TODO update shipengine package

	response := responses.GenerateUpdateBoxResponse(*box)
	util.JSONResponse(w, response, http.StatusOK)
}

func DeleteBox(w http.ResponseWriter, r *http.Request) {
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

	boxID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "box_id is required", http.StatusBadRequest)
		return
	}

	box, err := models.GetBoxByID(ctx, boxID)
	if err != nil {
		util.ErrorResponse(w, "failed to get box", http.StatusBadRequest)
		return
	}

	if box.OrganizationID != user.Organization.ID {
		util.ErrorResponse(w, "box does not belong to organization", http.StatusUnauthorized)
		return
	}

	//TODO logic around when we prevent a box from being deleted

	err = box.Delete(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to delete box", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}
