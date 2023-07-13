package responses

/* -------------------------------- AuthLogin ------------------------------- */

// AuthLoginResponse represents the expected response body for the AuthLogin endpoint
type AuthLoginResponse struct {
	Token string `json:"token"`
}

// GenerateAuthLoginResponse generates the response body for the AuthLogin endpoint
func GenerateAuthLoginResponse(token string) *AuthLoginResponse {
	return &AuthLoginResponse{
		Token: token,
	}
}

/* -------------------------- AuthNoExpirationToken ------------------------- */

// AuthNoExpirationTokenResponse represents the expected response body for the AuthNoExpirationToken endpoint
type AuthNoExpirationTokenResponse struct {
	Token string `json:"token"`
}

// GenerateAuthNoExpirationTokenResponse generates the response body for the AuthNoExpirationToken endpoint
func GenerateAuthNoExpirationTokenResponse(token string) *AuthNoExpirationTokenResponse {
	return &AuthNoExpirationTokenResponse{
		Token: token,
	}
}
