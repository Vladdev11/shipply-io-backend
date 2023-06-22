package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type ProductLot struct {
	ID         int
	LotNumber  string
	ProductID  int
	ExpiryDate time.Time
	Active     bool
	CreatedBy  int

	Product       Product
	CreatedByUser *User `gorm:"foreignkey:CreatedBy"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func (pl *ProductLot) Create() error {
	return PGDB.Create(pl).Error
}

func (pl *ProductLot) Update() error {
	return PGDB.Save(pl).Error
}

type ProductLotReturnJSON struct {
	ID         int            `json:"id"`
	LotNumber  string         `json:"lot_number"`
	ProductID  int            `json:"product_id"`
	ExpiryDate string         `json:"expiry_date"`
	Active     bool           `json:"active"`
	CreatedBy  UserReturnJSON `json:"created_by"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (pl *ProductLot) ConvertToReturnJSON() *ProductLotReturnJSON {

	createdByUser, err := GetUserByID(pl.CreatedBy)
	if err != nil {
		return nil
	}

	return &ProductLotReturnJSON{
		ID:         pl.ID,
		LotNumber:  pl.LotNumber,
		ProductID:  pl.ProductID,
		ExpiryDate: pl.ExpiryDate.Format("2006-01-02"),
		Active:     pl.Active,
		CreatedBy:  *createdByUser.ConvertToReturnJSON(),
		CreatedAt:  pl.CreatedAt,
		UpdatedAt:  pl.UpdatedAt,
	}
}

type ProductLotCreateRequest struct {
	LotNumber  string     `json:"lot_number"`
	ProductID  int        `json:"product_id"`
	ExpiryDate SingleDate `json:"expiry_date"`
}

func (plcr *ProductLotCreateRequest) ParseAndValidateRequest(r *http.Request) []string {
	var errs []string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := struct {
		LotNumber  json.RawMessage `json:"lot_number"`
		ProductID  json.RawMessage `json:"product_id"`
		ExpiryDate json.RawMessage `json:"expiry_date"`
	}{}

	if err := json.Unmarshal(body, &aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.LotNumber != nil {
		if err := json.Unmarshal(aux.LotNumber, &plcr.LotNumber); err != nil {
			errs = append(errs, "lot_number must be a string")
		} else if len(plcr.LotNumber) == 0 {
			errs = append(errs, "lot_number cannot be empty")
		} else if len(plcr.LotNumber) > 255 {
			errs = append(errs, "lot_number cannot be longer than 255 characters")
		}
	} else {
		errs = append(errs, "lot_number is required")
	}

	if aux.ProductID != nil {
		if err := json.Unmarshal(aux.ProductID, &plcr.ProductID); err != nil {
			errs = append(errs, "product_id must be an integer")
		} else if plcr.ProductID <= 0 {
			errs = append(errs, "product_id must be a positive integer")
		}
	} else {
		errs = append(errs, "product_id is required")
	}

	if aux.ExpiryDate != nil {
		if err := json.Unmarshal(aux.ExpiryDate, &plcr.ExpiryDate); err != nil {
			errs = append(errs, "expiry_date must be a date in the format YYYY-MM-DD")
		}
	} else {
		errs = append(errs, "expiry_date is required")
	}

	return errs

}

func (pl *ProductLot) DetermineActive() {
	now := time.Now()
	if now.After(pl.ExpiryDate) {
		pl.Active = false
	} else {
		pl.Active = true
	}
}
