package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

var (
	AllowedFilterTables = map[string]struct{}{
		"products": {},
	}
)

// UserSavedFilter is a user's saved filter.
// Table represents the table that the filter is for.
// Params is a JSON-encoded representation of the filter.
// The filter is quaried by a combination of the id, table and user_id.
type UserSavedFilter struct {
	gorm.Model
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Table  string `json:"table" gorm:"column:tb"`
	Params string `json:"params"`
}

type UserSavedFilterCreateRequest struct {
	Name   string         `json:"name"`
	Params map[string]any `json:"params"`
}

func (usfcr *UserSavedFilterCreateRequest) Execute(ctx context.Context, user int, table string) (*UserSavedFilter, error) {
	params, err := json.Marshal(usfcr.Params)
	if err != nil {
		return nil, fmt.Errorf("can't encode filter params: %w", err)
	}

	if usfcr.Name == "" {
		// TODO: Do we even want to give filters names?
		// If so what should the default be?
		usfcr.Name = "Untitled Filter"
	}

	ua := &UserSavedFilter{
		UserID: user,
		Name:   usfcr.Name,
		Table:  table,
		Params: string(params),
	}
	return ua, util.DBFromContext(ctx).Create(ua).Error
}

// Update updates a user's saved filter.
// You can update the name and/or params.
type UserSavedFilterUpdateRequest struct {
	Name   *string         `json:"name"`
	Params *map[string]any `json:"params"`
}

func (usfur *UserSavedFilterUpdateRequest) Execute(ctx context.Context, user int, table string, filterID int) (*UserSavedFilter, error) {
	db := util.DBFromContext(ctx)
	var usf UserSavedFilter
	if err := db.Where("user_id = ? AND tb = ? AND id = ?", user, table, filterID).First(&usf).Error; err != nil {
		return nil, err
	}

	if usfur.Name != nil {
		usf.Name = *usfur.Name
	}
	if usfur.Params != nil {
		params, err := json.Marshal(usfur.Params)
		if err != nil {
			return nil, fmt.Errorf("can't encode filter params: %w", err)
		}
		usf.Params = string(params)
	}

	return &usf, db.Save(&usf).Error
}

// LoadUserSavedFilter loads a user's saved filter into the target.
// Table is required to be one of the following:
// - products
// Target must be a pointer to one of the following:
// - ProductListFilters
func LoadUserSavedFilter(ctx context.Context, userID int, table string, filterID int, target any) error {
	var usf UserSavedFilter
	if err := util.DBFromContext(ctx).Where("user_id = ? AND tb = ? AND id = ?", userID, table, filterID).First(&usf).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("%s filter %d not found", table, filterID)
		}
		return err
	}

	switch t := target.(type) {
	case *ProductListFilters:
		var plf ProductListFilters
		if err := json.Unmarshal([]byte(usf.Params), &plf); err != nil {
			return fmt.Errorf("can't decode filter params: %w", err)
		}
		*t = plf
	default:
		return fmt.Errorf("invalid filter type: %T", target)
	}

	return nil
}

type UserSavedFilterJSON struct {
	ID     int            `json:"id"`
	Name   string         `json:"name"`
	Params map[string]any `json:"params"`
}

func (usf *UserSavedFilter) AsJSON() (*UserSavedFilterJSON, error) {
	var params map[string]any
	if err := json.Unmarshal([]byte(usf.Params), &params); err != nil {
		return nil, fmt.Errorf("can't decode filter params: %w", err)
	}
	return &UserSavedFilterJSON{
		ID:     int(usf.ID),
		Name:   usf.Name,
		Params: params,
	}, nil
}

func GetUserSavedFiltersForTable(ctx context.Context, userID int, table string) ([]UserSavedFilterJSON, error) {
	var usfs []UserSavedFilter
	if err := util.DBFromContext(ctx).Where("user_id = ? AND tb = ?", userID, table).Find(&usfs).Error; err != nil {
		return nil, err
	}
	usfsJSON := make([]UserSavedFilterJSON, 0, len(usfs))
	for i, usf := range usfs {
		usfJSON, err := usf.AsJSON()
		if err != nil {
			return nil, fmt.Errorf("can't decode filter params for filter #%d: %w", i, err)
		}
		usfsJSON = append(usfsJSON, *usfJSON)
	}
	return usfsJSON, nil
}

func GetUserSavedFilter(ctx context.Context, userID int, table string, filterID int) (*UserSavedFilterJSON, error) {
	var usf UserSavedFilter
	if err := util.DBFromContext(ctx).Where("user_id = ? AND tb = ? AND id = ?", userID, table, filterID).First(&usf).Error; err != nil {
		return nil, err
	}
	return usf.AsJSON()
}

func DeleteUserSavedFilter(ctx context.Context, userID int, table string, filterID int) error {
	db := util.DBFromContext(ctx)
	var usf UserSavedFilter
	if err := db.Where("user_id = ? AND tb = ? AND id = ?", userID, table, filterID).First(&usf).Error; err != nil {
		return err
	}
	return db.Delete(usf).Error
}
