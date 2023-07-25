package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"

	SendgridAPI "github.com/shipply-io/shipply-io-backend/api/sendgrid"
)

// Credentials is the struct for the login credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthLogin is the handler for the login route
// We verfiy the credentials and return a signed JWT token
func AuthLogin(w http.ResponseWriter, r *http.Request) {

	var credentials Credentials

	// Decode the incoming credentials into the struct
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		util.ErrResponse(w, ErrFailedToReceiveCredentials, http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByEmail(r.Context(), credentials.Email)
	if err != nil {
		util.ErrResponse(w, ErrInvalidCredentials, http.StatusUnauthorized)
		return
	}

	// validate credentials
	if !models.ComparePassword(credentials.Password, user.Salt, user.Password) {
		util.ErrResponse(w, ErrInvalidCredentials, http.StatusUnauthorized)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)

	// Initialize and add the claims to store data in the token
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["user_id"] = user.ID

	// sign the token
	// signed string is required to be a byte array
	tokenString, err := token.SignedString([]byte(util.AuthSecretFromContext(r.Context())))
	if err != nil {
		util.ErrResponse(w, ErrFailedTokenSignature, http.StatusInternalServerError)
		return
	}

	response := responses.GenerateAuthLoginResponse(tokenString)
	util.JSONResponse(w, response, http.StatusOK)
}

func AuthNoExpirationToken(w http.ResponseWriter, r *http.Request) {

	var credentials Credentials

	// Decode the incoming credentials into the struct
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		util.ErrResponse(w, ErrFailedToReceiveCredentials, http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByEmail(r.Context(), credentials.Email)
	if err != nil {
		util.ErrResponse(w, ErrInvalidCredentials, http.StatusUnauthorized)
		return
	}

	// validate credentials
	if !models.ComparePassword(credentials.Password, user.Salt, user.Password) {
		util.ErrResponse(w, ErrInvalidCredentials, http.StatusUnauthorized)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)

	// initialize and add the claims to store data in the token
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 2000000).Unix()
	claims["user_id"] = user.ID

	// sign the token
	// signed string is required to be a byte array
	tokenString, err := token.SignedString([]byte(util.AuthSecretFromContext(r.Context())))
	if err != nil {
		util.ErrResponse(w, ErrFailedTokenSignature, http.StatusInternalServerError)
		return
	}

	response := responses.GenerateAuthNoExpirationTokenResponse(tokenString)
	util.JSONResponse(w, response, http.StatusOK)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type ResetPasswordCredentials struct {
		Email string `json:"email"`
	}

	var resetPasswordCredentials ResetPasswordCredentials

	// Decode the incoming credentials into the struct
	err := json.NewDecoder(r.Body).Decode(&resetPasswordCredentials)
	if err != nil {
		util.ErrResponse(w, ErrFailedToReceiveCredentials, http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByEmail(ctx, resetPasswordCredentials.Email)
	if err != nil {
		util.ErrResponse(w, ErrUserEmailDoesNotExist, http.StatusUnauthorized)
		return
	}

	//create a password reset token
	passwordResetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     models.GenerateToken(),
		ExpiresAt: time.Now().Add(time.Hour * 1),
	}

	err = passwordResetToken.Create(ctx)
	if err != nil {
		util.ErrResponse(w, ErrResetPasswordTokenGeneration, http.StatusInternalServerError)
		return
	}

	//send email with password reset token
	err = SendgridAPI.SendPasswordResetEmail(ctx, user.FirstName, user.Email, passwordResetToken.Token)

	util.SuccessResponse(w, http.StatusOK)

}

func ValidateResetPasswordToken(w http.ResponseWriter, r *http.Request) {

	token, err := util.GetStringQueryParam(r, "token")
	if err != nil {
		util.ErrResponse(w, ErrTokenRequired, http.StatusBadRequest)
		return
	}

	passwordResetToken, err := models.GetPasswordResetTokenByToken(r.Context(), token)
	if err != nil {
		util.ErrResponse(w, ErrInvalidToken, http.StatusUnauthorized)
		return
	}

	if passwordResetToken.ExpiresAt.Before(time.Now()) {
		util.ErrResponse(w, ErrTokenExpired, http.StatusUnauthorized)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token, err := util.GetStringQueryParam(r, "token")
	if err != nil {
		util.ErrResponse(w, ErrTokenRequired, http.StatusBadRequest)
		return
	}

	passwordResetToken, err := models.GetPasswordResetTokenByToken(ctx, token)
	if err != nil {
		util.ErrResponse(w, ErrInvalidToken, http.StatusUnauthorized)
		return
	}

	if passwordResetToken.ExpiresAt.Before(time.Now()) {
		util.ErrResponse(w, ErrTokenExpired, http.StatusUnauthorized)
		return
	}

	type UpdatePasswordCredentials struct {
		Password string `json:"password"`
	}

	var updatePasswordCredentials UpdatePasswordCredentials

	err = json.NewDecoder(r.Body).Decode(&updatePasswordCredentials)
	if err != nil {
		util.ErrResponse(w, ErrFailedToReceiveCredentials, http.StatusBadRequest)
		return
	}

	if !util.IsValidPassword(updatePasswordCredentials.Password) {
		util.ErrResponse(w, ErrInvalidPassword, http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByID(ctx, passwordResetToken.UserID)
	if err != nil {
		util.ErrResponse(w, ErrInvalidUser, http.StatusUnauthorized)
		return
	}

	hashedPassword, salt, err := models.GenerateUserPasswordAndSalt(updatePasswordCredentials.Password)
	if err != nil {
		util.ErrResponse(w, ErrFailedPasswordGeneration, http.StatusInternalServerError)
		return
	}

	user.Password = hashedPassword
	user.Salt = salt

	err = user.Update(ctx)
	if err != nil {
		util.ErrResponse(w, ErrUpdatePassword, http.StatusInternalServerError)
		return
	}

	err = passwordResetToken.Delete(ctx)
	if err != nil {
		util.ErrResponse(w, ErrDeletePasswordResetToken, http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}
