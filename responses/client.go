package responses

import (
	"context"

	"github.com/shipply-io/shipply-io-backend/models"
)

/* -------------------------------- GetClient ------------------------------- */

// GetClientResponse represents the expected response body for the GetClient endpoint
type GetClientResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
	AvatarURL string `json:"avatar_url"`
}

// GenerateGetClientResponse generates the response body for the GetClient endpoint
func GenerateGetClientResponse(ctx context.Context, client models.Client) *GetClientResponse {

	return &GetClientResponse{
		ID:        client.ID,
		Name:      client.Name,
		Active:    client.Active,
		AvatarURL: client.GetAvatarFileURL(ctx),
	}

}

/* ------------------------------- ListClients ------------------------------ */

// ListClientsResponse represents the expected response body for the ListClients endpoint
type ListClientsResponse struct {
	Clients []ClientResponseForListClients `json:"clients"`
}

// ClientResponseForListClients represents the expected response body for an individual client in the ListClients endpoint
type ClientResponseForListClients struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
	AvatarURL string `json:"avatar_url"`
}

// GenerateListClientsResponse generates the response body for the ListClients endpoint
func GenerateListClientsResponse(ctx context.Context, clients []models.Client) *ListClientsResponse {

	clientResponses := []ClientResponseForListClients{}
	for _, client := range clients {

		clientResponses = append(clientResponses, ClientResponseForListClients{
			ID:        client.ID,
			Name:      client.Name,
			Active:    client.Active,
			AvatarURL: client.GetAvatarFileURL(ctx),
		})
	}

	return &ListClientsResponse{
		Clients: clientResponses,
	}

}

/* ------------------------------ CreateClient ------------------------------ */

// CreateClientResponse represents the expected response body for the CreateClient endpoint
type CreateClientResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
	AvatarURL string `json:"avatar_url"`
}

// GenerateCreateClientResponse generates the response body for the CreateClient endpoint
func GenerateCreateClientResponse(ctx context.Context, client models.Client) *CreateClientResponse {

	return &CreateClientResponse{
		ID:        client.ID,
		Name:      client.Name,
		Active:    client.Active,
		AvatarURL: client.GetAvatarFileURL(ctx),
	}

}

/* ------------------------------ UpdateClient ------------------------------ */

// UpdateClientResponse represents the expected response body for the UpdateClient endpoint
type UpdateClientResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
	AvatarURL string `json:"avatar_url"`
}

// GenerateUpdateClientResponse generates the response body for the UpdateClient endpoint
func GenerateUpdateClientResponse(ctx context.Context, client models.Client) *UpdateClientResponse {

	return &UpdateClientResponse{
		ID:        client.ID,
		Name:      client.Name,
		Active:    client.Active,
		AvatarURL: client.GetAvatarFileURL(ctx),
	}

}

/* --------------------------- UpdateClientAvatar --------------------------- */

// UpdateClientAvatarResponse represents the expected response body for the UpdateClientAvatar endpoint
type UpdateClientAvatarResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
	AvatarURL string `json:"avatar_url"`
}

// GenerateUpdateClientAvatarResponse generates the response body for the UpdateClientAvatar endpoint
func GenerateUpdateClientAvatarResponse(ctx context.Context, client models.Client) *UpdateClientAvatarResponse {

	return &UpdateClientAvatarResponse{
		ID:        client.ID,
		Name:      client.Name,
		Active:    client.Active,
		AvatarURL: client.GetAvatarFileURL(ctx),
	}

}
