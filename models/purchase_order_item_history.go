package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type PurchaseOrderItemHistory struct {
	ID                  int    `json:"id"`
	PurchaseOrderItemID int    `json:"purchase_order_item_id"`
	Note                string `json:"note"`
	CreatedBy           int    `json:"created_by"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	PurchaseOrderItem *PurchaseOrderItem `json:"purchase_order"`
	CreatedByUser     *User              `json:"created_by_user" gorm:"foreignKey:CreatedBy"`
}

type PurchaseOrderItemHistoryReturnJSON struct {
	Id                  int            `json:"id"`
	PurchaseOrderItemID int            `json:"purchase_order_item_id,omitempty"`
	Note                string         `json:"note"`
	CreatedBy           UserReturnJSON `json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (poih *PurchaseOrderItemHistory) Create(ctx context.Context) error {
	err := util.DBFromContext(ctx).Save(poih).Error
	if err != nil {
		return err
	}
	return nil
}

func GetPurchaseOrderItemHistorysByID(ctx context.Context, id int) ([]PurchaseOrderItemHistory, error) {

	var PurchaseOrderItemHistorys []PurchaseOrderItemHistory
	err := util.DBFromContext(ctx).Where("purchase_order_item_id = ?", id).Order("id desc").Find(&PurchaseOrderItemHistorys).Error
	if err != nil {
		return nil, err
	}

	return PurchaseOrderItemHistorys, nil
}

func (pon *PurchaseOrderItemHistory) ConvertToReturnJSON(ctx context.Context) PurchaseOrderItemHistoryReturnJSON {
	user, _ := GetUserByID(ctx, pon.CreatedBy)
	createdByUser := *user.ConvertToReturnJSON(ctx)

	return PurchaseOrderItemHistoryReturnJSON{
		Id:                  pon.ID,
		PurchaseOrderItemID: pon.PurchaseOrderItemID,
		Note:                pon.Note,
		CreatedBy:           createdByUser,
		CreatedAt:           pon.CreatedAt,
		UpdatedAt:           pon.UpdatedAt,
	}
}
