package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"gorm.io/gorm"
)

type PickSession struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	WarehouseID int       `json:"warehouse_id"`
	Completed   bool      `json:"completed"`
	CompletedAt time.Time `json:"completed_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	User              User
	Warehouse         Warehouse
	PickSessionOrders []PickSessionOrder
}

type PickSessionCreateRequest struct {
	WarehouseID int `json:"warehouse_id"`
	ToteCount   int `json:"tote_count"`
}

type PickSessionSelectItemRequest struct {
	Barcode string `json:"barcode"`
}

type PickSessionSelectItemResponse struct {
	ToteName                string `json:"tote_name"`
	PickSessionOrderID      int    `json:"pick_session_order_id"`
	QuantityRemainingToPick int    `json:"quantity_remaining_to_pick"`
}

type PickSessionAssignToteRequest struct {
	ToteBarcode        string `json:"tote_barcode"`
	PickSessionOrderID int    `json:"pick_session_order_id"`
}

type PickSessionConfirmToteRequest struct {
	PickSessionOrderID int    `json:"pick_session_order_id"`
	ToteBarcode        string `json:"tote_barcode"`
}

type PickSessionPickRequest struct {
	ProductBarcode string `json:"product_barcode"`
	ToteBarcode    string `json:"tote_barcode"`
}

type PickSessionPickResponse struct {
	IsProductPickingComplete bool   `json:"is_product_picking_complete"`
	IsNewTote                bool   `json:"is_new_tote"`
	ToteName                 string `json:"tote_name"`
	QuantityRemainingToPick  int    `json:"quantity_remaining_to_pick"`
	PickSessionOrderID       int    `json:"pick_session_order_id"`
}

// PickSessionResponse is the response for any scan or creation of a pick session.
type PickSessionResponse struct {
	ID        int                    `json:"id"`
	Progress  []PickSessionProgress  `json:"progress"`
	Remaining []PickSessionRemaining `json:"remaining"`
}

type PickSessionProgress struct {
	ToteName           string                    `json:"tote_name"`
	OrderID            int                       `json:"order_id"`
	PickSessionOrderID int                       `json:"pick_session_order_id"`
	Items              []PickSessionProgressItem `json:"items"`
}

type PickSessionProgressItem struct {
	ProductID      int    `json:"product_id"`
	ProductName    string `json:"product_name"`
	ImageURL       string `json:"image_url"`
	QuantityPicked int    `json:"quantity_picked"`
	QuantityToPick int    `json:"quantity_to_pick"`
}

type PickSessionRemaining struct {
	ImageURL     string `json:"image_url"`
	Quantity     int    `json:"quantity"`
	ProductName  string `json:"product_name"`
	LocationName string `json:"location_name"`
	ProductID    int    `json:"product_id"`
}

func (ps *PickSession) Create() error {
	return PGDB.Create(ps).Error
}

func (ps *PickSession) Complete() error {
	ps.Completed = true
	ps.CompletedAt = time.Now()
	return PGDB.Save(ps).Error
}

func (ps *PickSession) GetPickSessionOrders() error {
	return PGDB.Model(ps).Association("PickSessionOrders").Find(&ps.PickSessionOrders)
}

func (ps *PickSession) CreatePickSessionOrders(toteCount int) ([]PickSessionOrder, error) {

	// TODO pick session order rules (partial picks, etc)

	// TODO prefilter to avoid having two products with the same barcode if possible

	// TODO optimize query
	var orders []Order
	err := PGDB.Where("NOT EXISTS (?)", PGDB.Table("pick_session_orders").
		Where("pick_session_orders.order_id = orders.id")).
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Where("order_items.quantity = order_items.allocated").
		Limit(toteCount).
		Find(&orders).Error
	if err != nil {
		return nil, errors.New("failed to get orders for pick session")
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("no orders ready for picking")
	}

	// Attempt to get orders without overlapping product barcodes
	var pickSessionOrders []PickSessionOrder
	for _, order := range orders {

		pickSessionOrder := PickSessionOrder{
			PickSessionID: ps.ID,
			OrderID:       order.ID,
		}
		err := pickSessionOrder.Create()
		if err != nil {
			return nil, err
		}

		pickSessionOrders = append(pickSessionOrders, pickSessionOrder)
	}

	return pickSessionOrders, nil

}

func (ps *PickSession) ConvertToPickSessionResponse() (*PickSessionResponse, error) {

	pickSessionResponse := &PickSessionResponse{
		ID: ps.ID,
	}

	err := ps.GetPickSessionOrders()
	if err != nil {
		return nil, err
	}

	for _, pickSessionOrder := range ps.PickSessionOrders {

		err := pickSessionOrder.GetPickSessionOrderItems()
		if err != nil {
			return nil, err
		}

		progress := PickSessionProgress{
			OrderID:            pickSessionOrder.OrderID,
			PickSessionOrderID: pickSessionOrder.ID,
		}

		if pickSessionOrder.LocationID != nil {
			tote, err := GetLocationByID(*pickSessionOrder.LocationID)
			if err != nil {
				return nil, err
			}

			if !tote.IsTote {
				return nil, errors.New("location is not a tote")
			}

			progress.ToteName = tote.Name
		}

		for _, pickSessionOrderItem := range pickSessionOrder.PickSessionOrderItems {

			product, err := GetProductByID(pickSessionOrderItem.ProductID)
			if err != nil {
				return nil, err
			}

			progress.Items = append(progress.Items, PickSessionProgressItem{
				ProductID:      pickSessionOrderItem.ProductID,
				ProductName:    product.Name,
				ImageURL:       product.ImageURL,
				QuantityPicked: pickSessionOrderItem.QuantityPicked,
				QuantityToPick: pickSessionOrderItem.QuantityToPick,
			})

			if pickSessionOrderItem.QuantityPicked >= pickSessionOrderItem.QuantityToPick {
				continue
			}

			alreadyHasProduct := false
			for _, pickSessionRemainingItem := range pickSessionResponse.Remaining {

				if pickSessionRemainingItem.ProductID == pickSessionOrderItem.ProductID {
					pickSessionRemainingItem.Quantity = pickSessionRemainingItem.Quantity + pickSessionOrderItem.QuantityToPick - pickSessionOrderItem.QuantityPicked
					alreadyHasProduct = true
				}

			}

			if !alreadyHasProduct {

				location, err := GetLocationByID(pickSessionOrderItem.LocationID)
				if err != nil {
					return nil, err
				}

				pickSessionResponse.Remaining = append(pickSessionResponse.Remaining, PickSessionRemaining{
					Quantity:     pickSessionOrderItem.QuantityToPick - pickSessionOrderItem.QuantityPicked,
					ProductID:    pickSessionOrderItem.ProductID,
					ImageURL:     product.ImageURL,
					ProductName:  product.Name,
					LocationName: location.Name,
				})

			}

		}

		pickSessionResponse.Progress = append(pickSessionResponse.Progress, progress)

	}

	// sort by location name
	sort.Slice(pickSessionResponse.Remaining, func(i, j int) bool {
		return pickSessionResponse.Remaining[i].LocationName < pickSessionResponse.Remaining[j].LocationName
	})

	return pickSessionResponse, nil

}

func (pscr *PickSessionCreateRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errs []string

	user, err := GetRequestingUser(r)
	if err != nil {
		return []string{"invalid user"}
	}
	user.GetOrganization()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		WarehouseID json.RawMessage `json:"warehouse_id"`
		ToteCount   json.RawMessage `json:"tote_count"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.WarehouseID == nil {
		errs = append(errs, "warehouse_id is required")
	} else if err := json.Unmarshal(aux.WarehouseID, &pscr.WarehouseID); err != nil {
		errs = append(errs, "warehouse_id must be an integer")
	} else {
		if !user.Organization.IsWarehouseOwner(pscr.WarehouseID) {
			return []string{"you do not have permission to create a pick session at that warehouse"}
		}
	}

	if aux.ToteCount == nil {
		errs = append(errs, "tote_count is required")
	} else if err := json.Unmarshal(aux.ToteCount, &pscr.ToteCount); err != nil {
		errs = append(errs, "tote_count must be an integer")
	} else if pscr.ToteCount < 1 {
		errs = append(errs, "tote_count must be greater than 0")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func GetCurrentPickSessionByUserID(userID int) (*PickSession, error) {

	var pickSession PickSession
	err := PGDB.Where("user_id = ? AND completed IS NOT TRUE", userID).First(&pickSession).Error
	if err != nil {
		return nil, err
	}

	return &pickSession, nil

}

func (pssir *PickSessionSelectItemRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errs []string

	user, err := GetRequestingUser(r)
	if err != nil {
		return []string{"invalid user"}
	}
	user.GetOrganization()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Barcode json.RawMessage `json:"barcode"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Barcode == nil {
		errs = append(errs, "barcode is required")
	} else if err := json.Unmarshal(aux.Barcode, &pssir.Barcode); err != nil {
		errs = append(errs, "barcode must be a string")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (psatr *PickSessionAssignToteRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errs []string

	user, err := GetRequestingUser(r)
	if err != nil {
		return []string{"invalid user"}
	}
	user.GetOrganization()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		ToteBarcode        json.RawMessage `json:"tote_barcode"`
		PickSessionOrderID json.RawMessage `json:"pick_session_order_id"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.ToteBarcode == nil {
		errs = append(errs, "tote_barcode is required")
	} else if err := json.Unmarshal(aux.ToteBarcode, &psatr.ToteBarcode); err != nil {
		errs = append(errs, "tote_barcode must be an string")
	} else if len(psatr.ToteBarcode) < 1 {
		errs = append(errs, "tote_barcode must be at least 1 character")
	}

	if aux.PickSessionOrderID == nil {
		errs = append(errs, "pick_session_order_id is required")
	} else if err := json.Unmarshal(aux.PickSessionOrderID, &psatr.PickSessionOrderID); err != nil {
		errs = append(errs, "pick_session_order_id must be an integer")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (psctr *PickSessionConfirmToteRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errs []string

	user, err := GetRequestingUser(r)
	if err != nil {
		return []string{"invalid user"}
	}
	user.GetOrganization()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		PickSessionOrderID json.RawMessage `json:"pick_session_order_id"`
		ToteBarcode        json.RawMessage `json:"tote_barcode"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.PickSessionOrderID == nil {
		errs = append(errs, "pick_session_order_id is required")
	} else if err := json.Unmarshal(aux.PickSessionOrderID, &psctr.PickSessionOrderID); err != nil {
		errs = append(errs, "pick_session_order_id must be an integer")
	}

	if aux.ToteBarcode == nil {
		errs = append(errs, "tote_barcode is required")
	} else if err := json.Unmarshal(aux.ToteBarcode, &psctr.ToteBarcode); err != nil {
		errs = append(errs, "tote_barcode must be a string")
	} else if len(psctr.ToteBarcode) < 1 {
		errs = append(errs, "tote_barcode must be at least 1 character")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (pspr *PickSessionPickRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errs []string

	user, err := GetRequestingUser(r)
	if err != nil {
		return []string{"invalid user"}
	}
	user.GetOrganization()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		ProductBarcode json.RawMessage `json:"product_barcode"`
		ToteBarcode    json.RawMessage `json:"tote_barcode"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.ProductBarcode == nil {
		errs = append(errs, "product_barcode is required")
	} else if err := json.Unmarshal(aux.ProductBarcode, &pspr.ProductBarcode); err != nil {
		errs = append(errs, "product_barcode must be a string")
	}

	if aux.ToteBarcode == nil {
		errs = append(errs, "tote_barcode is required")
	} else if err := json.Unmarshal(aux.ToteBarcode, &pspr.ToteBarcode); err != nil {
		errs = append(errs, "tote_barcode must be a string")
	} else if len(pspr.ToteBarcode) < 1 {
		errs = append(errs, "tote_barcode must be at least 1 character")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
