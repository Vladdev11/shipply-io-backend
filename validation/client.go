package validation

import (
	"mime/multipart"
	"net/http"
)

/* ------------------------------ CreateClient ------------------------------ */

// CreateClientRequestData represents the formatted and validated data for the CreateClient endpoint
type CreateClientRequestData struct {
	Name string `json:"name"`
}

// ParseRequestToCreateClientRequestData parses the request body to CreateClientRequestData
func ParseRequestToCreateClientRequestData(r *http.Request) (*CreateClientRequestData, []string) {

	createClientRequestData := CreateClientRequestData{}
	errors := []string{}

	rawData, err := ParseJSONRequestBody(r)
	if err != nil {
		errors = append(errors, err.Error())
		return nil, errors
	}

	createClientRequestData.Name, err = validateRequiredStringField(rawData["name"], "name", 1, 255)
	if err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return nil, errors
	}

	return &createClientRequestData, nil
}

/* ------------------------------- UpdateClient ------------------------------- */

// UpdateClientRequestData represents the formatted and validated data for the UpdateClient endpoint
type UpdateClientRequestData struct {
	Name *string `json:"name"`
}

// ParseRequestToUpdateClientRequestData parses the request body to UpdateClientRequestData
func ParseRequestToUpdateClientRequestData(r *http.Request) (*UpdateClientRequestData, []string) {

	updateClientRequestData := UpdateClientRequestData{}
	errors := []string{}

	rawData, err := ParseJSONRequestBody(r)
	if err != nil {
		errors = append(errors, err.Error())
		return nil, errors
	}

	updateClientRequestData.Name, err = validateOptionalStringField(rawData["name"], "name", 1, 255)
	if err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return nil, errors
	}

	return &updateClientRequestData, nil
}

/* --------------------------- ClientUpdateAvatar --------------------------- */

// ClientUpdateAvatarRequestData represents the formatted and validated data for the ClientUpdateAvatar endpoint
type ClientUpdateAvatarRequestData struct {
	Image    multipart.File `json:"image"`
	FileType string         `json:"file_type"`
	FileName string         `json:"file_name"`
}

// ParseRequestToClientUpdateAvatarRequestData parses the request body to ClientUpdateAvatarRequestData
func ParseRequestToClientUpdateAvatarRequestData(r *http.Request) (*ClientUpdateAvatarRequestData, []string) {

	errors := []string{}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errors = append(errors, err.Error())
		return nil, errors
	}

	multipartImageData, err := ParseMultipartImage(r, "file")
	if err != nil && err != http.ErrMissingFile {
		errors = append(errors, err.Error())
		return nil, errors
	}

	if multipartImageData == nil {
		errors = append(errors, "file is required")
		return nil, errors
	}

	return &ClientUpdateAvatarRequestData{
		Image:    multipartImageData.FileData,
		FileType: multipartImageData.FileType,
		FileName: multipartImageData.FileName,
	}, nil
}
