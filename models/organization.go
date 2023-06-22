package models

import (
	"errors"
	"fmt"

	"github.com/shipply-io/shipply-io-backend/util"
)

type Organization struct {
	ID   int
	Type int
	Name string

	Clients []Client `gorm:"-"`
}

type OrganizationReturnJSON struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`

	Clients []ClientReturnJSON `json:"clients"`
}

func (o *Organization) GetClients() error {
	clients, err := GetClientsByOrganizationID(o.ID)
	if err != nil {
		return err
	}
	o.Clients = clients
	return nil
}

func (o *Organization) ConvertToReturnJSON() *OrganizationReturnJSON {

	if o == nil {
		return nil
	}

	clients := make([]ClientReturnJSON, len(o.Clients))
	for i := range o.Clients {
		clients[i] = *o.Clients[i].ConvertToReturnJSON()
	}

	return &OrganizationReturnJSON{
		ID:      o.ID,
		Type:    util.OrganizationType[o.Type],
		Name:    o.Name,
		Clients: clients,
	}
}

func (o *Organization) IsClientOwner(clientID int) bool {

	if o.Clients == nil {
		o.GetClients()
	}

	for _, client := range o.Clients {
		if client.ID == clientID {
			return true
		}
	}
	return false
}

func (o *Organization) IsWarehouseOwner(warehouseID int) bool {

	if o.Clients == nil {
		o.GetClients()
	}

	var warehouse Warehouse
	err := PGDB.Where("id = ?", warehouseID).Where("organization_id = ?", o.ID).First(&warehouse).Error
	if err != nil {
		return false
	}

	return true
}

func GetOrganizationByID(id int) (Organization, error) {
	var organization Organization
	err := PGDB.Where("id = ?", id).First(&organization).Error
	return organization, err
}

func (organization *Organization) GetPurchaseOrders(polr PurchaseOrderListRequest) ([]PurchaseOrder, int, int, error) {

	query := polr.ConvertToOrganizationQuery()
	countQuery := polr.ConvertToOrganizationQuery()

	totalQuery := PGDB.Model(&PurchaseOrder{}).
		Joins("LEFT JOIN clients ON clients.id = purchase_orders.client_id").
		Joins("LEFT JOIN organizations ON organizations.id = clients.organization_id")
	query = query.Where("organizations.id = ?", organization.ID)

	var purchaseOrders []PurchaseOrder
	if err := query.Offset(polr.Offset).Limit(polr.Limit).Find(&purchaseOrders).Error; err != nil {
		return nil, 0, 0, err
	}

	var count int64
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	var total int64
	if err := totalQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	return purchaseOrders, int(count), int(total), nil
}

func (organization *Organization) GetOrders(olr OrdersListRequest) ([]Order, int, int, error) {

	query := olr.ConvertToOrganizationQuery()
	countQuery := olr.ConvertToOrganizationQuery()

	totalQuery := PGDB.Model(&Order{}).
		Joins("LEFT JOIN stores ON stores.id = orders.store_id").
		Joins("LEFT JOIN clients ON clients.id = stores.client_id").
		Joins("LEFT JOIN organizations ON organizations.id = clients.organization_id")
	query = query.Where("organizations.id = ?", organization.ID)

	var orders []Order
	if err := query.Offset(olr.Offset).Limit(olr.Limit).Find(&orders).Error; err != nil {
		return nil, 0, 0, err
	}

	var count int64
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	var total int64
	if err := totalQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	return orders, int(count), int(total), nil
}

func (organization *Organization) SearchProducts(spr ProductSearchRequest) ([]Product, error) {

	var products []Product

	query := PGDB.Model(&Product{})

	query = query.Joins("LEFT JOIN clients ON clients.id = products.client_id").
		Joins("LEFT JOIN organizations ON organizations.id = clients.organization_id")

	if spr.ClientID != 0 {
		query = query.Where("client_id = ?", spr.ClientID)
	}

	query = query.Where("organizations.id = ?", organization.ID)

	if spr.SearchValue != "" {
		query = query.Where(PGDB.Where("to_tsvector('english', products.name) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", spr.SearchValue)).
			Or("to_tsvector('english', products.sku) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", spr.SearchValue)).
			Or("to_tsvector('english', products.barcode) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", spr.SearchValue)))
	}

	query = query.Limit(10)
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil

}

func (organization *Organization) GetVendors(vlr VendorListRequest) ([]Vendor, int, int, error) {

	query := vlr.ConvertToOrganizationQuery()
	countQuery := vlr.ConvertToOrganizationQuery()

	totalQuery := PGDB.Model(&Vendor{}).
		Joins("LEFT JOIN clients ON clients.id = vendors.client_id").
		Joins("LEFT JOIN organizations ON organizations.id = clients.organization_id").
		Where("organizations.id = ?", organization.ID)

	var vendors []Vendor
	if err := query.Offset(vlr.Offset).Limit(vlr.Limit).Find(&vendors).Error; err != nil {
		return nil, 0, 0, err
	}

	var count int64
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	var total int64
	if err := totalQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	return vendors, int(total), int(count), nil

}

func (organization *Organization) GetPurchaseOrderStatuses(poslr PurchaseOrderStatusListRequest) ([]PurchaseOrderStatus, int, int, error) {

	query := poslr.ConvertToOrganizationQuery()
	countQuery := poslr.ConvertToOrganizationQuery()

	totalQuery := PGDB.Model(&PurchaseOrderStatus{}).
		Joins("LEFT JOIN clients ON clients.id = purchase_order_statuses.client_id").
		Joins("LEFT JOIN organizations ON organizations.id = clients.organization_id").
		Where("organizations.id = ?", organization.ID)

	var purchaseOrderStatuses []PurchaseOrderStatus
	if err := query.Offset(poslr.Offset).Limit(poslr.Limit).Find(&purchaseOrderStatuses).Error; err != nil {
		return nil, 0, 0, err
	}

	var count int64
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	var total int64
	if err := totalQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	return purchaseOrderStatuses, int(total), int(count), nil

}

func (organization *Organization) GetWarehouses(wlr WarehouseListRequest) ([]Warehouse, int, int, error) {

	query := wlr.ConvertToOrganizationQuery()
	countQuery := wlr.ConvertToOrganizationQuery()

	totalQuery := PGDB.Model(&Warehouse{}).Where("organization_id = ?", organization.ID)

	var warehouses []Warehouse
	if err := query.Offset(wlr.Offset).Limit(wlr.Limit).Find(&warehouses).Error; err != nil {
		return nil, 0, 0, err
	}

	var count int64
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	var total int64
	if err := totalQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	return warehouses, int(total), int(count), nil

}

func (organization *Organization) GetLocationTypes() ([]LocationType, error) {

	var locationTypes []LocationType
	if err := PGDB.Model(&LocationType{}).Where("organization_id = ?", organization.ID).Find(&locationTypes).Error; err != nil {
		return nil, err
	}

	return locationTypes, nil

}

func (organization *Organization) GetCarrierConnections() ([]CarrierConnection, error) {
	err := organization.GetClients()
	if err != nil {
		return nil, errors.New("failed to get clients for organization")
	}

	var clientIDs []int
	for _, client := range organization.Clients {
		clientIDs = append(clientIDs, client.ID)
	}

	var carrierConnections []CarrierConnection
	if err := PGDB.Where("(owner_id = ? AND owner_type = 1) OR (owner_id IN (?) AND owner_type = 2)", organization.ID, clientIDs).Find(&carrierConnections).Error; err != nil {
		return nil, err
	}

	return carrierConnections, nil
}

func (organization *Organization) GetBoxes() ([]Box, error) {
	var boxes []Box
	if err := PGDB.Model(&Box{}).Where("organization_id = ?", organization.ID).Find(&boxes).Error; err != nil {
		return nil, err
	}

	return boxes, nil
}

func (o *Organization) GetLocations(llr LocationListRequest) ([]Location, int, int, error) {

	llr.OrganizationID = o.ID
	query := llr.ConvertToOrganizationQuery()
	countQuery := llr.ConvertToOrganizationQuery()

	totalQuery := PGDB.Model(&Location{}).
		Joins("LEFT JOIN warehouses ON warehouses.id = locations.warehouse_id").
		Where("warehouses.organization_id = ?", o.ID)

	var locations []Location
	if err := query.Offset(llr.Offset).Limit(llr.Limit).Find(&locations).Error; err != nil {
		return nil, 0, 0, err
	}

	var count int64
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	var total int64
	if err := totalQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	return locations, int(count), int(total), nil

}

func (organization *Organization) GetShippingMethods(request ShippingMethodListRequest) ([]ShippingMethod, error) {

	var shippingMethods []ShippingMethod

	stores := []int{}

	//get organization clients
	err := organization.GetClients()
	if err != nil {
		return nil, err
	}

	for _, client := range organization.Clients {
		err := client.GetStores()
		if err != nil {
			return nil, err
		}

		//get stores
		for _, store := range client.Stores {
			if request.ClientID != nil {
				if store.ClientID == *request.ClientID {
					stores = append(stores, store.ID)
				}
			} else {
				stores = append(stores, store.ID)
			}
		}

	}

	query := PGDB.Model(&ShippingMethod{})
	query = query.Where("store_id IN (?)", stores)

	if request.Mapped != nil {
		if *request.Mapped {
			query = query.Where("shipping_method_id IS NOT FALSE")
		} else {
			query = query.Where("shipping_method_id IS NOT TRUE")
		}
	}

	err = query.Find(&shippingMethods).Error
	if err != nil {
		return nil, err
	}

	return shippingMethods, nil

}

func (organization *Organization) GetStores(request StoreListRequest) ([]Store, error) {

	var stores []Store

	err := PGDB.Joins("inner join clients on stores.client_id = clients.id").Where("clients.organization_id = ?", request.OrganizationID).Find(&stores).Error
	if err != nil {
		return nil, err
	}

	return stores, nil

}
