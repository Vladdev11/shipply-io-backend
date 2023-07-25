package OrganizationHandlers

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/shipply-io/shipply-io-backend/api"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
	"github.com/shipply-io/shipply-io-backend/validation"
)

func ListClients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := models.UserFromContext(ctx)

	clients, err := models.GetClientsByOrganizationID(ctx, user.Organization.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to get clients", http.StatusUnauthorized)
		return
	}

	response := responses.GenerateListClientsResponse(ctx, clients)
	util.JSONResponse(w, response, http.StatusOK)
}

func GetClient(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	getClientRequestData, err := validation.ParseRequestToGetClientRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	client, err := models.GetClientByID(ctx, getClientRequestData.ClientID)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	if !user.Organization.IsClientOwner(ctx, client.ID) {
		util.ErrorResponse(w, "client does not belong to organization", http.StatusUnauthorized)
		return
	}

	response := responses.GenerateGetClientResponse(ctx, client)
	util.JSONResponse(w, response, http.StatusOK)
}

func CreateClient(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	createClientRequestData, err := validation.ParseRequestToCreateClientRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	client := models.Client{
		Name:           createClientRequestData.Name,
		OrganizationID: user.Organization.ID,
		Active:         true,
	}

	err = client.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create box", http.StatusBadRequest)
		return
	}

	//TODO create all initial things for client

	response := responses.GenerateCreateClientResponse(ctx, client)
	util.JSONResponse(w, response, http.StatusOK)
}

func UpdateClient(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	updateClientRequestData, err := validation.ParseRequestToUpdateClientRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	client, err := models.GetClientByID(ctx, updateClientRequestData.ClientID)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	if !user.Organization.IsClientOwner(ctx, client.ID) {
		util.ErrorResponse(w, "client does not belong to organization", http.StatusUnauthorized)
		return
	}

	if updateClientRequestData.Name != nil {
		client.Name = *updateClientRequestData.Name
	}

	err = client.Update(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to update client", http.StatusBadRequest)
		return
	}

	response := responses.GenerateUpdateClientResponse(ctx, client)
	util.JSONResponse(w, response, http.StatusOK)

}

func UpdateClientAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := models.UserFromContext(ctx)

	clientUpdateAvatarRequestData, err := validation.ParseRequestToClientUpdateAvatarRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	client, err := models.GetClientByID(ctx, clientUpdateAvatarRequestData.ClientID)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	if client.OrganizationID != user.Organization.ID {
		util.ErrorResponse(w, "client does not belong to organization", http.StatusUnauthorized)
		return
	}

	fileUUID := uuid.New()
	fileExtension := strings.Split(clientUpdateAvatarRequestData.FileType, "/")[1]

	err = api.S3FromContext(ctx).UploadFileToCDN(clientUpdateAvatarRequestData.Image, fileUUID.String(), fileExtension, clientUpdateAvatarRequestData.FileType)
	if err != nil {
		util.ErrorResponse(w, "failed to upload attachment to s3", http.StatusInternalServerError)
		return
	}

	client.AvatarFileName = fileUUID.String() + "." + fileExtension

	err = client.Update(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to update client", http.StatusBadRequest)
		return
	}

	response := responses.GenerateUpdateClientResponse(ctx, client)
	util.JSONResponse(w, response, http.StatusOK)
}
