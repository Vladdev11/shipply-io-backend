package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type Location struct {
	ID             int
	WarehouseID    int
	Name           string
	LocationTypeID int
	Pickable       bool
	Sellable       bool
	IsTote         bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	LocationType LocationType
	Warehouse    Warehouse
}

type LocationReturnJSON struct {
	ID             int    `json:"id"`
	WarehouseID    int    `json:"warehouse_id"`
	Name           string `json:"name"`
	LocationTypeID int    `json:"location_type_id"`
	Pickable       bool   `json:"pickable"`
	Sellable       bool   `json:"sellable"`
	IsTote         bool   `json:"is_tote"`

	LocationType *LocationTypeReturnJSON `json:"location_type"`
}

type LocationListRequest struct {
	OrganizationID int    `json:"organization_id"`
	WarehouseID    int    `json:"warehouse_id"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	OrderBy        string `json:"order_by"`
	OrderByColumn  string `json:"order_by_column"`
	SearchValue    string `json:"search_value"`
	// TODO sellable, pickable, is_tote filters
}

type LocationCreateRequest struct {
	Name           string `json:"name"`
	WarehouseID    int    `json:"warehouse_id"`
	LocationTypeID int    `json:"location_type_id"`
	Pickable       bool   `json:"pickable"`
	Sellable       bool   `json:"sellable"`
	IsTote         bool   `json:"is_tote"`
}

type LocationUpdateRequest struct {
	Name           string `json:"name"`
	WarehouseID    int    `json:"warehouse_id"`
	LocationTypeID int    `json:"location_type_id"`
	Pickable       bool   `json:"pickable"`
	Sellable       bool   `json:"sellable"`
	IsTote         bool   `json:"is_tote"`
}

func (l *Location) Create() error {
	err := PGDB.Create(l).Error
	if err != nil {
		return err
	}

	return nil
}

func (l *Location) Delete() error {
	err := PGDB.Delete(l).Error
	if err != nil {
		return err
	}

	return nil
}

func (l *Location) HasInventory() bool {
	var count int64
	PGDB.Model(&Inventory{}).Where("location_id = ?", l.ID).Count(&count)
	if count > 0 {
		return true
	}
	return false
}

func (l *Location) GetLocationType() error {
	err := PGDB.Model(l).Association("LocationType").Find(&l.LocationType)
	if err != nil {
		return err
	}

	return nil
}

func (l *Location) UpdateWithRequest(request *LocationUpdateRequest) error {

	if request.Name != "" {
		l.Name = request.Name
	}

	if request.LocationTypeID != 0 {
		l.LocationTypeID = request.LocationTypeID
	}

	l.Pickable = request.Pickable
	l.Sellable = request.Sellable
	l.IsTote = request.IsTote

	err := PGDB.Save(l).Error
	if err != nil {
		return err
	}

	return nil
}

func (l *Location) ConvertToReturnJSON() *LocationReturnJSON {
	return &LocationReturnJSON{
		ID:             l.ID,
		WarehouseID:    l.WarehouseID,
		Name:           l.Name,
		LocationTypeID: l.LocationTypeID,
		Pickable:       l.Pickable,
		Sellable:       l.Sellable,
		IsTote:         l.IsTote,
		LocationType:   l.LocationType.ConvertToReturnJSON(),
	}
}

func (llr *LocationListRequest) ParseAndValidateRequest(r *http.Request) []string {

	errors := []string{}

	llr.WarehouseID = 0
	warehouseID, err := util.GetIntQueryParam(r, "warehouse_id")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "warehouse_id must be an integer")
		} else if warehouseID < 0 {
			errors = append(errors, "warehouse_id must be greater than or equal to 0")
		} else {
			llr.WarehouseID = warehouseID
		}
	}

	llr.Limit = 100
	limit, err := util.GetIntQueryParam(r, "limit")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "limit must be an integer")
		} else if limit < 0 {
			errors = append(errors, "limit must be greater than or equal to 0")
		} else {
			llr.Limit = limit
		}
	}

	llr.Offset = 0
	offset, err := util.GetIntQueryParam(r, "offset")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "offset must be an integer")
		} else if offset < 0 {
			errors = append(errors, "offset must be greater than or equal to 0")
		} else {
			llr.Offset = offset
		}
	}

	llr.OrderBy = "asc"
	orderBy, err := util.GetStringQueryParam(r, "order")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "order must be a string")
		} else if orderBy != "asc" && orderBy != "desc" {
			errors = append(errors, "order must be either asc or desc")
		} else {
			llr.OrderBy = orderBy
		}
	}

	searchValue, err := util.GetStringQueryParam(r, "search_value")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "search_value must be a string")
		} else {
			llr.SearchValue = searchValue
		}
	}

	orderByColumn, err := util.GetStringQueryParam(r, "order_by_column")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "order_by_column must be a string")
		} else {
			llr.OrderByColumn = orderByColumn
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (llr *LocationListRequest) ConvertToOrganizationQuery() *gorm.DB {

	query := PGDB.Model(&Location{}).
		Select("DISTINCT locations.*").
		Joins("LEFT JOIN warehouses ON warehouses.id = locations.warehouse_id")

	query = query.Where("warehouses.organization_id = ?", llr.OrganizationID)

	if llr.WarehouseID != 0 {
		query = query.Where("locations.warehouse_id = ?", llr.WarehouseID)
	}

	if llr.SearchValue != "" {
		query = query.Where("locations.name ILIKE ?", "%"+llr.SearchValue+"%")
	}

	if llr.OrderByColumn != "" {
		query = query.Order(llr.OrderByColumn + " " + llr.OrderBy)
	} else {
		query = query.Order("locations.id " + llr.OrderBy)
	}

	return query
}

func (lcr *LocationCreateRequest) ParseAndValidateRequest(r *http.Request) []string {

	errors := []string{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Name           json.RawMessage `json:"name"`
		WarehouseID    json.RawMessage `json:"warehouse_id"`
		LocationTypeID json.RawMessage `json:"location_type_id"`
		Pickable       json.RawMessage `json:"pickable"`
		Sellable       json.RawMessage `json:"sellable"`
		IsTote         json.RawMessage `json:"is_tote"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	user, err := GetRequestingUser(r)
	if err != nil {
		return []string{"invalid user"}
	}

	if aux.WarehouseID == nil {
		errors = append(errors, "warehouse_id is required")
	} else if err := json.Unmarshal(aux.WarehouseID, &lcr.WarehouseID); err != nil {
		errors = append(errors, "warehouse_id must be an integer")
	} else if lcr.WarehouseID < 0 {
		errors = append(errors, "warehouse_id must be greater than or equal to 0")
	} else {
		warehouse, err := GetWarehouseByID(lcr.WarehouseID)
		if err != nil {
			errors = append(errors, "warehouse_id must be a valid warehouse")
		} else if warehouse.OrganizationID != user.OwnerID {
			errors = append(errors, "user does not have access to that warehouse")
		}
	}

	if aux.Name == nil {
		errors = append(errors, "name is required")
	} else if err := json.Unmarshal(aux.Name, &lcr.Name); err != nil {
		errors = append(errors, "name must be a string")
	} else if len(lcr.Name) > 255 {
		errors = append(errors, "name must be less than 255 characters")
	} else {
		locations, err := GetLocationsByWarehouseID(lcr.WarehouseID)
		if err != nil {
			return []string{"error getting locations"}
		}
		for _, location := range locations {
			if strings.ToLower(location.Name) == strings.ToLower(lcr.Name) {
				errors = append(errors, "a location with that name already exists for that warehouse")
				break
			}
		}

	}

	if aux.LocationTypeID == nil {
		errors = append(errors, "location_type_id is required")
	} else if err := json.Unmarshal(aux.LocationTypeID, &lcr.LocationTypeID); err != nil {
		errors = append(errors, "location_type_id must be an integer")
	} else {
		locationType, err := GetLocationTypeByID(lcr.LocationTypeID)
		if err != nil {
			errors = append(errors, "location_type_id does not exist")
		} else {
			if locationType.OrganizationID != user.OwnerID {
				errors = append(errors, "location_type_id does not belong to your organization")
			}
		}
	}

	if aux.Pickable != nil {
		if err := json.Unmarshal(aux.Pickable, &lcr.Pickable); err != nil {
			errors = append(errors, "pickable must be a boolean")
		}
	} else {
		lcr.Pickable = false
	}

	if aux.Sellable != nil {
		if err := json.Unmarshal(aux.Sellable, &lcr.Sellable); err != nil {
			errors = append(errors, "sellable must be a boolean")
		}
	} else {
		lcr.Sellable = false
	}

	if aux.IsTote != nil {
		if err := json.Unmarshal(aux.IsTote, &lcr.IsTote); err != nil {
			errors = append(errors, "is_tote must be a boolean")
		}
	} else {
		lcr.IsTote = false
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (lur *LocationUpdateRequest) ParseAndValidateRequest(r *http.Request) []string {

	errors := []string{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Name           json.RawMessage `json:"name"`
		WarehouseID    json.RawMessage `json:"warehouse_id"`
		LocationTypeID json.RawMessage `json:"location_type_id"`
		Pickable       json.RawMessage `json:"pickable"`
		Sellable       json.RawMessage `json:"sellable"`
		IsTote         json.RawMessage `json:"is_tote"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	locationID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		errors = append(errors, "invalid location id")
	}

	location, err := GetLocationByID(locationID)
	if err != nil {
		return []string{"failed to get location"}
	}

	user, err := GetRequestingUser(r)
	if err != nil {
		return []string{"invalid user"}
	}

	if aux.WarehouseID != nil {
		if err := json.Unmarshal(aux.WarehouseID, &lur.WarehouseID); err != nil {
			errors = append(errors, "warehouse_id must be an integer")
		} else if lur.WarehouseID < 0 {
			errors = append(errors, "warehouse_id must be greater than or equal to 0")
		} else {
			warehouse, err := GetWarehouseByID(lur.WarehouseID)
			if err != nil {
				errors = append(errors, "warehouse_id must be a valid warehouse")
			} else if warehouse.OrganizationID != user.OwnerID {
				errors = append(errors, "user does not have access to that warehouse")
			}
		}
	} else {
		lur.WarehouseID = location.WarehouseID
	}

	if aux.Name != nil {
		if err := json.Unmarshal(aux.Name, &lur.Name); err != nil {
			errors = append(errors, "name must be a string")
		} else if len(lur.Name) > 255 {
			errors = append(errors, "name must be less than 255 characters")
		} else {
			locations, err := GetLocationsByWarehouseID(lur.WarehouseID)
			if err != nil {
				return []string{"error getting locations"}
			}
			for _, location := range locations {
				if strings.ToLower(location.Name) == strings.ToLower(lur.Name) && location.ID != locationID {
					errors = append(errors, "a location with that name already exists for that warehouse")
					break
				}
			}
		}
	}

	if aux.LocationTypeID != nil {
		if err := json.Unmarshal(aux.LocationTypeID, &lur.LocationTypeID); err != nil {
			errors = append(errors, "location_type_id must be an integer")
		} else {
			warehouse, err := GetWarehouseByID(location.WarehouseID)
			if err != nil {
				errors = append(errors, "failed to get warehouse")
			} else {
				if warehouse.OrganizationID != user.OwnerID {
					errors = append(errors, "user does not have access to that warehouse")
				}
			}

			locationType, err := GetLocationTypeByID(lur.LocationTypeID)
			if err != nil {
				errors = append(errors, "location_type_id does not exist")
			} else {
				if locationType.OrganizationID != warehouse.OrganizationID {
					errors = append(errors, "location_type_id does not belong to your organization")
				}
			}
		}
	}

	if aux.Pickable != nil {
		if err := json.Unmarshal(aux.Pickable, &lur.Pickable); err != nil {
			errors = append(errors, "pickable must be a boolean")
		}
	} else {
		lur.Pickable = location.Pickable
	}

	if aux.Sellable != nil {
		if err := json.Unmarshal(aux.Sellable, &lur.Sellable); err != nil {
			errors = append(errors, "sellable must be a boolean")
		}
	} else {
		lur.Sellable = location.Sellable
	}

	if aux.IsTote != nil {
		if err := json.Unmarshal(aux.IsTote, &lur.IsTote); err != nil {
			errors = append(errors, "is_tote must be a boolean")
		}
	} else {
		lur.IsTote = location.IsTote
	}

	if len(errors) > 0 {
		return errors
	}

	return nil

}

func GetLocationByID(locationID int) (Location, error) {
	var location Location
	err := PGDB.Where("id = ?", locationID).First(&location).Error
	return location, err
}

func LocationExists(locationID int) bool {
	var location Location
	err := PGDB.Where("id = ?", locationID).First(&location).Error
	if err != nil {
		return false
	}
	return true
}

func ConvertLocationsToSearchResults(matchingLocations []Location, total int, count int) (*SearchResults, error) {

	var locations []*LocationReturnJSON

	for _, location := range matchingLocations {
		err := location.GetLocationType()
		if err != nil {
			return nil, err
		}

		locations = append(locations, location.ConvertToReturnJSON())
	}

	results, err := util.ConvertStructsToInterfaces(locations)
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

func GetLocationsByWarehouseID(warehouseID int) ([]Location, error) {
	var locations []Location
	err := PGDB.Where("warehouse_id = ?", warehouseID).Find(&locations).Error
	return locations, err
}

func (l *Location) HasActivePickSessionOrder() bool {
	var pickSessionOrder PickSessionOrder
	if err := PGDB.Where("location_id = ? AND shipped IS NOT TRUE", l.ID).
		First(&pickSessionOrder).Error; err != nil {
		return false
	}
	return true
}

func (l *Location) GetActivePickSessionOrder() (*PickSessionOrder, error) {
	var pickSessionOrder PickSessionOrder
	if err := PGDB.Where("location_id = ? AND shipped IS NOT TRUE", l.ID).
		First(&pickSessionOrder).Error; err != nil {
		return nil, err
	}
	return &pickSessionOrder, nil
}

func (l *Location) GetWarehouse() error {
	var warehouse Warehouse
	err := PGDB.Where("id = ?", l.WarehouseID).First(&warehouse).Error
	if err != nil {
		return err
	}
	l.Warehouse = warehouse
	return nil
}
