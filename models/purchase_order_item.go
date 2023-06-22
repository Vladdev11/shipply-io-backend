package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type PurchaseOrderItem struct {
	ID              int
	PurchaseOrderID int
	ProductID       int
	Ordered         int
	Received        int
	Damaged         int
	UnitPrice       float64
	Notes           string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	TotalPrice              float64 `gorm:"-"`
	Product                 *Product
	PurchaseOrder           *PurchaseOrder
	InventoryLocationLevels []InventoryLocationLevel     `gorm:"-"`
	History                 []PurchaseOrderItemHistory   `gorm:"-"`
	Rejections              []PurchaseOrderItemRejection `gorm:"-"`
	ProductLots             []ProductLot                 `gorm:"-"`
}

type InventoryLocationLevel struct {
	LocationID   int    `json:"location_id"`
	LocationName string `json:"location_name"`
	Count        int    `json:"count"`
}

type PurchaseOrderItemReturnJSON struct {
	ID                      int                                    `json:"id"`
	PurchaseOrderID         int                                    `json:"purchase_order_id,omitempty"`
	ProductID               int                                    `json:"product_id,omitempty"`
	Ordered                 int                                    `json:"ordered"`
	Received                int                                    `json:"received"`
	Damaged                 int                                    `json:"damaged"`
	UnitPrice               float64                                `json:"unit_price"`
	Notes                   string                                 `json:"notes"`
	Product                 *ProductReturnJSON                     `json:"product"`
	InventoryLocationLevels []InventoryLocationLevel               `json:"inventory_location_levels,omitempty"`
	History                 []PurchaseOrderItemHistoryReturnJSON   `json:"history,omitempty"`
	Rejections              []PurchaseOrderItemRejectionReturnJSON `json:"rejections,omitempty"`
	ProductLots             []ProductLotReturnJSON                 `json:"product_lots,omitempty"`
}

type PurchaseOrderItemCreateRequest struct {
	ProductID int     `json:"product_id"`
	Ordered   int     `json:"ordered"`
	UnitPrice float64 `json:"unit_price"`
}

type PurchaseOrderItemUpdateRequest struct {
	ID        int     `json:"id"`
	Ordered   int     `json:"ordered"`
	UnitPrice float64 `json:"unit_price"`
	Notes     string  `json:"notes"`
}

type PurchaseOrderItemReceiveRequest struct {
	Quantity            int `json:"quantity"`
	LocationID          int `json:"location_id"`
	PurchaseOrderItemID int `json:"purchase_order_item_id"`
	ProductLotID        int `json:"product_lot_id"`
}

type PurchaseOrderItemReceiveBatchRequest []PurchaseOrderItemReceiveRequest

type PurchaseOrderItemBulkUpdateRequest []PurchaseOrderItemUpdateRequest

func (poi *PurchaseOrderItem) Create() error {
	err := PGDB.Create(poi).Error
	if err != nil {
		return err
	}

	return nil
}

func (poi *PurchaseOrderItem) Delete() error {

	err := PGDB.Delete(poi).Error
	if err != nil {
		return err
	}

	return nil

}

func (poi *PurchaseOrderItem) GetProduct() error {
	product, err := GetProductByID(poi.ProductID)
	if err != nil {
		return err
	}

	poi.Product = product

	return nil
}

func GetPurchaseOrderItemByID(purchaseOrderItemID int) (*PurchaseOrderItem, error) {
	//define purchase order item
	purchaseOrderItem := &PurchaseOrderItem{}

	//get purchase order item from database
	err := PGDB.Where("id = ?", purchaseOrderItemID).First(purchaseOrderItem).Error
	if err != nil {
		return nil, err
	}

	//return purchase order item
	return purchaseOrderItem, nil
}

func UpdatePurchaseOrderItem(purchaseOrderItem *PurchaseOrderItem) error {
	//update purchase order item in database
	err := PGDB.Save(purchaseOrderItem).Error
	if err != nil {
		return err
	}

	return nil
}

func (poicr *PurchaseOrderItemCreateRequest) ParseAndValidateRequest(r *http.Request) []string {
	var errs []string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := struct {
		ProductID json.RawMessage `json:"product_id"`
		Ordered   json.RawMessage `json:"ordered"`
		UnitPrice json.RawMessage `json:"unit_price"`
	}{}

	if err := json.Unmarshal(body, &aux); err != nil {
		return []string{"invalid JSON"}
	}

	_, err = util.GetIntFromPath(r, "id")
	if err != nil {
		errs = append(errs, "invalid purchase order id")
	}

	if aux.ProductID == nil {
		errs = append(errs, "product_id is required")
	} else if err := json.Unmarshal(aux.ProductID, &poicr.ProductID); err != nil {
		errs = append(errs, "product_id must be of type int")
	} else if poicr.ProductID <= 0 {
		errs = append(errs, "product_id must be greater than 0")
	}

	if aux.Ordered != nil {
		if err := json.Unmarshal(aux.Ordered, &poicr.Ordered); err != nil {
			errs = append(errs, "ordered must be of type int")
		}
		if poicr.Ordered < 1 {
			errs = append(errs, "ordered must be greater than 0")
		}
	}

	if aux.UnitPrice != nil {
		if err := json.Unmarshal(aux.UnitPrice, &poicr.UnitPrice); err != nil {
			errs = append(errs, "unit_price must be of type float64")
		} else if aux.UnitPrice != nil && poicr.UnitPrice < 1 {
			errs = append(errs, "unit_price must be greater than 0")
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (p *PurchaseOrderItemUpdateRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errs []string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Ordered   json.RawMessage `json:"ordered"`
		UnitPrice json.RawMessage `json:"unit_price"`
		Notes     json.RawMessage `json:"notes"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Ordered != nil {
		if err := json.Unmarshal(aux.Ordered, &p.Ordered); err != nil {
			errs = append(errs, "ordered must be of type int")
		}
	}

	if aux.UnitPrice != nil {
		if err := json.Unmarshal(aux.UnitPrice, &p.UnitPrice); err != nil {
			errs = append(errs, "unit price must be of type float64")
		}
	}

	if aux.Notes != nil {
		if err := json.Unmarshal(aux.Notes, &p.Notes); err != nil {
			errs = append(errs, "notes must be of type string")
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func ParseAndValidatePurchaseOrderItemUpdateRequests(r *http.Request) ([]PurchaseOrderItemUpdateRequest, []string) {

	var errs []string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, []string{"invalid JSON"}
	}

	poiur := []PurchaseOrderItemUpdateRequest{}

	if err := json.Unmarshal(body, &poiur); err != nil {
		return nil, []string{"invalid JSON"}
	}

	for _, updateRequest := range poiur {

		updateRequestBytes, err := json.Marshal(updateRequest)
		if err != nil {
			errs = append(errs, "invalid JSON")
			continue
		}

		aux := &struct {
			ID        json.RawMessage `json:"id"`
			Ordered   json.RawMessage `json:"ordered"`
			UnitPrice json.RawMessage `json:"unit_price"`
			Notes     json.RawMessage `json:"notes"`
		}{}

		if err := json.Unmarshal(updateRequestBytes, aux); err != nil {
			errs = append(errs, "invalid JSON")
			continue
		}

		if err := json.Unmarshal(aux.ID, &updateRequest.ID); err != nil {
			errs = append(errs, "id must be of type int")
		} else if aux.ID != nil && updateRequest.ID < 1 {
			errs = append(errs, "id must be greater than 0")
		} else if aux.ID == nil {
			errs = append(errs, "id is required")
		}

		if err := json.Unmarshal(aux.Ordered, &updateRequest.Ordered); err != nil {
			errs = append(errs, "ordered must be of type int")
		}

		if err := json.Unmarshal(aux.UnitPrice, &updateRequest.UnitPrice); err != nil {
			errs = append(errs, "unit price must be of type float64")
		}

		if err := json.Unmarshal(aux.Notes, &updateRequest.Notes); err != nil {
			errs = append(errs, "notes must be of type string")
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return poiur, nil

}

func (p *PurchaseOrderItem) ConvertToReturnJSON() PurchaseOrderItemReturnJSON {

	if p == nil {
		return PurchaseOrderItemReturnJSON{}
	}

	PurchaseOrderItemReturnJSON := PurchaseOrderItemReturnJSON{
		ID:              p.ID,
		PurchaseOrderID: p.PurchaseOrderID,
		ProductID:       p.ProductID,
		Ordered:         p.Ordered,
		Received:        p.Received,
		Damaged:         p.Damaged,
		UnitPrice:       p.UnitPrice,
		Notes:           p.Notes,
		Product:         p.Product.ConvertToReturnJSON(),
	}

	if len(p.InventoryLocationLevels) > 0 {
		PurchaseOrderItemReturnJSON.InventoryLocationLevels = p.InventoryLocationLevels
	}

	historyReturnJSON := []PurchaseOrderItemHistoryReturnJSON{}
	if len(p.History) > 0 {
		for _, history := range p.History {
			historyReturnJSON = append(historyReturnJSON, history.ConvertToReturnJSON())
		}
		PurchaseOrderItemReturnJSON.History = historyReturnJSON
	}

	rejectionReturnJSON := []PurchaseOrderItemRejectionReturnJSON{}
	if len(p.Rejections) > 0 {
		for _, rejection := range p.Rejections {
			rejectionReturnJSON = append(rejectionReturnJSON, *rejection.ConvertToReturnJSON())
		}
		PurchaseOrderItemReturnJSON.Rejections = rejectionReturnJSON
	}

	ProductLotReturnJSON := []ProductLotReturnJSON{}
	if len(p.ProductLots) > 0 {
		for _, productLot := range p.ProductLots {
			ProductLotReturnJSON = append(ProductLotReturnJSON, *productLot.ConvertToReturnJSON())
		}
		PurchaseOrderItemReturnJSON.ProductLots = ProductLotReturnJSON
	}

	return PurchaseOrderItemReturnJSON
}

func DeletePurchaseOrderItem(purchaseOrderItem *PurchaseOrderItem) error {
	//delete purchase order item in database
	err := PGDB.Delete(purchaseOrderItem).Error
	if err != nil {
		return err
	}

	return nil
}

func DeletePurchaeOrderItemByPurchaseOrderID(purchaseOrderID int) error {
	//delete all purchase order items from database where purchase order id matches
	err := PGDB.Where("purchase_order_id = ?", purchaseOrderID).Delete(&PurchaseOrderItem{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *PurchaseOrderItem) UpdateWithRequest(request *PurchaseOrderItemUpdateRequest) error {

	if request.Ordered != 0 {
		p.Ordered = request.Ordered
	}

	if request.UnitPrice != 0 {
		p.UnitPrice = request.UnitPrice
	}

	if request.Notes != "" {
		p.Notes = request.Notes
	}

	if err := PGDB.Save(p).Error; err != nil {
		return err
	}

	return nil
}

func ValidateBatchPurchaseOrderItemReceiveRequest(r *http.Request) ([]PurchaseOrderItemReceiveRequest, []string) {

	var errs []string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, []string{"invalid JSON"}
	}

	aux := &[]struct {
		Quantity            json.RawMessage `json:"quantity"`
		LocationID          json.RawMessage `json:"location_id"`
		PurchaseOrderItemID json.RawMessage `json:"purchase_order_item_id"`
		ProductLotID        json.RawMessage `json:"product_lot_id"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return nil, []string{"invalid JSON"}
	}

	batchReceiveRequests := []PurchaseOrderItemReceiveRequest{}

	for _, batchReceiveRequest := range *aux {

		p := PurchaseOrderItemReceiveRequest{}

		if batchReceiveRequest.Quantity == nil {
			errs = append(errs, "quantity is required")
		} else if err := json.Unmarshal(batchReceiveRequest.Quantity, &p.Quantity); err != nil {
			errs = append(errs, "quantity must be of type int")
		} else if p.Quantity == 0 {
			errs = append(errs, "quantity must be a positive or negative number")
		}

		if batchReceiveRequest.LocationID == nil {
			errs = append(errs, "location_id is required")
		} else if err := json.Unmarshal(batchReceiveRequest.LocationID, &p.LocationID); err != nil {
			errs = append(errs, "location_id must be of type int")
		} else if p.LocationID < 1 {
			errs = append(errs, "location_id must be greater than 0")
		}

		if batchReceiveRequest.PurchaseOrderItemID == nil {
			errs = append(errs, "purchase_order_item_id is required")
		} else if err := json.Unmarshal(batchReceiveRequest.PurchaseOrderItemID, &p.PurchaseOrderItemID); err != nil {
			errs = append(errs, "purchase_order_item_id must be of type int")
		} else if p.PurchaseOrderItemID < 1 {
			errs = append(errs, "purchase_order_item_id must be greater than 0")
		}

		if batchReceiveRequest.ProductLotID != nil {
			if err := json.Unmarshal(batchReceiveRequest.ProductLotID, &p.ProductLotID); err != nil {
				errs = append(errs, "product_lot_id must be of type int")
			} else if p.ProductLotID < 1 {
				errs = append(errs, "product_lot_id must be greater than 0")
			}
		}

		batchReceiveRequests = append(batchReceiveRequests, p)

		if len(errs) > 0 {
			return nil, errs
		}

	}

	return batchReceiveRequests, nil

}

func (poi *PurchaseOrderItem) Update() error {
	if err := PGDB.Save(poi).Error; err != nil {
		return err
	}

	return nil
}

func (poi *PurchaseOrderItem) GetLocationsAndLevels() error {

	inventoryLocationLevels, err := GetProductLocationsAndLevelsByProductID(poi.ProductID)
	if err != nil {
		return err
	}

	poi.InventoryLocationLevels = inventoryLocationLevels

	return nil

}

func (poi *PurchaseOrderItem) GetHistory() error {

	var purchaseOrderItemHistories []PurchaseOrderItemHistory

	err := PGDB.Where("purchase_order_item_id = ?", poi.ID).Find(&purchaseOrderItemHistories).Error
	if err != nil {
		return err
	}

	poi.History = purchaseOrderItemHistories

	return nil

}

func (poi *PurchaseOrderItem) GetRejections() error {

	var purchaseOrderItemRejections []PurchaseOrderItemRejection

	err := PGDB.Where("purchase_order_item_id = ?", poi.ID).Find(&purchaseOrderItemRejections).Error
	if err != nil {
		return err
	}

	if len(purchaseOrderItemRejections) > 0 {
		//add images to rejections
		for i, _ := range purchaseOrderItemRejections {
			err = purchaseOrderItemRejections[i].GetImages()
			if err != nil {
				return err
			}
		}
	}

	poi.Rejections = purchaseOrderItemRejections

	return nil

}

func (poi *PurchaseOrderItem) GetProductLots() error {

	var productLots []ProductLot

	err := PGDB.Where("product_id = ?", poi.ProductID).Find(&productLots).Error
	if err != nil {
		return err
	}

	poi.ProductLots = productLots

	return nil

}

type PurchaseOrderItemUpdateIPARequest struct {
	Weight float64 `json:"weight"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Length float64 `json:"length"`
}

func (poiuipar *PurchaseOrderItemUpdateIPARequest) ParseAndValidateRequest(r *http.Request) []string {

	var errs []string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := struct {
		Weight json.RawMessage `json:"weight"`
		Width  json.RawMessage `json:"width"`
		Height json.RawMessage `json:"height"`
		Length json.RawMessage `json:"length"`
	}{}

	if err := json.Unmarshal(body, &aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Weight != nil {
		if err := json.Unmarshal(aux.Weight, &poiuipar.Weight); err != nil {
			errs = append(errs, "weight must be of type float")
		} else if poiuipar.Weight < 0 {
			errs = append(errs, "weight must be a positive number")
		}
	}

	if aux.Width != nil {
		if err := json.Unmarshal(aux.Width, &poiuipar.Width); err != nil {
			errs = append(errs, "width must be of type float")
		} else if poiuipar.Width < 0 {
			errs = append(errs, "width must be a positive number")
		}
	}

	if aux.Height != nil {
		if err := json.Unmarshal(aux.Height, &poiuipar.Height); err != nil {
			errs = append(errs, "height must be of type float")
		} else if poiuipar.Height < 0 {
			errs = append(errs, "height must be a positive number")
		}
	}

	if aux.Length != nil {
		if err := json.Unmarshal(aux.Length, &poiuipar.Length); err != nil {
			errs = append(errs, "length must be of type float")
		} else if poiuipar.Length < 0 {
			errs = append(errs, "length must be a positive number")
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
