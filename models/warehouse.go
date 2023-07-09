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

type Warehouse struct {
	ID                    int
	OrganizationID        int
	Name                  string
	ShipFromAddressID     int
	ReturnAddressID       int
	ShipengineWarehouseID string
	Default               bool

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	ShipFromAddress Address
	ReturnAddress   Address
	Organization    Organization
	Locations       []Location
}

type WarehouseListRequest struct {
	OrganizationID int    `json:"organization_id"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	OrderBy        string `json:"order_by"`
	OrderByColumn  string `json:"order_by_column"`
	SearchValue    string `json:"search_value"`
}

type WarehouseCreateRequest struct {
	Name            string   `json:"name"`
	ShipFromAddress *Address `json:"ship_from_address"`
	ReturnAddress   *Address `json:"return_address"`
}

type WarehouseUpdateRequest struct {
	Name            string  `json:"name"`
	ShipFromAddress Address `json:"ship_from_address"`
	ReturnAddress   Address `json:"return_address"`
}

type WarehouseReturnJSON struct {
	ID             int    `json:"id"`
	OrganizationID int    `json:"organization_id"`
	Name           string `json:"name"`
	UpdatedAt      string `json:"updated_at"`

	ShipFromAddress *AddressReturnJSON `json:"ship_from_address"`
	ReturnAddress   *AddressReturnJSON `json:"return_address"`
}

func (w *Warehouse) Create(ctx context.Context) error {

	err := util.DBFromContext(ctx).Create(w).Error
	if err != nil {
		return err
	}

	return nil
}

func (w *Warehouse) Delete(ctx context.Context) error {

	err := util.DBFromContext(ctx).Delete(w).Error
	if err != nil {
		return err
	}

	return nil
}

func (w *Warehouse) UpdateWithRequest(ctx context.Context, request *WarehouseUpdateRequest) error {

	if request.Name != "" {
		w.Name = request.Name
	}

	err := request.ShipFromAddress.Create(ctx)
	if err != nil {
		return err
	}

	err = request.ReturnAddress.Create(ctx)
	if err != nil {
		return err
	}

	if request.ShipFromAddress.ID != 0 {
		w.ShipFromAddressID = request.ShipFromAddress.ID
	}

	if err := util.DBFromContext(ctx).Save(w).Error; err != nil {
		return err
	}

	return nil
}

func (w *Warehouse) GetShipFromAddress(ctx context.Context) error {

	var address Address

	err := util.DBFromContext(ctx).Where("id = ?", w.ShipFromAddressID).First(&address).Error
	if err != nil {
		return err
	}

	w.ShipFromAddress = address

	return nil
}

func (w *Warehouse) GetReturnAddress(ctx context.Context) error {
	var address Address

	err := util.DBFromContext(ctx).Where("id = ?", w.ReturnAddressID).First(&address).Error
	if err != nil {
		return err
	}

	w.ReturnAddress = address

	return nil
}

func (w *Warehouse) ConvertToReturnJSON() *WarehouseReturnJSON {

	if w == nil {
		return &WarehouseReturnJSON{}
	}

	return &WarehouseReturnJSON{
		ID:              w.ID,
		OrganizationID:  w.OrganizationID,
		Name:            w.Name,
		ShipFromAddress: w.ShipFromAddress.ConvertToReturnJSON(),
		ReturnAddress:   w.ReturnAddress.ConvertToReturnJSON(),
		UpdatedAt:       w.UpdatedAt.Format(time.RFC3339),
	}
}

func (wlr *WarehouseListRequest) ParseAndValidateRequest(r *http.Request) []string {

	errors := []string{}

	wlr.Limit = 100
	limit, err := util.GetIntQueryParam(r, "limit")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "limit must be an integer")
		} else if limit < 0 {
			errors = append(errors, "limit must be greater than or equal to 0")
		} else {
			wlr.Limit = limit
		}
	}

	wlr.Offset = 0
	offset, err := util.GetIntQueryParam(r, "offset")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "offset must be an integer")
		} else if offset < 0 {
			errors = append(errors, "offset must be greater than or equal to 0")
		} else {
			wlr.Offset = offset
		}
	}

	wlr.OrderBy = "asc"
	orderBy, err := util.GetStringQueryParam(r, "order_by")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "order must be a string")
		} else if orderBy != "asc" && orderBy != "desc" {
			errors = append(errors, "order must be either asc or desc")
		} else {
			wlr.OrderBy = orderBy
		}
	}

	searchValue, err := util.GetStringQueryParam(r, "search_value")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "search_value must be a string")
		} else {
			wlr.SearchValue = searchValue
		}
	}

	orderByColumn, err := util.GetStringQueryParam(r, "order_by_column")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "order_by_column must be a string")
		} else {
			wlr.OrderByColumn = orderByColumn
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (wlr *WarehouseListRequest) ConvertToOrganizationQuery(ctx context.Context) *gorm.DB {

	query := util.DBFromContext(ctx).Model(&Warehouse{}).
		Select("DISTINCT warehouses.*").
		Where("warehouses.organization_id = ?", wlr.OrganizationID)

	if wlr.SearchValue != "" {
		query = query.Where("warehouses.name ILIKE ?", "%"+wlr.SearchValue+"%")
	}

	if wlr.OrderByColumn != "" {
		query = query.Order(wlr.OrderByColumn + " " + wlr.OrderBy)
	} else {
		query = query.Order("warehouses.id " + wlr.OrderBy)
	}

	return query
}

func (wcr *WarehouseCreateRequest) ParseAndValidateRequest(r *http.Request) []string {
	ctx := r.Context()

	var errs []string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Name            json.RawMessage `json:"name"`
		ShipFromAddress json.RawMessage `json:"ship_from_address"`
		ReturnAddress   json.RawMessage `json:"return_address"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Name == nil {
		errs = append(errs, "name is required")
	} else if err := json.Unmarshal(aux.Name, &wcr.Name); err != nil {
		errs = append(errs, "name must be a string")
	} else if len(wcr.Name) > 255 {
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
		warehouses, err := GetWarehousesByOrganizationID(ctx, user.Organization.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return []string{"invalid warehouses for organization"}
		}
		for _, warehouse := range warehouses {
			if warehouse.Name == wcr.Name {
				errs = append(errs, "warehouse with that name already exists")
			}
		}
	}

	if aux.ShipFromAddress == nil {
		errs = append(errs, "ship_from_address is required")
	} else {
		var address Address
		if err := json.Unmarshal(aux.ShipFromAddress, &address); err != nil {
			errs = append(errs, "ship_from_address is invalid")
		} else if err := address.Validate(); err != nil {
			errs = append(errs, err.Error())
		} else {
			wcr.ShipFromAddress = &address
		}
	}

	if aux.ReturnAddress == nil {
		errs = append(errs, "return_address is required")
	} else {
		var address Address
		if err := json.Unmarshal(aux.ReturnAddress, &address); err != nil {
			errs = append(errs, "return_address is invalid")
		} else if err := address.Validate(); err != nil {
			errs = append(errs, err.Error())
		} else {
			wcr.ReturnAddress = &address
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (wur *WarehouseUpdateRequest) ParseAndValidateRequest(r *http.Request) []string {
	ctx := r.Context()

	var errs []string

	warehouseID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		return []string{"invalid warehouse_id"}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Name            json.RawMessage `json:"name"`
		ShipFromAddress json.RawMessage `json:"ship_from_address"`
		ReturnAddress   json.RawMessage `json:"return_address"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if err := json.Unmarshal(aux.Name, &wur.Name); err != nil {
		errs = append(errs, "name must be a string")
	} else if len(wur.Name) > 255 {
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

		warehouses, err := GetWarehousesByOrganizationID(ctx, user.Organization.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return []string{"invalid warehouses for organization"}
		}
		for _, warehouse := range warehouses {
			if warehouse.Name == wur.Name && warehouse.ID != warehouseID {
				errs = append(errs, "warehouse with that name already exists")
			}
		}
	}

	if aux.ShipFromAddress != nil {
		var address Address
		if err := json.Unmarshal(aux.ShipFromAddress, &address); err != nil {
			errs = append(errs, "ship_from_address is invalid")
		} else if err := address.Validate(); err != nil {
			errs = append(errs, err.Error())
		} else {
			wur.ShipFromAddress = address
		}
	}

	if aux.ReturnAddress != nil {
		var address Address
		if err := json.Unmarshal(aux.ReturnAddress, &address); err != nil {
			errs = append(errs, "return_address is invalid")
		} else if err := address.Validate(); err != nil {
			errs = append(errs, err.Error())
		} else {
			wur.ReturnAddress = address
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func ConvertWarehousesToSearchResults(matchingWarehouses []Warehouse, total int, count int) (*SearchResults, error) {

	var warehouses []*WarehouseReturnJSON

	for _, warehouse := range matchingWarehouses {
		warehouses = append(warehouses, warehouse.ConvertToReturnJSON())
	}

	results, err := util.ConvertStructsToInterfaces(warehouses)
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

func GetWarehouseByID(ctx context.Context, id int) (*Warehouse, error) {
	warehouse := &Warehouse{}

	err := util.DBFromContext(ctx).First(warehouse, id).Error
	if err != nil {
		return nil, err
	}

	return warehouse, nil
}

func GetWarehousesByOrganizationID(ctx context.Context, organizationID int) ([]Warehouse, error) {
	warehouses := []Warehouse{}

	err := util.DBFromContext(ctx).Where("organization_id = ?", organizationID).Find(&warehouses).Error
	if err != nil {
		return nil, err
	}

	return warehouses, nil
}
