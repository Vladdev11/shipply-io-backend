package models

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type LocationType struct {
	ID             int
	Name           string
	OrganizationID int

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	Organization Organization
	Locations    []Location
}

type LocationTypeReturnJSON struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	OrganizationID int    `json:"organization_id"`
}

type LocationTypeCreateRequest struct {
	Name string `json:"name"`
}

type LocationTypeUpdateRequest struct {
	Name string `json:"name"`
}

func (lt *LocationType) ConvertToReturnJSON() *LocationTypeReturnJSON {
	return &LocationTypeReturnJSON{
		ID:             lt.ID,
		Name:           lt.Name,
		OrganizationID: lt.OrganizationID,
	}
}

func (lt *LocationType) Create(ctx context.Context) error {
	err := util.DBFromContext(ctx).Create(lt).Error
	if err != nil {
		return err
	}

	return nil
}

func (lt *LocationType) Delete(ctx context.Context) error {
	err := util.DBFromContext(ctx).Delete(lt).Error
	if err != nil {
		return err
	}

	return nil
}

func (lt *LocationType) GetLocations(ctx context.Context) error {
	err := util.DBFromContext(ctx).Model(lt).Association("Locations").Find(&lt.Locations)
	if err != nil {
		return err
	}

	return nil
}

func (lt *LocationType) UpdateWithRequest(request LocationTypeUpdateRequest) error {

	if request.Name != "" {
		lt.Name = request.Name
	}

	return nil
}

func (ltcr *LocationTypeCreateRequest) ParseAndValidateRequest(r *http.Request) []string {
	ctx := r.Context()

	var errs []string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Name json.RawMessage `json:"name"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Name == nil {
		errs = append(errs, "name is required")
	} else if err := json.Unmarshal(aux.Name, &ltcr.Name); err != nil {
		errs = append(errs, "name must be a string")
	} else if len(ltcr.Name) > 255 {
		errs = append(errs, "name must be less than 255 characters")
	} else {
		user, err := GetRequestingUser(r)
		if err != nil {
			return []string{"invalid user"}
		}
		err = user.GetOrganization(ctx)
		if err != nil {
			return []string{"invalid organization"}
		}
		locationTypes, err := user.Organization.GetLocationTypes(ctx)
		if err != nil && err != gorm.ErrRecordNotFound {
			return []string{"invalid location types"}
		}
		for _, locationType := range locationTypes {
			if strings.ToLower(locationType.Name) == strings.ToLower(ltcr.Name) {
				errs = append(errs, "a location type with that name already exists")
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func (ltur *LocationTypeUpdateRequest) ParseAndValidateRequest(r *http.Request) []string {
	ctx := r.Context()

	var errs []string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Name json.RawMessage `json:"name"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	locationTypeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		return []string{"invalid location type id"}
	}

	if aux.Name != nil {
		if err := json.Unmarshal(aux.Name, &ltur.Name); err != nil {
			errs = append(errs, "name must be a string")
		} else if len(ltur.Name) > 255 {
			errs = append(errs, "name must be less than 255 characters")
		} else {
			user, err := GetRequestingUser(r)
			if err != nil {
				return []string{"invalid user"}
			}
			err = user.GetOrganization(ctx)
			if err != nil {
				return []string{"invalid organization"}
			}
			locationTypes, err := user.Organization.GetLocationTypes(ctx)
			if err != nil && err != gorm.ErrRecordNotFound {
				return []string{"invalid location types"}
			}
			for _, locationType := range locationTypes {
				if strings.ToLower(locationType.Name) == strings.ToLower(ltur.Name) && locationType.ID != locationTypeID {
					errs = append(errs, "a location type with that name already exists")
				}
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func GetLocationTypeByID(ctx context.Context, id int) (*LocationType, error) {
	locationType := &LocationType{}
	err := util.DBFromContext(ctx).First(locationType, id).Error
	if err != nil {
		return nil, err
	}

	return locationType, nil
}
