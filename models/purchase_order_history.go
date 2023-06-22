package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type PurchaseOrderHistory struct {
	ID              int    `json:"id"`
	PurchaseOrderID int    `json:"purchase_order_id"`
	Note            string `json:"note"`
	CreatedBy       int    `json:"created_by"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	PurchaseOrder *PurchaseOrder `json:"purchase_order"`
	CreatedByUser *User          `json:"created_by_user" gorm:"foreignKey:CreatedBy"`
}

type PurchaseOrderHistoryReturnJSON struct {
	Id              int            `json:"id"`
	PurchaseOrderId int            `json:"purchase_order_id,omitempty"`
	Note            string         `json:"note"`
	CreatedBy       UserReturnJSON `json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PurchaseOrderHistoryCreateRequest struct {
	PurchaseOrderId int    `json:"purchase_order_id"`
	Note            string `json:"note"`
	CreatedBy       int    `json:"created_by"`
}

func (pon *PurchaseOrderHistoryCreateRequest) ConvertToPurchaseOrderHistory() *PurchaseOrderHistory {
	return &PurchaseOrderHistory{
		PurchaseOrderID: pon.PurchaseOrderId,
		Note:            pon.Note,
		CreatedBy:       pon.CreatedBy,
	}
}

func (pon *PurchaseOrderHistoryCreateRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errs []string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Note json.RawMessage `json:"note"`
	}{}

	if err := json.Unmarshal(body, &aux); err != nil {
		return []string{"invalid JSON"}
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		errs = append(errs, "invalid purchase order id")
	} else if purchaseOrderID <= 0 {
		errs = append(errs, "purchase order id must be greater than 0")
	} else {
		_, err := GetPurchaseOrderByID(purchaseOrderID)
		if err != nil {
			errs = append(errs, "purchase order id is not a valid purchase order")
		} else {
			pon.PurchaseOrderId = purchaseOrderID
		}
	}

	if aux.Note == nil {
		errs = append(errs, "note is required")
	} else if err := json.Unmarshal(aux.Note, &pon.Note); err != nil {
		errs = append(errs, "note must be a string")
	} else if len(pon.Note) == 0 {
		errs = append(errs, "note must be greater than 0 characters")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func (pon *PurchaseOrderHistory) Create() error {
	err := PGDB.Save(pon).Error
	if err != nil {
		return err
	}
	return nil
}

func GetPurchaseOrderHistorysByID(id int) ([]PurchaseOrderHistory, error) {

	var PurchaseOrderHistorys []PurchaseOrderHistory
	err := PGDB.Where("purchase_order_id = ?", id).Order("id desc").Find(&PurchaseOrderHistorys).Error
	if err != nil {
		return nil, err
	}

	return PurchaseOrderHistorys, nil
}

func (pon *PurchaseOrderHistory) ConvertToReturnJSON() PurchaseOrderHistoryReturnJSON {

	user, _ := GetUserByID(pon.CreatedBy)
	createdByUser := *user.ConvertToReturnJSON()

	return PurchaseOrderHistoryReturnJSON{
		Id:              pon.ID,
		PurchaseOrderId: pon.PurchaseOrderID,
		Note:            pon.Note,
		CreatedBy:       createdByUser,
		CreatedAt:       pon.CreatedAt,
		UpdatedAt:       pon.UpdatedAt,
	}
}
