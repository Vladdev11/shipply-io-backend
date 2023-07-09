package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type User struct {
	ID             int
	FirstName      string
	LastName       string
	Email          string
	Password       []byte `gorm:"type:bytea"`
	Salt           []byte `gorm:"type:bytea"`
	OwnerID        int
	Role           int
	AvatarFileName string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Organization *Organization `gorm:"-"`
	Client       *Client       `gorm:"-"`
}

type UserReturnJSON struct {
	ID            int                     `json:"id"`
	FirstName     string                  `json:"first_name"`
	LastName      string                  `json:"last_name"`
	Email         string                  `json:"email"`
	Role          string                  `json:"role"`
	AvatarFileURL string                  `json:"avatar_file_name"`
	Organization  *OrganizationReturnJSON `json:"organization,omitempty"`
}

type UserCreateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Admin     bool   `json:"admin"`
}

func (u *User) GetRole() string {
	return util.UserRoleMap[u.Role]
}

func (u *User) GetOrganization(ctx context.Context) error {

	organization, err := GetOrganizationByID(ctx, u.OwnerID)
	if err != nil {
		return err
	}
	u.Organization = &organization
	return nil
}

func (u *User) GetClient(ctx context.Context) error {
	if u.GetRole() != util.ClientAdmin && u.GetRole() != util.ClientUser {
		return fmt.Errorf("user does not belong to a client")
	}

	client, err := GetClientByID(ctx, u.OwnerID)
	if err != nil {
		return err
	}
	u.Client = &client
	return nil
}

func (u *User) Delete(ctx context.Context) error {
	err := util.DBFromContext(ctx).Delete(u).Error
	return err
}

func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := util.DBFromContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func CreateUser(ctx context.Context, user *User) (*User, error) {
	err := util.DBFromContext(ctx).Create(&user).Error
	return user, err
}

func (u *User) Update(ctx context.Context) error {
	err := util.DBFromContext(ctx).Save(u).Error
	return err
}

func GetUserByID(ctx context.Context, id int) (User, error) {
	var user User
	err := util.DBFromContext(ctx).Where("id = ?", id).First(&user).Error
	return user, err
}

func GenerateUserPasswordAndSalt(password string) ([]byte, []byte, error) {
	//generate salt
	salt, err := util.GenerateSalt(64)
	if err != nil {
		return nil, nil, err
	}

	//hash password
	hashedPassword := HashPassword(password, salt)

	return hashedPassword, salt, err
}

func (u *User) ConvertToReturnJSON(ctx context.Context) *UserReturnJSON {
	return &UserReturnJSON{
		ID:            u.ID,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Email:         u.Email,
		Role:          u.GetRole(),
		Organization:  u.Organization.ConvertToReturnJSON(ctx),
		AvatarFileURL: fmt.Sprintf("%s/%s", util.CDNFromContext(ctx), u.AvatarFileName),
	}
}

func (u *UserCreateRequest) UnmarshalJSON(data []byte) error {
	aux := struct {
		FirstName json.RawMessage `json:"first_name"`
		LastName  json.RawMessage `json:"last_name"`
		Email     json.RawMessage `json:"email"`
		Password  json.RawMessage `json:"password"`
	}{}

	//unmarshal into aux struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	//required field, must be string
	if aux.FirstName != nil {
		if err := json.Unmarshal(aux.FirstName, &u.FirstName); err != nil {
			return fmt.Errorf("first name must be of type string")
		}
	} else {
		return fmt.Errorf("first name is required")
	}

	//required field, must be string
	if aux.LastName != nil {
		if err := json.Unmarshal(aux.LastName, &u.LastName); err != nil {
			return fmt.Errorf("last name must be of type string")
		}
	} else {
		return fmt.Errorf("last name is required")
	}

	//required field, must be string
	if aux.Email != nil {
		if err := json.Unmarshal(aux.Email, &u.Email); err != nil {
			return fmt.Errorf("email must be of type string")
		}
	} else {
		return fmt.Errorf("email is required")
	}

	//required field, must be string
	if aux.Password != nil {
		if err := json.Unmarshal(aux.Password, &u.Password); err != nil {
			return fmt.Errorf("password must be of type string")
		}
	} else {
		return fmt.Errorf("password is required")
	}

	return nil

}

func (u *User) Create(ctx context.Context) error {

	if err := util.DBFromContext(ctx).Create(&u).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserCreateRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errs []string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		FirstName json.RawMessage `json:"first_name"`
		LastName  json.RawMessage `json:"last_name"`
		Email     json.RawMessage `json:"email"`
		Password  json.RawMessage `json:"password"`
		Admin     json.RawMessage `json:"admin"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.FirstName != nil {
		if err := json.Unmarshal(aux.FirstName, &u.FirstName); err != nil {
			errs = append(errs, "first name must be of type string")
		}
	} else {
		errs = append(errs, "first name is required")
	}

	if aux.LastName != nil {
		if err := json.Unmarshal(aux.LastName, &u.LastName); err != nil {
			errs = append(errs, "last name must be of type string")
		}
	} else {
		errs = append(errs, "last name is required")
	}

	if aux.Email != nil {
		if err := json.Unmarshal(aux.Email, &u.Email); err != nil {
			errs = append(errs, "email must be of type string")
		}
	} else {
		errs = append(errs, "email is required")
	}

	user, err := GetUserByEmail(r.Context(), u.Email)
	//check if email is already in use
	if err == nil {
		if user != nil {
			errs = append(errs, "email is already in use")
		}
	}

	if aux.Password != nil {
		if err := json.Unmarshal(aux.Password, &u.Password); err != nil {
			errs = append(errs, "password must be of type string")
		} else {
			if !util.IsValidPassword(u.Password) {
				errs = append(errs, "password must be at least 8 characters, 1 uppercase, 1 lowercase, 1 number, 1 special character")
			}
		}
	} else {
		errs = append(errs, "password is required")
	}

	if aux.Admin != nil {
		if err := json.Unmarshal(aux.Admin, &u.Admin); err != nil {
			errs = append(errs, "admin must be of type boolean")
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func (u *UserCreateRequest) ConvertToUser() (*User, error) {

	password, salt, err := GenerateUserPasswordAndSalt(u.Password)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  password,
		Salt:      salt,
	}, nil
}

type UserUpdateAvatarRequest struct {
	File     multipart.File `json:"file"`
	FileType string         `json:"file_type"`
}

func (uur *UserUpdateAvatarRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errors []string

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errors = append(errors, err.Error())
		return errors
	}

	file, fileHeader, _ := r.FormFile("file")
	//check if there is a file
	if fileHeader != nil {
		//check if file is a valid image
		if fileHeader.Header.Get("Content-Type") != "image/jpeg" && fileHeader.Header.Get("Content-Type") != "image/png" {
			errors = append(errors, "file must be a valid image")
		}
		defer file.Close()
		uur.File = file
		uur.FileType = fileHeader.Header.Get("Content-Type")
	}

	if len(errors) > 0 {
		return errors
	}

	return nil

}

type UserUpdatePasswordRequest struct {
	Password    string `json:"password"`
	OldPassword string `json:"old_password"`
}

func (uupr *UserUpdatePasswordRequest) ParseAndValidateRequest(r *http.Request) []string {
	var errs []string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Password    json.RawMessage `json:"password"`
		OldPassword json.RawMessage `json:"old_password"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.OldPassword != nil {
		if err := json.Unmarshal(aux.OldPassword, &uupr.OldPassword); err != nil {
			errs = append(errs, "old password must be of type string")
		}
	} else {
		errs = append(errs, "password is required")
	}

	if aux.Password != nil {
		if err := json.Unmarshal(aux.Password, &uupr.Password); err != nil {
			errs = append(errs, "password must be of type string")
		} else {
			if !util.IsValidPassword(uupr.Password) {
				errs = append(errs, "password must be at least 8 characters, 1 uppercase, 1 lowercase, 1 number, 1 special character")
			}
		}
	} else {
		errs = append(errs, "password is required")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Admin     *bool  `json:"admin"`
}

func (uur *UserUpdateRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errs []string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		FirstName json.RawMessage `json:"first_name"`
		LastName  json.RawMessage `json:"last_name"`
		Admin     json.RawMessage `json:"admin"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.FirstName != nil {
		if err := json.Unmarshal(aux.FirstName, &uur.FirstName); err != nil {
			errs = append(errs, "first name must be of type string")
		}
	}

	if aux.LastName != nil {
		if err := json.Unmarshal(aux.LastName, &uur.LastName); err != nil {
			errs = append(errs, "last name must be of type string")
		}
	}

	if aux.Admin != nil {
		if err := json.Unmarshal(aux.Admin, &uur.Admin); err != nil {
			errs = append(errs, "admin must be of type boolean")
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}
