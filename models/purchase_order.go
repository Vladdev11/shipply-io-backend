package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

// PurchaseOrder represents a purchase order
type PurchaseOrder struct {
	ID             int       `json:"id"`
	ClientID       int       `json:"client_id"`
	PONumber       string    `json:"po_number"`
	StatusID       int       `json:"status_id"`
	ExpectedDate   time.Time `json:"expected_date"`
	ShipDate       time.Time `json:"ship_date"`
	ClosedDate     time.Time `json:"closed_date"`
	Closed         bool      `json:"closed"`
	VendorID       int       `json:"vendor_id"`
	WarehouseID    int       `json:"warehouse_id"`
	TrackingNumber string    `json:"tracking_number"`
	TrackingURL    string    `json:"tracking_url"`
	WarehouseNotes string    `json:"warehouse_notes"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	Items       []PurchaseOrderItem       `json:"items" gorm:"-"`
	Tags        []PurchaseOrderTag        `json:"tags"`
	History     []PurchaseOrderHistory    `gorm:"-"`
	Attachments []PurchaseOrderAttachment `json:"attachments" gorm:"-"`
	Status      PurchaseOrderStatus       `json:"status" gorm:"-"`
	UniqueItems int                       `json:"unique_items" gorm:"-"`
	TotalItems  int                       `json:"total_items" gorm:"-"`

	Client    Client    `json:"client"`
	Warehouse Warehouse `json:"warehouse"`
}

func (po *PurchaseOrder) GetStatus(ctx context.Context) error {

	err := util.DBFromContext(ctx).Where("id = ?", po.StatusID).First(&po.Status).Error
	if err != nil {
		return err
	}

	return nil

}

// PurchaseOrderJSON represents a purchase order in json format
// There are some fields that are altered for json formatting
type PurchaseOrderReturnJSON struct {
	ID             int                              `json:"id"`
	Client         ClientReturnJSON                 `json:"client"`
	PONumber       string                           `json:"po_number"`
	Status         *PurchaseOrderStatusReturnJSON   `json:"status"`
	ExpectedDate   string                           `json:"expected_date"`
	ShipDate       string                           `json:"ship_date"`
	ClosedDate     string                           `json:"closed_date"`
	VendorID       int                              `json:"vendor_id"`
	WarehouseID    int                              `json:"warehouse_id"`
	TrackingNumber string                           `json:"tracking_number"`
	TrackingURL    string                           `json:"tracking_url"`
	WarehouseNotes string                           `json:"warehouse_notes"`
	CreatedAt      time.Time                        `json:"created_at"`
	UpdatedAt      time.Time                        `json:"updated_at"`
	Items          []PurchaseOrderItemReturnJSON    `json:"items" `
	Tags           []string                         `json:"tags" `
	UniqueItems    int                              `json:"unique_items" `
	TotalItems     int                              `json:"total_items" `
	History        []PurchaseOrderHistoryReturnJSON `json:"history" `
}

// PurchaseOrderCreateRequest represents a create request we receive from the client
type PurchaseOrderCreateRequest struct {
	ClientID       int        `json:"client_id"`
	VendorID       int        `json:"vendor_id"`
	WarehouseID    int        `json:"warehouse_id"`
	PONumber       string     `json:"po_number"`
	StatusID       int        `json:"status_id"`
	ExpectedDate   SingleDate `json:"expected_date"`
	ShipDate       SingleDate `json:"ship_date"`
	ClosedDate     SingleDate `json:"closed_date"`
	TrackingNumber string     `json:"tracking_number"`
	TrackingURL    string     `json:"tracking_url"`
	WarehouseNotes string     `json:"warehouse_notes"`
}

// PurchaseOrderUpdateRequest represents a update request we receive from the client
type PurchaseOrderUpdateRequest struct {
	PONumber       string     `json:"po_number"`
	StatusID       int        `json:"status_id"`
	ExpectedDate   SingleDate `json:"expected_date"`
	ShipDate       SingleDate `json:"ship_date"`
	ClosedDate     SingleDate `json:"closed_date"`
	VendorID       int        `json:"vendor_id"`
	WarehouseID    int        `json:"warehouse_id"`
	TrackingNumber string     `json:"tracking_number"`
	TrackingURL    string     `json:"tracking_url"`
	WarehouseNotes string     `json:"warehouse_notes"`
	Tags           []string   `json:"tags"`
}

type PurchaseOrderListRequest struct {
	ClientID       int    `json:"client_id"`
	OrganizationID int    `json:"organization_id"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	OrderBy        string `json:"order_by"`
	OrderByColumn  string `json:"order_by_column"`
	SearchValue    string `json:"search_value"`
}

// CreatePurchaseOrder creates a purchase order and returns the purchase order
func CreatePurchaseOrder(ctx context.Context, purchaseOrder *PurchaseOrder) (*PurchaseOrder, error) {
	//create purchase order
	err := util.DBFromContext(ctx).Create(purchaseOrder).Error
	if err != nil {
		return nil, err
	}

	//return purchase order
	return purchaseOrder, nil
}

// GetPurchaseOrderByID returns a purchase order by id
func GetPurchaseOrderByID(ctx context.Context, purchaseOrderID int) (*PurchaseOrder, error) {
	//get purchase order from database
	purchaseOrder := &PurchaseOrder{}
	err := util.DBFromContext(ctx).First(purchaseOrder, purchaseOrderID).Error
	if err != nil {
		return nil, err
	}

	//return purchase order
	return purchaseOrder, nil
}

// UpdatePurchaseOrder updates a purchase order and returns the purchase order
func UpdatePurchaseOrder(ctx context.Context, purchaseOrder *PurchaseOrder) (*PurchaseOrder, error) {
	//update purchase order
	err := util.DBFromContext(ctx).Save(purchaseOrder).Error
	if err != nil {
		return nil, err
	}

	//return purchase order
	return purchaseOrder, nil
}

func (pour *PurchaseOrderUpdateRequest) ParseAndValidateRequest(r *http.Request) []string {
	ctx := r.Context()

	var errs []string

	user, err := GetRequestingUser(r)
	if err != nil {
		errs = append(errs, "invalid user")
		return errs
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		PONumber       json.RawMessage `json:"po_number"`
		StatusID       json.RawMessage `json:"status_id"`
		ExpectedDate   json.RawMessage `json:"expected_date"`
		ShipDate       json.RawMessage `json:"ship_date"`
		ClosedDate     json.RawMessage `json:"closed_date"`
		VendorID       json.RawMessage `json:"vendor_id"`
		WarehouseID    json.RawMessage `json:"warehouse_id"`
		TrackingNumber json.RawMessage `json:"tracking_number"`
		TrackingURL    json.RawMessage `json:"tracking_url"`
		WarehouseNotes json.RawMessage `json:"warehouse_notes"`
		Tags           json.RawMessage `json:"tags"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		errs = append(errs, "invalid purchase order id")
	}

	purchaseOrder, err := GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		errs = append(errs, "invalid purchase order id")
		return []string{"invalid purchase order id"}
	}

	if purchaseOrder.ClientID == 0 {
		errs = append(errs, "purchase order does not belong to a client")
	}

	if aux.PONumber != nil {
		if err := json.Unmarshal(aux.PONumber, &pour.PONumber); err != nil {
			errs = append(errs, "po_number must be a string")
		} else if len(pour.PONumber) > 255 {
			errs = append(errs, "po_number must be less than 255 characters")
		}
	}

	if aux.StatusID != nil {
		if err := json.Unmarshal(aux.StatusID, &pour.StatusID); err != nil {
			errs = append(errs, "status_id must be an integer")
		} else {
			status, err := GetPurchaseOrderStatusByID(ctx, pour.StatusID)
			if err != nil {
				errs = append(errs, "invalid status_id")
			} else {
				if status.ClientID != purchaseOrder.ClientID {
					errs = append(errs, "status does not belong to the client")
				}
			}
		}
	}

	if aux.ExpectedDate != nil {
		if err := json.Unmarshal(aux.ExpectedDate, &pour.ExpectedDate); err != nil {
			errs = append(errs, "expected_date must be in the format YYYY-MM-DD")
		}
	}

	if aux.ShipDate != nil {
		if err := json.Unmarshal(aux.ShipDate, &pour.ShipDate); err != nil {
			errs = append(errs, "ship_date must be in the format YYYY-MM-DD")
		}
	}

	if aux.ClosedDate != nil {
		if err := json.Unmarshal(aux.ClosedDate, &pour.ClosedDate); err != nil {
			errs = append(errs, "closed_date must be in the format YYYY-MM-DD")
		}
	}

	if aux.VendorID != nil {
		if err := json.Unmarshal(aux.VendorID, &pour.VendorID); err != nil {
			errs = append(errs, "vendor_id must be an integer")
		} else {
			vendor, err := GetVendorByID(ctx, pour.VendorID)
			if err != nil {
				errs = append(errs, fmt.Sprintf("vendor_id %d does not exist", pour.VendorID))
			} else if vendor.ClientID != purchaseOrder.ClientID {
				errs = append(errs, fmt.Sprintf("vendor_id %d does not belong to client_id %d", pour.VendorID, purchaseOrder.ClientID))
			}
		}
	}

	if aux.WarehouseID != nil {
		if err := json.Unmarshal(aux.WarehouseID, &pour.WarehouseID); err != nil {
			errs = append(errs, "warehouse_id must be an integer")
		} else {
			warehouse, err := GetWarehouseByID(ctx, pour.WarehouseID)
			if err != nil {
				errs = append(errs, "warehouse_id must be a valid warehouse")
			} else if user.GetRole() == "organization_admin" || user.GetRole() == "organization_user" {
				if warehouse.OrganizationID != user.OwnerID {
					errs = append(errs, fmt.Sprintf("warehouse_id %d does not belong to your organization", pour.WarehouseID))
				}
			} else {
				user.GetClient(ctx)
				if warehouse.OrganizationID != user.Client.OrganizationID {
					errs = append(errs, "you do not have access to this warehouse")
				}
			}
		}
	}

	if aux.TrackingNumber != nil {
		if err := json.Unmarshal(aux.TrackingNumber, &pour.TrackingNumber); err != nil {
			errs = append(errs, "tracking_number must be a string")
		} else if len(pour.TrackingNumber) > 255 {
			errs = append(errs, "tracking_number must be less than 255 characters")
		}
	}

	if aux.TrackingURL != nil {
		if err := json.Unmarshal(aux.TrackingURL, &pour.TrackingURL); err != nil {
			errs = append(errs, "tracking_url must be a string")
		} else if len(pour.TrackingURL) > 255 {
			errs = append(errs, "tracking_url must be less than 255 characters")
		}
	}

	if aux.WarehouseNotes != nil {
		if err := json.Unmarshal(aux.WarehouseNotes, &pour.WarehouseNotes); err != nil {
			errs = append(errs, "warehouse_notes must be a string")
		} else if len(pour.WarehouseNotes) > 255 {
			errs = append(errs, "warehouse_notes must be less than 255 characters")
		}
	}

	if aux.Tags != nil {
		if err := json.Unmarshal(aux.Tags, &pour.Tags); err != nil {
			errs = append(errs, "tags must be an array of strings")
		}

		for _, tag := range pour.Tags {
			if len(tag) > 255 {
				errs = append(errs, "each tag must be less than 255 characters")
				break
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func (po *PurchaseOrder) Delete(ctx context.Context) error {
	db := util.DBFromContext(ctx)

	err := db.Delete(po).Error
	if err != nil {
		return err
	}

	err = db.Where("purchase_order_id = ?", po.ID).Delete(&PurchaseOrderItem{}).Error
	if err != nil {
		return err
	}

	return nil

}

func (po *PurchaseOrder) Update(ctx context.Context) error {
	err := util.DBFromContext(ctx).Save(po).Error
	if err != nil {
		return err
	}
	return nil
}

func (po *PurchaseOrder) GetHistory(ctx context.Context) error {

	histories, err := GetPurchaseOrderHistorysByID(ctx, po.ID)
	if err != nil {
		return err
	}

	for i := range histories {
		err = histories[i].GetCreatedByUser(ctx)
		if err != nil {
			return err
		}
	}

	po.History = histories

	return nil

}

func (po *PurchaseOrder) GetItems(ctx context.Context) error {

	err := util.DBFromContext(ctx).Where("purchase_order_id = ?", po.ID).Order("id ASC").Find(&po.Items).Error
	if err != nil {
		return err
	}

	totalItems := 0
	for i, item := range po.Items {
		//get product
		product, err := GetProductByID(ctx, item.ProductID)
		if err != nil {
			return err
		}
		po.Items[i].Product = product
		totalItems += item.Ordered
	}

	po.UniqueItems = len(po.Items)
	po.TotalItems = totalItems

	return nil
}

func (p *PurchaseOrder) GetAttachments(ctx context.Context) error {

	err := util.DBFromContext(ctx).Where("purchase_order_id = ?", p.ID).Order("id ASC").Find(&p.Attachments).Error
	if err != nil {
		return err
	}

	return nil

}

func (p *PurchaseOrder) GetClient(ctx context.Context) error {

	err := util.DBFromContext(ctx).Where("id = ?", p.ClientID).First(&p.Client).Error
	if err != nil {
		return err
	}

	return nil

}

func (p *PurchaseOrder) GetTags(ctx context.Context) error {

	var purchaseOrderTags []PurchaseOrderTag

	err := util.DBFromContext(ctx).Where("purchase_order_id = ?", p.ID).Find(&purchaseOrderTags).Error
	if err != nil {
		return err
	}

	p.Tags = purchaseOrderTags

	return nil
}

func (p *PurchaseOrder) HasTag(ctx context.Context, tagName string) bool {

	var purchaseOrderTag PurchaseOrderTag

	err := util.DBFromContext(ctx).Where("tag = ?", tagName).Where("purchase_order_id = ?", p.ID).First(&purchaseOrderTag).Error
	return err == nil
}

func (p *PurchaseOrder) AddTag(ctx context.Context, tagName string) error {

	purchaseOrderTag := &PurchaseOrderTag{
		PurchaseOrderID: p.ID,
		Tag:             tagName,
	}

	err := util.DBFromContext(ctx).Create(purchaseOrderTag).Error
	if err != nil {
		return err
	}

	p.Tags = append(p.Tags, *purchaseOrderTag)

	return nil
}

func (p *PurchaseOrder) RemoveTag(ctx context.Context, tagName string) error {

	err := util.DBFromContext(ctx).Where("tag = ?", tagName).Where("purchase_order_id = ?", p.ID).Delete(&PurchaseOrderTag{}).Error
	if err != nil {
		return err
	}

	for i, tag := range p.Tags {
		if tag.Tag == tagName {
			p.Tags = append(p.Tags[:i], p.Tags[i+1:]...)
			break
		}
	}

	return nil
}

// UnmarshalJSON unmarshals the purchase order create request and returns clean specific errors for fields
func (p *PurchaseOrderCreateRequest) ParseAndValidateRequest(r *http.Request) []string {
	ctx := r.Context()

	var errs []string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		ClientID       json.RawMessage `json:"client_id"`
		WarehouseID    json.RawMessage `json:"warehouse_id"`
		VendorID       json.RawMessage `json:"vendor_id"`
		PONumber       json.RawMessage `json:"po_number"`
		StatusID       json.RawMessage `json:"status_id"`
		ExpectedDate   json.RawMessage `json:"expected_date"`
		ShipDate       json.RawMessage `json:"ship_date"`
		ClosedDate     json.RawMessage `json:"closed_date"`
		TrackingNumber json.RawMessage `json:"tracking_number"`
		TrackingURL    json.RawMessage `json:"tracking_url"`
		WarehouseNotes json.RawMessage `json:"warehouse_notes"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	user, err := GetRequestingUser(r)
	if err != nil {
		return []string{"invalid user"}
	}

	if user.GetRole() == "organization_admin" || user.GetRole() == "organization_user" {
		if aux.ClientID == nil {
			errs = append(errs, "client_id is required")
		} else if err := json.Unmarshal(aux.ClientID, &p.ClientID); err != nil {
			errs = append(errs, "client_id must be an integer")
		} else if p.ClientID <= 0 {
			errs = append(errs, "client_id must be greater than 0")
		}
	} else {
		p.ClientID = user.OwnerID
	}

	if aux.WarehouseID == nil {
		errs = append(errs, "warehouse_id is required")
	} else if err := json.Unmarshal(aux.WarehouseID, &p.WarehouseID); err != nil {
		errs = append(errs, "warehouse_id must be an integer")
	} else if p.WarehouseID <= 0 {
		errs = append(errs, "warehouse_id must be greater than 0")
	} else {
		warehouse, err := GetWarehouseByID(ctx, p.WarehouseID)
		if err != nil {
			errs = append(errs, "warehouse_id must be a valid warehouse")
		} else if user.GetRole() == "organization_admin" || user.GetRole() == "organization_user" {
			if warehouse.OrganizationID != user.OwnerID {
				errs = append(errs, fmt.Sprintf("warehouse_id %d does not belong to your organization", p.WarehouseID))
			}
		} else {
			user.GetClient(ctx)
			if warehouse.OrganizationID != user.Client.OrganizationID {
				errs = append(errs, "you do not have access to this warehouse")
			}
		}
	}

	if aux.VendorID != nil {
		if err := json.Unmarshal(aux.VendorID, &p.VendorID); err != nil {
			errs = append(errs, "vendor_id must be an integer")
		} else if p.VendorID <= 0 {
			errs = append(errs, "vendor_id must be greater than 0")
		} else {
			vendor, err := GetVendorByID(ctx, p.VendorID)
			if err != nil {
				errs = append(errs, "vendor_id must be a valid vendor")
			} else if vendor.ClientID != p.ClientID {
				errs = append(errs, fmt.Sprintf("vendor_id %d does not belong to client id %d", p.VendorID, p.ClientID))
			}
		}
	}

	if aux.PONumber != nil {
		if err := json.Unmarshal(aux.PONumber, &p.PONumber); err != nil {
			errs = append(errs, "po_number must be a string")
		} else if len(p.PONumber) > 255 {
			errs = append(errs, "po_number must be less than 255 characters")
		}
	}

	if aux.StatusID != nil {
		if err := json.Unmarshal(aux.StatusID, &p.StatusID); err != nil {
			errs = append(errs, "status_id must be an int")
		} else {
			status, err := GetPurchaseOrderStatusByID(ctx, p.StatusID)
			if err != nil {
				errs = append(errs, "invalid status")
			} else {
				if status.ClientID != p.ClientID {
					errs = append(errs, "status does not belong to the client")
				}
			}
		}
	}

	if aux.ExpectedDate != nil {
		if err := json.Unmarshal(aux.ExpectedDate, &p.ExpectedDate); err != nil {
			errs = append(errs, "expected date must be a date of the format YYYY-MM-DD")
		}
	}

	if aux.ShipDate != nil {
		if err := json.Unmarshal(aux.ShipDate, &p.ShipDate); err != nil {
			errs = append(errs, "ship date must be a date of the format YYYY-MM-DD")
		}
	}

	if aux.ClosedDate != nil {
		if err := json.Unmarshal(aux.ClosedDate, &p.ClosedDate); err != nil {
			errs = append(errs, "closed date must be a date of the format YYYY-MM-DD")
		}
	}

	if aux.TrackingNumber != nil {
		if err := json.Unmarshal(aux.TrackingNumber, &p.TrackingNumber); err != nil {
			errs = append(errs, "tracking number must be a string")
		} else if len(p.TrackingNumber) > 255 {
			errs = append(errs, "tracking number must be less than 255 characters")
		}
	}

	if aux.TrackingURL != nil {
		if err := json.Unmarshal(aux.TrackingURL, &p.TrackingURL); err != nil {
			errs = append(errs, "tracking url must be a string")
		} else if len(p.TrackingURL) > 255 {
			errs = append(errs, "tracking url must be less than 255 characters")
		}
	}

	if aux.WarehouseNotes != nil {
		if err := json.Unmarshal(aux.WarehouseNotes, &p.WarehouseNotes); err != nil {
			errs = append(errs, "warehouse notes must be a string")
		} else if len(p.WarehouseNotes) > 255 {
			errs = append(errs, "warehouse notes must be less than 255 characters")
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// UpdateWithPurchaseOrderRequest updates the purchase order with the purchase order request and returns the updated purchase order
func (p *PurchaseOrder) UpdateWithPurchaseOrderUpdateRequest(poRequest *PurchaseOrderUpdateRequest) *PurchaseOrder {
	//update purchase order
	p.PONumber = poRequest.PONumber
	p.StatusID = poRequest.StatusID
	p.ExpectedDate = poRequest.ExpectedDate.Time
	p.ShipDate = poRequest.ShipDate.Time
	p.ClosedDate = poRequest.ClosedDate.Time
	p.VendorID = poRequest.VendorID
	p.WarehouseID = poRequest.WarehouseID
	p.TrackingNumber = poRequest.TrackingNumber
	p.TrackingURL = poRequest.TrackingURL
	p.WarehouseNotes = poRequest.WarehouseNotes

	return p
}

func (p *PurchaseOrder) ConvertToReturnJSON(ctx context.Context) *PurchaseOrderReturnJSON {
	//make sure items is an empty array if it is nil
	if p.Items == nil {
		p.Items = []PurchaseOrderItem{}
	}

	purchaseOrderItems := []PurchaseOrderItemReturnJSON{}
	for _, item := range p.Items {
		//set to 0 to omit from json response
		item.PurchaseOrderID = 0
		item.ProductID = 0
		purchaseOrderItems = append(purchaseOrderItems, item.ConvertToReturnJSON(ctx))
	}

	//make sure notes is an empty array if it is nil
	if p.History == nil {
		p.History = []PurchaseOrderHistory{}
	}

	PurchaseOrderHistorys := []PurchaseOrderHistoryReturnJSON{}
	for _, note := range p.History {
		note.PurchaseOrderID = 0
		PurchaseOrderHistorys = append(PurchaseOrderHistorys, note.ConvertToReturnJSON(ctx))
	}

	//convert purchase order tags to slice of strings
	tags := []string{}
	if p.Tags != nil {
		for _, tag := range p.Tags {
			tags = append(tags, tag.Tag)
		}
	}

	//get status
	status, _ := GetPurchaseOrderStatusByID(ctx, p.StatusID)

	//get client
	client, _ := GetClientByID(ctx, p.ClientID)

	return &PurchaseOrderReturnJSON{
		ID:             p.ID,
		Client:         *client.ConvertToReturnJSON(ctx),
		PONumber:       p.PONumber,
		Status:         status.ConvertToReturnJSON(),
		ExpectedDate:   p.ExpectedDate.Format("2006-01-02"),
		ShipDate:       p.ShipDate.Format("2006-01-02"),
		ClosedDate:     p.ClosedDate.Format("2006-01-02"),
		VendorID:       p.VendorID,
		WarehouseID:    p.WarehouseID,
		TrackingNumber: p.TrackingNumber,
		TrackingURL:    p.TrackingURL,
		WarehouseNotes: p.WarehouseNotes,
		Tags:           tags,
		UniqueItems:    p.UniqueItems,
		TotalItems:     p.TotalItems,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
		Items:          purchaseOrderItems,
		History:        PurchaseOrderHistorys,
	}
}

func DeletePurchaseOrder(ctx context.Context, purchaseOrder *PurchaseOrder) error {
	//delete purchase order
	if err := util.DBFromContext(ctx).Delete(purchaseOrder).Error; err != nil {
		return err
	}

	return nil
}

func (por *PurchaseOrderListRequest) ParseAndValidateRequest(r *http.Request) []string {

	errors := []string{}

	clientID, err := util.GetIntQueryParam(r, "client_id")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "client_id must be an integer")
		} else {
			por.ClientID = clientID
		}
	}

	por.Limit = 100
	limit, err := util.GetIntQueryParam(r, "limit")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "limit must be an integer")
		} else if limit < 0 {
			errors = append(errors, "limit must be greater than or equal to 0")
		} else {
			por.Limit = limit
		}
	}

	por.Offset = 0
	offset, err := util.GetIntQueryParam(r, "offset")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "offset must be an integer")
		} else if offset < 0 {
			errors = append(errors, "offset must be greater than or equal to 0")
		} else {
			por.Offset = offset
		}
	}

	por.OrderBy = "asc"
	orderBy, err := util.GetStringQueryParam(r, "order_by")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "order must be a string")
		} else if orderBy != "asc" && orderBy != "desc" {
			errors = append(errors, "order must be either asc or desc")
		} else {
			por.OrderBy = orderBy
		}
	}

	searchValue, err := util.GetStringQueryParam(r, "search_value")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "search_value must be a string")
		} else {
			por.SearchValue = searchValue
		}
	}

	orderByColumn, err := util.GetStringQueryParam(r, "order_by_column")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "order_by_column must be a string")
		} else {
			por.OrderByColumn = orderByColumn
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (por *PurchaseOrder) Create(ctx context.Context) error {
	if err := util.DBFromContext(ctx).Create(por).Error; err != nil {
		return err
	}

	return nil
}

func (por *PurchaseOrderListRequest) ConvertToClientQuery(ctx context.Context) *gorm.DB {

	query := util.DBFromContext(ctx).Model(&PurchaseOrder{}).
		Select("DISTINCT purchase_orders.*").
		Joins("LEFT JOIN purchase_order_items ON purchase_order_items.purchase_order_id = purchase_orders.id").
		Joins("LEFT JOIN products ON products.id = purchase_order_items.product_id").
		Joins("LEFT JOIN vendors ON vendors.id = purchase_orders.vendor_id").
		Joins("LEFT JOIN purchase_order_statuses ON purchase_order_statuses.id = purchase_orders.status_id")

	query = query.Where("purchase_orders.client_id = ?", por.ClientID)

	if por.SearchValue != "" {
		query = query.Where(util.DBFromContext(ctx).Where("to_tsvector('english', po_number) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue)).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', purchase_orders.tracking_number) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', purchase_orders.tracking_url) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', purchase_orders.warehouse_notes) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', products.name) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', products.sku) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', vendors.name) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', purchase_order_statuses.name) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))))
	}

	if por.OrderByColumn != "" {
		query = query.Order(por.OrderByColumn + " " + por.OrderBy)
	} else {
		if por.OrderBy == "asc" {
			query = query.Order("id asc")
		} else {
			query = query.Order("id desc")
		}
	}

	return query
}

func (por *PurchaseOrderListRequest) ConvertToOrganizationQuery(ctx context.Context) *gorm.DB {

	query := util.DBFromContext(ctx).Model(&PurchaseOrder{}).
		Select("DISTINCT purchase_orders.*").
		Joins("LEFT JOIN purchase_order_items ON purchase_order_items.purchase_order_id = purchase_orders.id").
		Joins("LEFT JOIN products ON products.id = purchase_order_items.product_id").
		Joins("LEFT JOIN vendors ON vendors.id = purchase_orders.vendor_id").
		Joins("LEFT JOIN purchase_order_statuses ON purchase_order_statuses.id = purchase_orders.status_id").
		Joins("LEFT JOIN clients ON clients.id = purchase_orders.client_id").
		Joins("LEFT JOIN organizations ON organizations.id = clients.organization_id")

	query = query.Where("organizations.id = ?", por.OrganizationID)

	if por.ClientID != 0 {
		query = query.Where("purchase_orders.client_id = ?", por.ClientID)
	}

	if por.SearchValue != "" {
		query = query.Where(util.DBFromContext(ctx).Where("to_tsvector('english', po_number) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue)).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', purchase_orders.tracking_number) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', purchase_orders.tracking_url) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', purchase_orders.warehouse_notes) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', products.name) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', products.sku) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', vendors.name) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', clients.name) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', purchase_order_statuses.name) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue))))
	}

	if por.OrderByColumn != "" {
		query = query.Order(por.OrderByColumn + " " + por.OrderBy)
	} else {
		if por.OrderBy == "asc" {
			query = query.Order("id asc")
		} else {
			query = query.Order("id desc")
		}
	}

	return query
}

func (por *PurchaseOrderListRequest) Search(ctx context.Context) ([]PurchaseOrder, int, int, error) {

	purchaseOrders := []PurchaseOrder{}
	var count int64
	var total int64

	query := por.ConvertToClientQuery(ctx)
	countQuery := por.ConvertToClientQuery(ctx)
	totalQuery := util.DBFromContext(ctx).Model(&PurchaseOrder{})

	if err := query.Limit(por.Limit).Offset(por.Offset).Find(&purchaseOrders).Error; err != nil {
		return nil, 0, 0, err
	}

	if err := countQuery.Distinct("purchase_orders.id").Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	if err := totalQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	return purchaseOrders, int(total), int(count), nil
}

func ConvertPurchaseOrdersToSearchResults(ctx context.Context, matchingPurchaseOrders []PurchaseOrder, total int, count int) (*SearchResults, error) {

	var purchaseOrders []*PurchaseOrderReturnJSON

	for _, purchaseOrder := range matchingPurchaseOrders {
		purchaseOrders = append(purchaseOrders, purchaseOrder.ConvertToReturnJSON(ctx))
	}

	results, err := util.ConvertStructsToInterfaces(purchaseOrders)
	if err != nil {
		return nil, err
	}

	searchResults := &SearchResults{
		TotalCount:    total,
		FilteredCount: count,
		Data:          results,
	}

	return searchResults, nil
}
func (por *PurchaseOrder) UpdateWithRequest(ctx context.Context, request PurchaseOrderUpdateRequest) error {

	if request.PONumber != "" {
		por.PONumber = request.PONumber
	}

	if request.StatusID != 0 {
		por.StatusID = request.StatusID
	}

	if !request.ExpectedDate.Time.IsZero() {
		por.ExpectedDate = request.ExpectedDate.Time
	}

	if !request.ShipDate.Time.IsZero() {
		por.ShipDate = request.ShipDate.Time
	}

	if !request.ClosedDate.Time.IsZero() {
		por.ClosedDate = request.ClosedDate.Time
	}

	if request.VendorID != 0 {
		por.VendorID = request.VendorID
	}

	if request.WarehouseID != 0 {
		por.WarehouseID = request.WarehouseID
	}

	if request.TrackingNumber != "" {
		por.TrackingNumber = request.TrackingNumber
	}

	if request.TrackingURL != "" {
		por.TrackingURL = request.TrackingURL
	}

	if request.WarehouseNotes != "" {
		por.WarehouseNotes = request.WarehouseNotes
	}

	if err := util.DBFromContext(ctx).Save(&por).Error; err != nil {
		return err
	}

	if request.Tags != nil {
		por.GetTags(ctx)

		var currentTags []string
		for _, tag := range por.Tags {
			currentTags = append(currentTags, tag.Tag)
		}

		// Delete tags that are no longer in the request
		for _, tag := range currentTags {
			if !util.SliceContains(request.Tags, tag) {
				por.RemoveTag(ctx, tag)
			}
		}

		// Add tags that are in the request but not in the purchase order
		for _, tag := range request.Tags {
			if !util.SliceContains(currentTags, tag) {
				por.AddTag(ctx, tag)
			}
		}
	}

	return nil

}

type PurchaseOrderItemScanInputRequest struct {
	Value string `json:"value"`
}

type PurchaseOrderItemScanInputResponse struct {
	PurchaseOrderItemID int `json:"purchase_order_item_id"`
	Quantity            int `json:"quantity"`
}

func (porisir *PurchaseOrderItemScanInputRequest) ParseAndValidateRequest(r *http.Request) []string {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Value json.RawMessage `json:"value"`
	}{}

	if err := json.Unmarshal(body, &aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Value != nil {
		if err := json.Unmarshal(aux.Value, &porisir.Value); err != nil {
			return []string{"value must be a string"}
		}
	} else {
		return []string{"value is required"}
	}

	return nil
}
