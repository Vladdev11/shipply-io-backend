package ClientHandlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/shipply-io/shipply-io-backend/api"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
)

var (
	//ErrUserDoesNotHaveAccessToUser is returned when a user does not have access to a user
	ErrUserDoesNotHaveAccessToUser = errors.New("user does not have access to user")
	//ErrCannotDeleteSelf is returned when a user tries to delete themselves
	ErrCannotDeleteSelf = errors.New("cannot delete self")
	//ErrCannotChangeOwnAdminStatus is returned when a user tries to change their own admin status
	ErrCannotChangeOwnAdminStatus = errors.New("cannot change own admin status")
	//ErrOldPasswordIncorrect is returned when a user tries to change their password and the old password is incorrect
	ErrOldPasswordIncorrect = errors.New("old password is incorrect")
)

func GetUserSelf(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusBadRequest)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = user.Client.GetOrganization(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	response := responses.GenerateClientGetUserSelfResponse(ctx, *user)
	util.JSONResponse(w, response, http.StatusOK)
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	requestingUser, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "invalid login", http.StatusUnauthorized)
		return
	}

	userID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByID(ctx, userID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	//If the user is not the requesting user, they must be an admin
	if requestingUser.OwnerID != user.OwnerID || (user.Role != util.ClientAdminInt && user.Role != util.ClientUserInt) {
		util.ErrResponse(w, ErrUserDoesNotHaveAccessToUser, http.StatusForbidden)
		return
	}

	response := responses.GenerateGetUserResponse(ctx, user)
	util.JSONResponse(w, response, http.StatusOK)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestingUser, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	request := models.UserCreateRequest{}

	validationErrors := request.ParseAndValidateRequest(r)
	if len(validationErrors) > 0 {
		util.ErrorsResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	//generate password and salt
	hashedPassword, salt, err := models.GenerateUserPasswordAndSalt(request.Password)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	role := util.ClientUserInt
	if request.Admin {
		role = util.ClientAdminInt
	}

	user := &models.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  hashedPassword,
		Salt:      salt,
		OwnerID:   requestingUser.OwnerID,
		Role:      role,
	}

	err = user.Create(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	response := responses.GenerateCreateUserResponse(*user)
	util.JSONResponse(w, response, http.StatusOK)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	userID, err := util.GetIntFromPath(r, "user_id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	userForUpdate, err := models.GetUserByID(ctx, userID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if userForUpdate.OwnerID != user.OwnerID || (userForUpdate.Role != util.ClientAdminInt && userForUpdate.Role != util.ClientUserInt) {
		util.ErrResponse(w, ErrUserDoesNotHaveAccessToUser, http.StatusForbidden)
		return
	}

	request := models.UserUpdateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	if request.FirstName != "" {
		userForUpdate.FirstName = request.FirstName
	}

	if request.LastName != "" {
		userForUpdate.LastName = request.LastName
	}

	if request.Admin != nil {
		if *request.Admin {
			userForUpdate.Role = util.ClientAdminInt
		} else {
			if userForUpdate.ID == user.ID {
				util.ErrResponse(w, ErrCannotChangeOwnAdminStatus, http.StatusBadRequest)
				return
			}
			userForUpdate.Role = util.ClientUserInt
		}
	}

	err = userForUpdate.Update(ctx)

	response := responses.GenerateUpdateUserResponse(ctx, userForUpdate)
	util.JSONResponse(w, response, http.StatusOK)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	userID, err := util.GetIntFromPath(r, "user_id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	userForDeletion, err := models.GetUserByID(ctx, userID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if userForDeletion.ID == user.ID {
		util.ErrResponse(w, ErrCannotDeleteSelf, http.StatusBadRequest)
		return
	}

	if userForDeletion.OwnerID != user.OwnerID || (userForDeletion.Role != util.ClientAdminInt && userForDeletion.Role != util.ClientUserInt) {
		util.ErrResponse(w, ErrUserDoesNotHaveAccessToUser, http.StatusForbidden)
		return
	}

	err = userForDeletion.Delete(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	request := models.UserUpdatePasswordRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	if !models.ComparePassword(request.OldPassword, user.Salt, user.Password) {
		util.ErrResponse(w, ErrOldPasswordIncorrect, http.StatusBadRequest)
		return
	}

	//generate password and salt
	hashedPassword, salt, err := models.GenerateUserPasswordAndSalt(request.Password)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	user.Password = hashedPassword
	user.Salt = salt

	err = user.Update(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func UpdateUserAvatar(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	userID, err := util.GetIntFromPath(r, "user_id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	userForAvatarUpdate, err := models.GetUserByID(ctx, userID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if userForAvatarUpdate.OwnerID != user.OwnerID || (userForAvatarUpdate.Role != util.ClientAdminInt && userForAvatarUpdate.Role != util.ClientUserInt) {
		util.ErrResponse(w, ErrUserDoesNotHaveAccessToUser, http.StatusForbidden)
		return
	}

	request := models.UserUpdateAvatarRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	if request.File != nil {
		fileUUID := uuid.New()
		fileExtension := strings.Split(request.FileType, "/")[1]

		err = api.S3FromContext(ctx).UploadFileToCDN(request.File, fileUUID.String(), fileExtension, request.FileType)
		if err != nil {
			util.ErrResponse(w, err, http.StatusInternalServerError)
			return
		}

		userForAvatarUpdate.AvatarFileName = fileUUID.String() + "." + fileExtension
	}

	err = userForAvatarUpdate.Update(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	response := responses.GenerateUpdateUserAvatarResponse(ctx, userForAvatarUpdate)
	util.JSONResponse(w, response, http.StatusOK)
}
