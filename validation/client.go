package validation

import (
	"errors"
	"mime/multipart"
	"net/http"
)

/* -------------------------------- GetClient ------------------------------- */

// GetClientRequestData represents the formatted and validated data for the GetClient endpoint
type GetClientRequestData struct {
	ClientID int
}

// ParseRequestToGetClientRequestData parses the request body to GetClientRequestData
func ParseRequestToGetClientRequestData(r *http.Request) (*GetClientRequestData, error) {

	var err error
	var getClientRequestData GetClientRequestData

	getClientRequestData.ClientID, err = validateIntPathParameter(r, "client_id", 1)
	if err != nil {
		return nil, err
	}

	return &getClientRequestData, nil
}

/* ------------------------------ CreateClient ------------------------------ */

// CreateClientRequestData represents the formatted and validated data for the CreateClient endpoint
type CreateClientRequestData struct {
	Name string `json:"name"`
}

// ParseRequestToCreateClientRequestData parses the request body to CreateClientRequestData
func ParseRequestToCreateClientRequestData(r *http.Request) (*CreateClientRequestData, error) {

	var errs []error
	var createClientRequestData CreateClientRequestData

	rawData, err := ParseJSONRequestBody(r)
	if err != nil {
		return nil, err
	}

	createClientRequestData.Name, err = validateRequiredStringField(rawData["name"], "name", 1, 255)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &createClientRequestData, nil
}

/* ------------------------------- UpdateClient ------------------------------- */

// UpdateClientRequestData represents the formatted and validated data for the UpdateClient endpoint
type UpdateClientRequestData struct {
	ClientID int
	Name     *string `json:"name"`
}

// ParseRequestToUpdateClientRequestData parses the request body to UpdateClientRequestData
func ParseRequestToUpdateClientRequestData(r *http.Request) (*UpdateClientRequestData, error) {

	var errs []error
	var updateClientRequestData UpdateClientRequestData

	rawData, err := ParseJSONRequestBody(r)
	if err != nil {
		return nil, err
	}

	updateClientRequestData.ClientID, err = validateIntPathParameter(r, "client_id", 1)
	if err != nil {
		errs = append(errs, err)
	}

	updateClientRequestData.Name, err = validateOptionalStringField(rawData["name"], "name", 1, 255)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &updateClientRequestData, nil
}

/* --------------------------- ClientUpdateAvatar --------------------------- */

// ClientUpdateAvatarRequestData represents the formatted and validated data for the ClientUpdateAvatar endpoint
type ClientUpdateAvatarRequestData struct {
	ClientID int

	Image    multipart.File `json:"image"`
	FileType string         `json:"file_type"`
	FileName string         `json:"file_name"`
}

// ParseRequestToClientUpdateAvatarRequestData parses the request body to ClientUpdateAvatarRequestData
func ParseRequestToClientUpdateAvatarRequestData(r *http.Request) (*ClientUpdateAvatarRequestData, error) {

	var err error
	var errs []error
	var clientUpdateAvatarRequestData ClientUpdateAvatarRequestData

	clientUpdateAvatarRequestData.ClientID, err = validateIntPathParameter(r, "client_id", 1)
	if err != nil {
		errs = append(errs, err)
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return nil, err
	}

	multipartImageData, err := ValidateRequiredMultipartImage(r, "file")
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	clientUpdateAvatarRequestData.Image = multipartImageData.FileData
	clientUpdateAvatarRequestData.FileType = multipartImageData.FileType
	clientUpdateAvatarRequestData.FileName = multipartImageData.FileName

	return &clientUpdateAvatarRequestData, nil
}
