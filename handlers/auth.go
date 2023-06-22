package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/shipply-io/shipply-io-backend/models"
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
		util.ErrorResponse(w, "Failed To Receive Credentials", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByEmail(credentials.Email)
	if err != nil {
		util.ErrorResponse(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	// validate credentials
	if !models.ComparePassword(credentials.Password, user.Salt, user.Password) {
		util.ErrorResponse(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)

	// Initialize and add the claims to store data in the token
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["user_id"] = user.ID

	// sign the token
	// signed string is required to be a byte array
	tokenString, err := token.SignedString([]byte(util.ConfigAuthSecretSigning))
	if err != nil {
		util.ErrorResponse(w, "Failed To Sign Token", http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, map[string]string{"token": tokenString}, http.StatusOK)
}

func AuthNoExpirationToken(w http.ResponseWriter, r *http.Request) {

	var credentials Credentials

	// Decode the incoming credentials into the struct
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		util.ErrorResponse(w, "Failed To Receive Credentials", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByEmail(credentials.Email)
	if err != nil {
		util.ErrorResponse(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	// validate credentials
	if !models.ComparePassword(credentials.Password, user.Salt, user.Password) {
		util.ErrorResponse(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)

	// initialize and add the claims to store data in the token
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 2000000).Unix()
	claims["user_id"] = user.ID

	// sign the token
	// signed string is required to be a byte array
	tokenString, err := token.SignedString([]byte(util.ConfigAuthSecretSigning))
	if err != nil {
		util.ErrorResponse(w, "Failed To Sign Token", http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, map[string]string{"token": tokenString}, http.StatusOK)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {

	type ResetPasswordCredentials struct {
		Email string `json:"email"`
	}

	var resetPasswordCredentials ResetPasswordCredentials

	// Decode the incoming credentials into the struct
	err := json.NewDecoder(r.Body).Decode(&resetPasswordCredentials)
	if err != nil {
		util.ErrorResponse(w, "Failed To Receive Credentials", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByEmail(resetPasswordCredentials.Email)
	if err != nil {
		util.ErrorResponse(w, "invalid email", http.StatusUnauthorized)
		return
	}

	//create a password reset token
	passwordResetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     models.GenerateToken(),
		ExpiresAt: time.Now().Add(time.Hour * 1),
	}

	err = passwordResetToken.Create()
	if err != nil {
		util.ErrorResponse(w, "failed to create password reset token", http.StatusInternalServerError)
		return
	}

	//send email with password reset token
	err = SendgridAPI.SendPasswordResetEmail(user.FirstName, user.Email, passwordResetToken.Token)

	util.SuccessResponse(w, http.StatusOK)

}

func ValidateResetPasswordToken(w http.ResponseWriter, r *http.Request) {

	token, err := util.GetStringQueryParam(r, "token")
	if err != nil {
		util.ErrorResponse(w, "token is required", http.StatusBadRequest)
		return
	}

	passwordResetToken, err := models.GetPasswordResetTokenByToken(token)
	if err != nil {
		util.ErrorResponse(w, "invalid token", http.StatusUnauthorized)
		return
	}

	if passwordResetToken.ExpiresAt.Before(time.Now()) {
		util.ErrorResponse(w, "token expired", http.StatusUnauthorized)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {

	token, err := util.GetStringQueryParam(r, "token")
	if err != nil {
		util.ErrorResponse(w, "token is required", http.StatusBadRequest)
		return
	}

	passwordResetToken, err := models.GetPasswordResetTokenByToken(token)
	if err != nil {
		util.ErrorResponse(w, "invalid token", http.StatusUnauthorized)
		return
	}

	if passwordResetToken.ExpiresAt.Before(time.Now()) {
		util.ErrorResponse(w, "token expired", http.StatusUnauthorized)
		return
	}

	type UpdatePasswordCredentials struct {
		Password string `json:"password"`
	}

	var updatePasswordCredentials UpdatePasswordCredentials

	err = json.NewDecoder(r.Body).Decode(&updatePasswordCredentials)
	if err != nil {
		util.ErrorResponse(w, "Failed To Receive Credentials", http.StatusBadRequest)
		return
	}

	if !util.IsValidPassword(updatePasswordCredentials.Password) {
		util.ErrorResponse(w, "invalid password", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByID(passwordResetToken.UserID)
	if err != nil {
		util.ErrorResponse(w, "invalid user", http.StatusUnauthorized)
		return
	}

	hashedPassword, salt, err := models.GenerateUserPasswordAndSalt(updatePasswordCredentials.Password)
	if err != nil {
		util.ErrorResponse(w, "failed to generate password", http.StatusInternalServerError)
		return
	}

	user.Password = hashedPassword
	user.Salt = salt

	err = user.Update()
	if err != nil {
		util.ErrorResponse(w, "failed to update password", http.StatusInternalServerError)
		return
	}

	err = passwordResetToken.Delete()
	if err != nil {
		util.ErrorResponse(w, "failed to delete password reset token", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}
