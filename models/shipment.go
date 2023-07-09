package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type Shipment struct {
	ID                  int       `json:"id" gorm:"primary_key"`
	Voided              bool      `json:"voided"`
	VoidedAt            time.Time `json:"voided_at"`
	QuotedCost          float64   `json:"quoted_cost"`
	ShipmentCost        float64   `json:"shipment_cost"`
	InsuranceCost       float64   `json:"insurance_cost"`
	TrackingNumber      string    `json:"tracking_number"`
	IsReturnLabel       bool      `json:"is_return_label"`
	MarketplaceNotified bool      `json:"marketplace_notified"`
	ShippingRateID      int       `json:"shipping_rate_id"`
	CreatedByUserID     int       `json:"created_by_user_id"`
	LabelPDFURL         string    `json:"label_pdf_url"`
	PackageCode         string    `json:"package_code"`
	PickSessionOrderID  int       `json:"pick_session_order_id"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// BatchID int
	// Batch   Batch

	PickSessionOrder PickSessionOrder `json:"pick_session_order"`
}

func (s *Shipment) Create(ctx context.Context) error {
	return util.DBFromContext(ctx).Create(s).Error
}
