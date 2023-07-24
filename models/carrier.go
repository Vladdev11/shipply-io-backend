package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type Carrier struct {
	ID                int
	ShipEngineID      string
	Name              string
	ThumbnailURL      string
	SmallThumbnailURL string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Organization *Organization `gorm:"-"`
	Client       *Client       `gorm:"-"`

	RequiredFields []util.FieldInfo `gorm:"-"`
}

type CarrierReturnJSON struct {
	ID                int                  `json:"id"`
	Name              string               `json:"name"`
	ThumbnailURL      string               `json:"thumbnail_url"`
	SmallThumbnailURL string               `json:"small_thumbnail_url"`
	RequiredFields    []util.FieldInfoJSON `json:"required_fields,omitempty"`
}

func (c *Carrier) ConvertToReturnJSON() CarrierReturnJSON {
	var requiredFields []util.FieldInfoJSON
	if c.RequiredFields != nil {
		for _, field := range c.RequiredFields {
			requiredFields = append(requiredFields, field.ConvertToReturnJSON())
		}
	}

	return CarrierReturnJSON{
		ID:                c.ID,
		Name:              c.Name,
		ThumbnailURL:      c.ThumbnailURL,
		SmallThumbnailURL: c.SmallThumbnailURL,
		RequiredFields:    requiredFields,
	}
}

func GetCarriers(ctx context.Context) ([]Carrier, error) {
	var carriers []Carrier
	err := util.DBFromContext(ctx).Find(&carriers).Error
	if err != nil {
		return nil, ErrQueryFailed{Err: err, Object: "carriers"}

	}

	return carriers, nil
}

func GetCarrierByID(ctx context.Context, id int) (Carrier, error) {
	var carrier Carrier
	err := util.DBFromContext(ctx).First(&carrier, id).Error
	if err != nil {
		return Carrier{}, ErrQueryFailed{Err: err, Object: "carrier"}
	}

	return carrier, nil
}

func (carrier *Carrier) Create(ctx context.Context) error {
	return util.DBFromContext(ctx).Create(&carrier).Error
}

func GetCarrierRequiredFields(carriers []Carrier) []Carrier {
	for i, carrier := range carriers {
		var fieldInfos []util.FieldInfo
		switch carrier.ShipEngineID {
		case "asendia":
			fieldInfos = util.GetFieldTypes(CarrierConnectAsendia{})
		case "dhl_ecommerce":
			fieldInfos = util.GetFieldTypes(CarrierConnectDHLeCommerce{})
		case "dhl_express":
			fieldInfos = util.GetFieldTypes(CarrierConnectDHLExpress{})
		case "ups":
			fieldInfos = util.GetFieldTypes(CarrierConnectUPSFrontend{})
		case "stamps_com":
			fieldInfos = util.GetFieldTypes(CarrierConnectStampsCom{})
		case "endicia":
			fieldInfos = util.GetFieldTypes(CarrierConnectEndicia{})
		case "ontrac":
			fieldInfos = util.GetFieldTypes(CarrierConnectOnTrac{})
		}
		carriers[i].RequiredFields = fieldInfos
	}

	return carriers
}
