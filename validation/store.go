package validation

import (
	"errors"
	"net/http"
)

/* -------------------------------- GetStore -------------------------------- */

// GetStoreRequestData represents the formatted and validated data for the GetStore endpoint
type GetStoreRequestData struct {
	StoreID int
}

// ParseRequestToGetStoreRequestData parses the request body to GetStoreRequestData
func ParseRequestToGetStoreRequestData(r *http.Request) (*GetStoreRequestData, error) {

	var getStoreRequestData GetStoreRequestData
	var err error

	getStoreRequestData.StoreID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		return nil, err
	}

	return &getStoreRequestData, nil
}

/* ------------------------------ ActivateStore ----------------------------- */

// ActivateStoreRequestData represents the formatted and validated data for the ActivateStore endpoint
type ActivateStoreRequestData struct {
	StoreID int
}

// ParseRequestToActivateStoreRequestData parses the request body to ActivateStoreRequestData
func ParseRequestToActivateStoreRequestData(r *http.Request) (*ActivateStoreRequestData, error) {

	var activateStoreRequestData ActivateStoreRequestData
	var err error

	activateStoreRequestData.StoreID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		return nil, err
	}

	return &activateStoreRequestData, nil
}

/* ------------------------------ DeactivateStore ----------------------------- */

// DeactivateStoreRequestData represents the formatted and validated data for the DeactivateStore endpoint
type DeactivateStoreRequestData struct {
	StoreID int
}

// ParseRequestToDeactivateStoreRequestData parses the request body to DeactivateStoreRequestData
func ParseRequestToDeactivateStoreRequestData(r *http.Request) (*DeactivateStoreRequestData, error) {

	var deactivateStoreRequestData DeactivateStoreRequestData
	var err error

	deactivateStoreRequestData.StoreID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		return nil, err
	}

	return &deactivateStoreRequestData, nil
}

/* ------------------------------- UpdateStore ------------------------------ */

// UpdateStoreRequestData represents the formatted and validated data for the UpdateStore endpoint
type UpdateStoreRequestData struct {
	StoreID int
	Name    *string `json:"name"`
	// TODO -- add settings
}

// ParseRequestToUpdateStoreRequestData parses the request body to UpdateStoreRequestData
func ParseRequestToUpdateStoreRequestData(r *http.Request) (*UpdateStoreRequestData, error) {

	var updateStoreRequestData UpdateStoreRequestData
	var err error
	var errs []error

	updateStoreRequestData.StoreID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		errs = append(errs, err)
	}

	rawData, err := ParseJSONRequestBody(r)
	if err != nil {
		return nil, err
	}

	updateStoreRequestData.Name, err = validateOptionalStringField(rawData["name"], "name", 1, 255)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &updateStoreRequestData, nil
}

/* ------------------------------- DeleteStore ------------------------------ */

// DeleteStoreRequestData represents the formatted and validated data for the DeleteStore endpoint
type DeleteStoreRequestData struct {
	StoreID int
}

// ParseRequestToDeleteStoreRequestData parses the request body to DeleteStoreRequestData
func ParseRequestToDeleteStoreRequestData(r *http.Request) (*DeleteStoreRequestData, error) {

	var deleteStoreRequestData DeleteStoreRequestData
	var err error

	deleteStoreRequestData.StoreID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		return nil, err
	}

	return &deleteStoreRequestData, nil
}
