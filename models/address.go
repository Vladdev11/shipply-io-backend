package models

import (
	"errors"
	"fmt"
	"time"
)

type Address struct {
	ID              int       `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Company         string    `json:"company"`
	Street1         string    `json:"street1"`
	Street2         string    `json:"street2"`
	Street3         string    `json:"street3"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	PostalCode      string    `json:"postal_code"`
	Country         string    `json:"country"`
	Phone           string    `json:"phone"`
	Residential     bool      `json:"residential"`
	AddressVerified bool      `json:"address_verified"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type AddressReturnJSON struct {
	ID              int    `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Company         string `json:"company"`
	Street1         string `json:"street1"`
	Street2         string `json:"street2"`
	Street3         string `json:"street3"`
	City            string `json:"city"`
	State           string `json:"state"`
	PostalCode      string `json:"postal_code"`
	Country         string `json:"country"`
	Phone           string `json:"phone"`
	Residential     bool   `json:"residential"`
	AddressVerified string `json:"address_verified"`
}

func (a *Address) Validate() error {
	if a.Street1 == "" {
		return errors.New("street1 is required")
	}
	if a.City == "" {
		return errors.New("city is required")
	}
	if a.State == "" {
		return errors.New("state is required")
	}
	if a.PostalCode == "" {
		return errors.New("postal_code is required")
	}
	if a.Country == "" {
		return errors.New("country is required")
	}
	return nil
}

func (a *Address) Create() error {
	err := PGDB.Create(a).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *Address) ConvertToReturnJSON() *AddressReturnJSON {

	if a == nil {
		return &AddressReturnJSON{}
	}

	return &AddressReturnJSON{
		ID:              a.ID,
		FirstName:       a.FirstName,
		LastName:        a.LastName,
		Company:         a.Company,
		Street1:         a.Street1,
		Street2:         a.Street2,
		Street3:         a.Street3,
		City:            a.City,
		State:           a.State,
		PostalCode:      a.PostalCode,
		Country:         a.Country,
		Phone:           a.Phone,
		Residential:     a.Residential,
		AddressVerified: fmt.Sprintf("%t", a.AddressVerified),
	}
}
