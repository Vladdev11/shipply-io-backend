package OrganizationHandlers

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/shipply-io/shipply-io-backend/api"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func UserGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusBadRequest)
		return
	}

	err = user.Organization.GetClients(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get clients", http.StatusBadRequest)
		return
	}

	userJSON := user.ConvertToReturnJSON(ctx)

	util.JSONResponse(w, userJSON, http.StatusOK)
}

func UserGetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestingUser, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "invalid login", http.StatusUnauthorized)
		return
	}

	userID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByID(ctx, userID)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	if requestingUser.OwnerID != user.OwnerID {
		util.ErrorResponse(w, "user does not belong to your organization", http.StatusForbidden)
		return
	}

	userJSON := user.ConvertToReturnJSON(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to convert user to json", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, userJSON, http.StatusOK)
}

func UserCreate(w http.ResponseWriter, r *http.Request) {
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
		util.ErrorResponse(w, "failed to generate password and salt", http.StatusBadRequest)
		return
	}

	role := util.OrganizationUserInt
	if request.Admin {
		role = util.OrganizationAdminInt
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
		util.ErrorResponse(w, "failed to create user", http.StatusBadRequest)
		return
	}

	userJSON := user.ConvertToReturnJSON(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to convert user to json", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, userJSON, http.StatusOK)
}

func UserUpdateAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	userID, err := util.GetIntFromPath(r, "user_id")
	if err != nil {
		util.ErrorResponse(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	userFromRequest, err := models.GetUserByID(ctx, userID)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	if userFromRequest.OwnerID != user.OwnerID {
		util.ErrorResponse(w, "user does not belong to your organization", http.StatusForbidden)
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
			util.ErrorResponse(w, "failed to upload attachment to s3", http.StatusInternalServerError)
			return
		}

		userFromRequest.AvatarFileName = fileUUID.String() + "." + fileExtension
	}

	err = userFromRequest.Update(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to update user", http.StatusBadRequest)
		return
	}

	userJSON := userFromRequest.ConvertToReturnJSON(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to convert user to json", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, userJSON, http.StatusOK)

}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	userID, err := util.GetIntFromPath(r, "user_id")
	if err != nil {
		util.ErrorResponse(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	userFromRequest, err := models.GetUserByID(ctx, userID)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	if userFromRequest.ID == user.ID {
		util.ErrorResponse(w, "cannot delete yourself", http.StatusBadRequest)
		return
	}

	if userFromRequest.OwnerID != user.OwnerID {
		util.ErrorResponse(w, "user does not belong to your organization", http.StatusForbidden)
		return
	}

	err = userFromRequest.Delete(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to delete user", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func UserUpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusBadRequest)
		return
	}

	request := models.UserUpdatePasswordRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	if !models.ComparePassword(request.OldPassword, user.Salt, user.Password) {
		util.ErrorResponse(w, "old password is incorrect", http.StatusBadRequest)
		return
	}

	//generate password and salt
	hashedPassword, salt, err := models.GenerateUserPasswordAndSalt(request.Password)
	if err != nil {
		util.ErrorResponse(w, "failed to generate password and salt", http.StatusBadRequest)
		return
	}

	user.Password = hashedPassword
	user.Salt = salt

	err = user.Update(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to update user", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusBadRequest)
		return
	}

	userID, err := util.GetIntFromPath(r, "user_id")
	if err != nil {
		util.ErrorResponse(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	userFromRequest, err := models.GetUserByID(ctx, userID)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	request := models.UserUpdateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	if request.FirstName != "" {
		userFromRequest.FirstName = request.FirstName
	}

	if request.LastName != "" {
		userFromRequest.LastName = request.LastName
	}

	if request.Admin != nil {
		if *request.Admin {
			userFromRequest.Role = util.OrganizationAdminInt
		} else {
			if userFromRequest.ID == user.ID {
				util.ErrorResponse(w, "cannot change your own admin status", http.StatusBadRequest)
				return
			}
			userFromRequest.Role = util.OrganizationUserInt
		}
	}

	err = userFromRequest.Update(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to update user", http.StatusBadRequest)
		return
	}

	userJSON := user.ConvertToReturnJSON(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to convert user to json", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, userJSON, http.StatusOK)

}
