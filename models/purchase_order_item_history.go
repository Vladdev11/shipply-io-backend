package models

import (
	"time"

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

func (poih *PurchaseOrderItemHistory) Create() error {
	err := PGDB.Save(poih).Error
	if err != nil {
		return err
	}
	return nil
}

func GetPurchaseOrderItemHistorysByID(id int) ([]PurchaseOrderItemHistory, error) {

	var PurchaseOrderItemHistorys []PurchaseOrderItemHistory
	err := PGDB.Where("purchase_order_item_id = ?", id).Order("id desc").Find(&PurchaseOrderItemHistorys).Error
	if err != nil {
		return nil, err
	}

	return PurchaseOrderItemHistorys, nil
}

func (pon *PurchaseOrderItemHistory) ConvertToReturnJSON() PurchaseOrderItemHistoryReturnJSON {

	user, _ := GetUserByID(pon.CreatedBy)
	createdByUser := *user.ConvertToReturnJSON()

	return PurchaseOrderItemHistoryReturnJSON{
		Id:                  pon.ID,
		PurchaseOrderItemID: pon.PurchaseOrderItemID,
		Note:                pon.Note,
		CreatedBy:           createdByUser,
		CreatedAt:           pon.CreatedAt,
		UpdatedAt:           pon.UpdatedAt,
	}
}
