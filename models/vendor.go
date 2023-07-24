package models

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type Vendor struct {
	ID              int
	ClientID        int
	Name            string
	VendorAccountID string

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	Client Client
}

type VendorReturnJSON struct {
	ID       int    `json:"id"`
	ClientID int    `json:"client_id"`
	Name     string `json:"name"`
}

type VendorListRequest struct {
	ClientID       int    `json:"client_id"`
	OrganizationID int    `json:"organization_id"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	OrderBy        string `json:"order_by"`
	OrderByColumn  string `json:"order_by_column"`
	SearchValue    string `json:"search_value"`
}

type VendorCreateRequest struct {
	ClientID        int    `json:"client_id"`
	Name            string `json:"name"`
	VendorAccountID string `json:"vendor_account_id"`
}

type VendorUpdateRequest struct {
	Name            string `json:"name"`
	VendorAccountID string `json:"vendor_account_id"`
}

func (v *Vendor) Create(ctx context.Context) error {
	err := util.DBFromContext(ctx).Create(v).Error
	if err != nil {
		return ErrCreateFailed{Object: "vendor", Err: err}
	}
	return nil
}

func (v *Vendor) Delete(ctx context.Context) error {
	err := util.DBFromContext(ctx).Delete(v).Error
	if err != nil {
		return ErrDeleteFailed{Object: "vendor", Err: err}
	}

	return nil
}

func (v *Vendor) UpdateWithRequest(ctx context.Context, request *VendorUpdateRequest) error {

	if request.Name != "" {
		v.Name = request.Name
	}

	if request.VendorAccountID != "" {
		v.VendorAccountID = request.VendorAccountID
	}

	if err := util.DBFromContext(ctx).Save(v).Error; err != nil {
		return err
	}

	return nil
}

func (v *Vendor) ConvertToReturnJSON() *VendorReturnJSON {

	if v == nil {
		return &VendorReturnJSON{}
	}

	return &VendorReturnJSON{
		ID:       v.ID,
		ClientID: v.ClientID,
		Name:     v.Name,
	}
}

func (v *Vendor) GetPurchaseOrders(ctx context.Context) ([]PurchaseOrder, error) {
	var purchaseOrders []PurchaseOrder

	err := util.DBFromContext(ctx).Where("vendor_id = ?", v.ID).Find(&purchaseOrders).Error
	if err != nil {
		return nil, err
	}

	return purchaseOrders, nil
}

func (v *VendorListRequest) ParseAndValidateRequest(r *http.Request) []string {

	errors := []string{}

	clientID, err := util.GetIntQueryParam(r, "client_id")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "client_id must be an integer")
		} else {
			v.ClientID = clientID
		}
	}

	v.Limit = 100
	limit, err := util.GetIntQueryParam(r, "limit")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "limit must be an integer")
		} else if limit < 0 {
			errors = append(errors, "limit must be greater than or equal to 0")
		} else {
			v.Limit = limit
		}
	}

	v.Offset = 0
	offset, err := util.GetIntQueryParam(r, "offset")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "offset must be an integer")
		} else if offset < 0 {
			errors = append(errors, "offset must be greater than or equal to 0")
		} else {
			v.Offset = offset
		}
	}

	v.OrderBy = "asc"
	orderBy, err := util.GetStringQueryParam(r, "order")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "order must be a string")
		} else if orderBy != "asc" && orderBy != "desc" {
			errors = append(errors, "order must be either asc or desc")
		} else {
			v.OrderBy = orderBy
		}
	}

	searchValue, err := util.GetStringQueryParam(r, "search_value")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "search_value must be a string")
		} else {
			v.SearchValue = searchValue
		}
	}

	orderByColumn, err := util.GetStringQueryParam(r, "order_by_column")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "order_by_column must be a string")
		} else {
			v.OrderByColumn = orderByColumn
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (v *VendorCreateRequest) ParseAndValidateRequest(r *http.Request) []string {
	ctx := r.Context()

	var errs []string

	user, err := GetRequestingUser(r)
	if err != nil {
		return []string{"invalid user"}
	}
	user.GetClient(ctx)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		ClientID        json.RawMessage `json:"client_id"`
		Name            json.RawMessage `json:"name"`
		VendorAccountID json.RawMessage `json:"vendor_account_id"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if user.GetRole() == "organization_admin" || user.GetRole() == "organization_user" {
		if aux.ClientID == nil {
			errs = append(errs, "client_id is required")
		} else if err := json.Unmarshal(aux.ClientID, &v.ClientID); err != nil {
			errs = append(errs, "client_id must be an integer")
		} else {
			if v.ClientID < 1 {
				errs = append(errs, "client_id must be greater than 0")
			} else {
				_, err := GetClientByID(ctx, v.ClientID)
				if err != nil {
					errs = append(errs, "client_id does not exist")
				}
			}
		}
	} else {
		v.ClientID = user.Client.ID
	}

	if aux.Name == nil {
		errs = append(errs, "name is required")
	} else if err := json.Unmarshal(aux.Name, &v.Name); err != nil {
		errs = append(errs, "name must be a string")
	} else if len(v.Name) < 1 {
		errs = append(errs, "name must be at least 1 character")
	} else {
		vendors, err := GetVendorsByClientID(ctx, v.ClientID)
		if err != nil {
			errs = append(errs, "error getting vendors")
		}

		for _, vendor := range vendors {
			if vendor.Name == v.Name {
				errs = append(errs, "vendor with that name already exists")
			}
		}
	}

	if aux.VendorAccountID != nil {
		if err := json.Unmarshal(aux.VendorAccountID, &v.VendorAccountID); err != nil {
			errs = append(errs, "vendor_account_id must be a string")
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (v *VendorUpdateRequest) ParseAndValidateRequest(r *http.Request) []string {
	ctx := r.Context()

	var errs []string

	vendorID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		return []string{"invalid vendor id"}
	}

	vendor, err := GetVendorByID(ctx, vendorID)
	if err != nil {
		return []string{"vendor does not exist"}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Name            json.RawMessage `json:"name"`
		VendorAccountID json.RawMessage `json:"vendor_account_id"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Name != nil {
		if err := json.Unmarshal(aux.Name, &v.Name); err != nil {
			errs = append(errs, "name must be a string")
		} else if len(v.Name) > 255 {
			errs = append(errs, "name must be less than 255 characters")
		} else {
			vendors, err := GetVendorsByClientID(ctx, vendor.ClientID)
			if err != nil {
				errs = append(errs, "error getting vendors")
			}

			for _, vendor := range vendors {
				if vendor.Name == v.Name {
					errs = append(errs, "vendor with that name already exists")
				}
			}
		}
	}

	if aux.VendorAccountID != nil {
		if err := json.Unmarshal(aux.VendorAccountID, &v.VendorAccountID); err != nil {
			errs = append(errs, "vendor_account_id must be a string")
		} else if len(v.VendorAccountID) > 255 {
			errs = append(errs, "vendor_account_id must be less than 255 characters")
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func GetVendorByID(ctx context.Context, id int) (*Vendor, error) {
	vendor := &Vendor{}

	err := util.DBFromContext(ctx).First(vendor, id).Error
	if err != nil {
		return nil, ErrQueryFailed{Object: "vendor", Err: err}
	}

	//return vendor
	return vendor, nil
}

func GetVendorsByClientID(ctx context.Context, clientID int) ([]*Vendor, error) {
	vendors := []*Vendor{}

	err := util.DBFromContext(ctx).Where("client_id = ?", clientID).Find(&vendors).Error
	if err != nil {
		return nil, err
	}

	return vendors, nil
}

func ConvertVendorsToSearchResults(matchingVendors []Vendor, total int, count int) (*SearchResults, error) {

	var vendors []*VendorReturnJSON

	for _, vendor := range matchingVendors {
		vendors = append(vendors, vendor.ConvertToReturnJSON())
	}

	results, err := util.ConvertStructsToInterfaces(vendors)
	if err != nil {
		return nil, err
	}

	searchResults := &SearchResults{
		TotalCount:    total,
		FilteredCount: count,
		Data:          results,
	}

	return searchResults, nil

}

func (vlr *VendorListRequest) ConvertToOrganizationQuery(ctx context.Context) *gorm.DB {

	query := util.DBFromContext(ctx).Model(&Vendor{}).
		Select("DISTINCT vendors.*").
		Joins("LEFT JOIN clients ON clients.id = vendors.client_id").
		Joins("LEFT JOIN organizations ON organizations.id = clients.organization_id")

	query = query.Where("organizations.id = ?", vlr.OrganizationID)

	if vlr.ClientID != 0 {
		query = query.Where("vendors.client_id = ?", vlr.ClientID)
	}

	if vlr.SearchValue != "" {
		query = query.Where("vendors.name ILIKE ?", "%"+vlr.SearchValue+"%")
	}

	if vlr.OrderByColumn != "" {
		query = query.Order(vlr.OrderByColumn + " " + vlr.OrderBy)
	} else {
		query = query.Order("vendors.id " + vlr.OrderBy)
	}

	return query
}

func (vlr *VendorListRequest) ConvertToClientQuery(ctx context.Context) *gorm.DB {

	query := util.DBFromContext(ctx).Model(&Vendor{}).
		Select("DISTINCT vendors.*")

	query = query.Where("vendors.client_id = ?", vlr.ClientID)

	if vlr.SearchValue != "" {
		query = query.Where("vendors.name ILIKE ?", "%"+vlr.SearchValue+"%")
	}

	if vlr.OrderByColumn != "" {
		query = query.Order(vlr.OrderByColumn + " " + vlr.OrderBy)
	} else {
		query = query.Order("vendors.id " + vlr.OrderBy)
	}

	return query
}
