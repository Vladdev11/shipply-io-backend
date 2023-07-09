package models

import (
	"context"

	"github.com/shipply-io/shipply-io-backend/util"
)

type ShippingRate struct {
	ID                        int     `json:"id" gorm:"primary_key"`
	CarrierConnectionID       int     `json:"carrier_connection_id,omitempty"`
	CarrierConnectionNickName string  `json:"carrier_name,omitempty"`
	ShipengineRateID          string  `json:"shipengine_rate_id,omitempty"`
	ShipengineServiceType     string  `json:"shipengine_service_type,omitempty"`
	ShipengineServiceCode     string  `json:"shipengine_service_code,omitempty"`
	Amount                    float64 `json:"amount,omitempty"`
	Currency                  string  `json:"currency,omitempty"`
	BoxID                     int     `json:"box_id,omitempty"`
	Weight                    float64 `json:"weight,omitempty"`
	PickSessionOrderID        int     `json:"pick_session_order_id,omitempty"`
	DeliveryDays              int     `json:"delivery_days,omitempty"`
}

func (sr *ShippingRate) Create(ctx context.Context) error {
	return util.DBFromContext(ctx).Create(sr).Error
}

func GetShippingRateByID(ctx context.Context, id int) (ShippingRate, error) {
	var shippingRate ShippingRate
	err := util.DBFromContext(ctx).Where("id = ?", id).First(&shippingRate).Error
	return shippingRate, err
}
