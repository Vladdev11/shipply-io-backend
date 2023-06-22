package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type Box struct {
	ID                    int
	Name                  string
	OrganizationID        int
	Barcode               string
	Length                float64
	Width                 float64
	Height                float64
	Weight                float64
	Type                  string
	Active                bool
	Cost                  float64
	ShipenginePackageID   string
	ShipenginePackageCode string
	CreatedBy             int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Organization Organization
}

type BoxReturnJSON struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Barcode string  `json:"barcode"`
	Length  float64 `json:"length"`
	Width   float64 `json:"width"`
	Height  float64 `json:"height"`
	Weight  float64 `json:"weight"`
	Type    string  `json:"type"`
	Active  bool    `json:"active"`
	Cost    float64 `json:"cost"`
}

type BoxCreateRequest struct {
	Name    string  `json:"name"`
	Barcode string  `json:"barcode"`
	Length  float64 `json:"length"`
	Width   float64 `json:"width"`
	Height  float64 `json:"height"`
	Weight  float64 `json:"weight"`
	Type    string  `json:"type"`
	Active  bool    `json:"active"`
	Cost    float64 `json:"cost"`
}

type BoxUpdateRequest struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Barcode string  `json:"barcode"`
	Length  float64 `json:"length"`
	Width   float64 `json:"width"`
	Height  float64 `json:"height"`
	Weight  float64 `json:"weight"`
	Type    string  `json:"type"`
	Active  bool    `json:"active"`
	Cost    float64 `json:"cost"`
}

func (b *Box) Create() error {
	err := PGDB.Create(b).Error
	if err != nil {
		return err
	}

	return nil
}

func (b *Box) Update() error {
	err := PGDB.Save(b).Error
	if err != nil {
		return err
	}

	return nil
}

func (b *Box) Delete() error {
	err := PGDB.Delete(b).Error
	if err != nil {
		return err
	}

	return nil
}

func GetBoxByBarcode(barcode string) (*Box, error) {
	box := &Box{}
	err := PGDB.Where("barcode = ?", barcode).First(box).Error
	if err != nil {
		return nil, err
	}

	return box, nil
}

func CheckForUniqueBarcode(barcode string, boxID int) (bool, error) {
	box := &Box{}
	err := PGDB.Where("barcode = ? AND id != ?", barcode, boxID).First(box).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

func (b *Box) ConvertToReturnJSON() *BoxReturnJSON {
	return &BoxReturnJSON{
		ID:      b.ID,
		Name:    b.Name,
		Barcode: b.Barcode,
		Length:  b.Length,
		Width:   b.Width,
		Height:  b.Height,
		Weight:  b.Weight,
		Type:    b.Type,
		Active:  b.Active,
		Cost:    b.Cost,
	}
}

func (b *BoxCreateRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errs []string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Name    json.RawMessage `json:"name"`
		Barcode json.RawMessage `json:"barcode"`
		Length  json.RawMessage `json:"length"`
		Width   json.RawMessage `json:"width"`
		Height  json.RawMessage `json:"height"`
		Weight  json.RawMessage `json:"weight"`
		Type    json.RawMessage `json:"type"`
		Active  json.RawMessage `json:"active"`
		Cost    json.RawMessage `json:"cost"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Name != nil {
		if err := json.Unmarshal(aux.Name, &b.Name); err != nil {
			errs = append(errs, "name must be a string")
		} else if len(b.Name) > 255 {
			errs = append(errs, "name must be less than 255 characters")
		} else if len(b.Name) == 0 {
			errs = append(errs, "name must be greater than 0 characters")
		}
	} else {
		errs = append(errs, "name is required")
	}

	if aux.Barcode != nil {
		if err := json.Unmarshal(aux.Barcode, &b.Barcode); err != nil {
			errs = append(errs, "barcode must be a string")
		} else if len(b.Barcode) > 255 {
			errs = append(errs, "barcode must be less than 255 characters")
		} else if len(b.Barcode) == 0 {
			errs = append(errs, "name must be greater than 0 characters")
		}

		//make sure barcode is unique
		box, err := GetBoxByBarcode(b.Barcode)
		if err != nil && err != gorm.ErrRecordNotFound {
			errs = append(errs, "error checking if barcode is unique")
		} else if box != nil {
			errs = append(errs, "barcode must be unique")
		}
	} else {
		errs = append(errs, "barcode is required")
	}

	if aux.Length != nil {
		if err := json.Unmarshal(aux.Length, &b.Length); err != nil {
			errs = append(errs, "length must be a float")
		} else if b.Length <= 0 {
			errs = append(errs, "length must be greater than 0")
		}
	} else {
		errs = append(errs, "length is required")
	}

	if aux.Width != nil {
		if err := json.Unmarshal(aux.Width, &b.Width); err != nil {
			errs = append(errs, "width must be a float")
		} else if b.Width <= 0 {
			errs = append(errs, "width must be greater than 0")
		}
	} else {
		errs = append(errs, "width is required")
	}

	if aux.Height != nil {
		if err := json.Unmarshal(aux.Height, &b.Height); err != nil {
			errs = append(errs, "height must be a float")
		} else if b.Height <= 0 {
			errs = append(errs, "height must be greater than 0")
		}
	} else {
		errs = append(errs, "height is required")
	}

	if aux.Weight != nil {
		if err := json.Unmarshal(aux.Weight, &b.Weight); err != nil {
			errs = append(errs, "weight must be a float")
		} else if b.Weight <= 0 {
			errs = append(errs, "weight must be greater than 0")
		}
	} else {
		errs = append(errs, "weight is required")
	}

	if aux.Type != nil {
		if err := json.Unmarshal(aux.Type, &b.Type); err != nil {
			errs = append(errs, "type must be a string")
		} else {
			if len(b.Type) > 255 {
				errs = append(errs, "type must be less than 255 characters")
			}
		}
	}

	if aux.Active != nil {
		if err := json.Unmarshal(aux.Active, &b.Active); err != nil {
			errs = append(errs, "active must be a boolean")
		}
	}

	if aux.Cost != nil {
		if err := json.Unmarshal(aux.Cost, &b.Cost); err != nil {
			errs = append(errs, "cost must be a float")
		} else if b.Cost <= 0 {
			errs = append(errs, "cost must be greater than 0")
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func GetBoxByID(id int) (*Box, error) {
	box := &Box{}
	err := PGDB.Where("id = ?", id).First(box).Error
	if err != nil {
		return nil, err
	}

	return box, nil
}

func (b *BoxUpdateRequest) ParseAndValidateRequest(r *http.Request) []string {
	var errs []string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Name    *json.RawMessage `json:"name"`
		Barcode *json.RawMessage `json:"barcode"`
		Length  *json.RawMessage `json:"length"`
		Width   *json.RawMessage `json:"width"`
		Height  *json.RawMessage `json:"height"`
		Weight  *json.RawMessage `json:"weight"`
		Type    *json.RawMessage `json:"type"`
		Active  *json.RawMessage `json:"active"`
		Cost    *json.RawMessage `json:"cost"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Name != nil {
		if err := json.Unmarshal(*aux.Name, &b.Name); err != nil {
			errs = append(errs, "name must be a string")
		} else if len(b.Name) > 255 {
			errs = append(errs, "name must be less than 255 characters")
		} else if len(b.Name) == 0 {
			errs = append(errs, "name must be greater than 0 characters")
		}
	}

	if aux.Barcode != nil {
		if err := json.Unmarshal(*aux.Barcode, &b.Barcode); err != nil {
			errs = append(errs, "barcode must be a string")
		} else if len(b.Barcode) > 255 {
			errs = append(errs, "barcode must be less than 255 characters")
		} else if len(b.Barcode) == 0 {
			errs = append(errs, "barcode must be greater than 0 characters")
		}

		// make sure barcode is unique
		exists, err := CheckForUniqueBarcode(b.Barcode, b.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			errs = append(errs, "error checking if barcode is unique")
		} else if !exists {
			errs = append(errs, "barcode must be unique")
		}
	}

	if aux.Length != nil {
		if err := json.Unmarshal(*aux.Length, &b.Length); err != nil {
			errs = append(errs, "length must be a float")
		} else if b.Length <= 0 {
			errs = append(errs, "length must be greater than 0")
		}
	}

	if aux.Width != nil {
		if err := json.Unmarshal(*aux.Width, &b.Width); err != nil {
			errs = append(errs, "width must be a float")
		} else if b.Width <= 0 {
			errs = append(errs, "width must be greater than 0")
		}
	}

	if aux.Height != nil {
		if err := json.Unmarshal(*aux.Height, &b.Height); err != nil {
			errs = append(errs, "height must be a float")
		} else if b.Height <= 0 {
			errs = append(errs, "height must be greater than 0")
		}
	}

	if aux.Weight != nil {
		if err := json.Unmarshal(*aux.Weight, &b.Weight); err != nil {
			errs = append(errs, "weight must be a float")
		} else if b.Weight <= 0 {
			errs = append(errs, "weight must be greater than 0")
		}
	}

	if aux.Type != nil {
		if err := json.Unmarshal(*aux.Type, &b.Type); err != nil {
			errs = append(errs, "type must be a string")
		} else {
			if len(b.Type) > 255 {
				errs = append(errs, "type must be less than 255 characters")
			}
		}
	}

	if aux.Active != nil {
		if err := json.Unmarshal(*aux.Active, &b.Active); err != nil {
			errs = append(errs, "active must be a boolean")
		}
	}

	if aux.Cost != nil {
		if err := json.Unmarshal(*aux.Cost, &b.Cost); err != nil {
			errs = append(errs, "cost must be a float")
		} else if b.Cost <= 0 {
			errs = append(errs, "cost must be greater than 0")
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (r *Box) UpdateWithRequest(request BoxUpdateRequest) {
	if request.Name != "" {
		r.Name = request.Name
	}
	if request.Barcode != "" {
		r.Barcode = request.Barcode
	}
	if request.Length != 0 {
		r.Length = request.Length
	}
	if request.Width != 0 {
		r.Width = request.Width
	}
	if request.Height != 0 {
		r.Height = request.Height
	}
	if request.Weight != 0 {
		r.Weight = request.Weight
	}
	if request.Type != "" {
		r.Type = request.Type
	}
	if request.Active != false {
		r.Active = true
	} else {
		r.Active = false
	}
	if request.Cost != 0 {
		r.Cost = request.Cost
	}
}
