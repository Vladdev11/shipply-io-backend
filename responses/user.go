package responses

import (
	"context"

	"github.com/shipply-io/shipply-io-backend/models"
)

/* --------------------------------- GetUserSelf -------------------------------- */

// OrganizationalGetUserSelfResponse represents the expected response body for the GetUser endpoint for an organizational user
type OrganizationalGetUserSelfResponse struct {
	ID            int                                              `json:"id"`
	FirstName     string                                           `json:"first_name"`
	LastName      string                                           `json:"last_name"`
	Email         string                                           `json:"email"`
	Role          string                                           `json:"role"`
	AvatarFileURL string                                           `json:"avatar_file_url"`
	Organization  OrganizationResponseForOrganizationalGetUserSelf `json:"organization"`
}

// OrganizationResponseForOrganizationalGetUserSelf represents the expected response body for the Organization in the GetUserSelf endpoint for an organizational user
type OrganizationResponseForOrganizationalGetUserSelf struct {
	ID      int                                                                 `json:"id"`
	Name    string                                                              `json:"name"`
	Type    string                                                              `json:"type"`
	Clients []ClientResponseForOrganizationResponseForOrganizationalGetUserSelf `json:"clients"`
}

// ClientResponseForOrganizationResponseForOrganizationalGetUserSelf represents the expected response body for an individual client in the Organization in the GetUserSelf endpoint for an organizational user
type ClientResponseForOrganizationResponseForOrganizationalGetUserSelf struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
	AvatarURL string `json:"avatar_url"`
}

// GenerateOrganizationalGetUserSelfResponse generates the response body for the GetUser endpoint for an organizational user
func GenerateOrganizationalGetUserSelfResponse(ctx context.Context, user models.User) *OrganizationalGetUserSelfResponse {

	clients := []ClientResponseForOrganizationResponseForOrganizationalGetUserSelf{}
	for _, client := range user.Organization.Clients {

		clients = append(clients, ClientResponseForOrganizationResponseForOrganizationalGetUserSelf{
			ID:        client.ID,
			Name:      client.Name,
			Active:    client.Active,
			AvatarURL: client.GetAvatarFileURL(ctx),
		})
	}

	return &OrganizationalGetUserSelfResponse{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Role:          user.GetRole(),
		AvatarFileURL: user.GetAvatarFileURL(ctx),
		Organization: OrganizationResponseForOrganizationalGetUserSelf{
			ID:      user.Organization.ID,
			Name:    user.Organization.Name,
			Type:    user.Organization.GetType(),
			Clients: clients,
		},
	}

}

// ClientGetUserSelfResponse represents the expected response body for the GetUser endpoint for a client user
type ClientGetUserSelfResponse struct {
	ID            int                                              `json:"id"`
	FirstName     string                                           `json:"first_name"`
	LastName      string                                           `json:"last_name"`
	Email         string                                           `json:"email"`
	Role          string                                           `json:"role"`
	AvatarFileURL string                                           `json:"avatar_file_url"`
	Client        ClientResponseForClientGetUserSelfResponse       `json:"client"`
	Organization  OrganizationResponseForClientGetUserSelfResponse `json:"organization"`
}

// ClientResponseForClientGetUserSelfResponse represents the expected response body for the Client in the GetUserSelf endpoint for a client user
type ClientResponseForClientGetUserSelfResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
	AvatarURL string `json:"avatar_url"`
}

// OrganizationResponseForClientGetUserSelfResponse represents the expected response body for the Organization in the GetUserSelf endpoint for a client user
type OrganizationResponseForClientGetUserSelfResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// GenerateClientGetUserSelfResponse generates the response body for the GetUser endpoint for a client user
func GenerateClientGetUserSelfResponse(ctx context.Context, user models.User) *ClientGetUserSelfResponse {

	return &ClientGetUserSelfResponse{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Role:          user.GetRole(),
		AvatarFileURL: user.GetAvatarFileURL(ctx),
		Client: ClientResponseForClientGetUserSelfResponse{
			ID:        user.Client.ID,
			Name:      user.Client.Name,
			Active:    user.Client.Active,
			AvatarURL: user.Client.GetAvatarFileURL(ctx),
		},
		Organization: OrganizationResponseForClientGetUserSelfResponse{
			ID:   user.Client.Organization.ID,
			Name: user.Client.Organization.Name,
			Type: user.Client.Organization.GetType(),
		},
	}

}

/* --------------------------------- GetUser -------------------------------- */

// GetUserResponse represents the expected response body for the GetUser endpoint
type GetUserResponse struct {
	ID            int    `json:"id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	AvatarFileURL string `json:"avatar_file_url"`
}

// GenerateGetUserResponse generates the response body for the GetUser endpoint
func GenerateGetUserResponse(ctx context.Context, user models.User) *GetUserResponse {

	return &GetUserResponse{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Role:          user.GetRole(),
		AvatarFileURL: user.GetAvatarFileURL(ctx),
	}

}

/* ------------------------------- CreateUser ------------------------------- */

// CreateUserResponse represents the expected response body for the CreateUser endpoint
type CreateUserResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

// GenerateCreateUserResponse generates the response body for the CreateUser endpoint
func GenerateCreateUserResponse(user models.User) *CreateUserResponse {
	return &CreateUserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.GetRole(),
	}
}

/* ------------------------------- UpdateUser ------------------------------- */

// UpdateUserResponse represents the expected response body for the UpdateUser endpoint
type UpdateUserResponse struct {
	ID            int    `json:"id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	AvatarFileURL string `json:"avatar_file_url"`
}

// GenerateUpdateUserResponse generates the response body for the UpdateUser endpoint
func GenerateUpdateUserResponse(ctx context.Context, user models.User) *UpdateUserResponse {

	return &UpdateUserResponse{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Role:          user.GetRole(),
		AvatarFileURL: user.GetAvatarFileURL(ctx),
	}

}

/* ---------------------------- UpdateUserAvatar ---------------------------- */

// UpdateUserAvatarResponse represents the expected response body for the UpdateUserAvatar endpoint
type UpdateUserAvatarResponse struct {
	AvatarFileURL string `json:"avatar_file_url"`
}

// GenerateUpdateUserAvatarResponse generates the response body for the UpdateUserAvatar endpoint
func GenerateUpdateUserAvatarResponse(ctx context.Context, user models.User) *UpdateUserAvatarResponse {

	return &UpdateUserAvatarResponse{
		AvatarFileURL: user.GetAvatarFileURL(ctx),
	}

}
