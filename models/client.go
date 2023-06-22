package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type Client struct {
	ID             int
	Name           string
	OrganizationID int
	Active         bool
	AvatarFileName string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Organization Organization
	Stores       []Store
}

func (c *Client) Create() error {
	err := PGDB.Create(c).Error
	return err
}

func (c *Client) Update() error {
	err := PGDB.Save(c).Error
	return err
}

type ClientReturnJSON struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
	AvatarURL string `json:"avatar_url"`
}

func (c *Client) GetOrganization() error {
	organization, err := GetOrganizationByID(c.OrganizationID)
	if err != nil {
		return err
	}
	c.Organization = organization
	return nil
}

func GetClientByID(id int) (Client, error) {
	var client Client
	err := PGDB.Where("id = ?", id).First(&client).Error
	return client, err
}

func GetClientByStoreID(storeID int) (Client, error) {
	var client Client
	err := PGDB.Where("id = (SELECT client_id FROM stores WHERE id = ?)", storeID).First(&client).Error
	return client, err
}

func DeleteClientByID(id int) error {
	err := PGDB.Delete(&Client{}, id).Error
	return err
}

func GetAllClients() ([]Client, error) {
	var clients []Client
	err := PGDB.Find(&clients).Error
	return clients, err
}

func GetClientsByOrganizationID(organizationID int) ([]Client, error) {
	var clients []Client
	err := PGDB.Where("organization_id = ?", organizationID).Find(&clients).Error
	return clients, err
}

func (c *Client) ConvertToReturnJSON() *ClientReturnJSON {

	if c.ID == 0 {
		return nil
	}

	avatarurl := ""
	if c.AvatarFileName != "" {
		avatarurl = fmt.Sprintf("%s/%s", util.ConfigCDNHost, c.AvatarFileName)
	}

	clientReturnJSON := ClientReturnJSON{
		ID:        c.ID,
		Name:      c.Name,
		Active:    c.Active,
		AvatarURL: avatarurl,
	}
	return &clientReturnJSON
}

func GetClientIDsByOrganizationID(organizationID int) ([]int, error) {
	clients, err := GetClientsByOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}
	clientIDs := make([]int, len(clients))
	for i, client := range clients {
		clientIDs[i] = client.ID
	}
	return clientIDs, nil
}

func (client *Client) GetPurchaseOrders(polr PurchaseOrderListRequest) ([]PurchaseOrder, int, int, error) {

	query := polr.ConvertToClientQuery()
	countQuery := polr.ConvertToClientQuery()
	totalQuery := PGDB.Model(&PurchaseOrder{}).Where("client_id = ?", client.ID)

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

func (client *Client) GetPurchaseOrderStatuses(poslr PurchaseOrderStatusListRequest) ([]PurchaseOrderStatus, int, int, error) {

	query := poslr.ConvertToClientQuery()
	countQuery := poslr.ConvertToClientQuery()
	totalQuery := PGDB.Model(&PurchaseOrderStatus{}).Where("client_id = ?", client.ID)

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

	return purchaseOrderStatuses, int(count), int(total), nil
}

func (client *Client) GetOrders(olr OrdersListRequest) ([]Order, int, int, error) {

	olr.ClientID = client.ID

	query := olr.ConvertToClientQuery()
	countQuery := olr.ConvertToClientQuery()
	totalQuery := PGDB.Model(&Order{}).
		Joins("JOIN stores ON stores.id = orders.store_id").
		Where("stores.client_id = ?", client.ID)

	var Orders []Order
	if err := query.Offset(olr.Offset).Limit(olr.Limit).Find(&Orders).Error; err != nil {
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

	return Orders, int(count), int(total), nil
}

func (client *Client) SearchProducts(spr ProductSearchRequest) ([]Product, error) {

	var products []Product

	query := PGDB.Model(&Product{})

	query = query.Where("client_id = ?", spr.ClientID)

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

func (client *Client) GetVendors(vlr VendorListRequest) ([]Vendor, int, int, error) {

	query := vlr.ConvertToClientQuery()
	countQuery := vlr.ConvertToClientQuery()

	totalQuery := PGDB.Model(&Vendor{}).Where("client_id = ?", client.ID)

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

	return vendors, int(count), int(total), nil

}

func (client *Client) GetCarrierConnections() ([]CarrierConnection, error) {

	var carrierConnections []CarrierConnection
	err := PGDB.Where("owner_id = ?", client.ID).
		Where("owner_type = ?", 2).
		Find(&carrierConnections).Error

	if err != nil {
		return nil, err
	}

	return carrierConnections, nil

}

func (client *Client) GetStores() error {

	var stores []Store
	err := PGDB.Where("client_id = ?", client.ID).Find(&stores).Error

	if err != nil {
		return err
	}

	client.Stores = stores

	return nil

}

type ClientCreateRequest struct {
	Name string `json:"name"`
}

func (c *ClientCreateRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errs []string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Name json.RawMessage `json:"name"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Name != nil {
		if err := json.Unmarshal(aux.Name, &c.Name); err != nil {
			errs = append(errs, "name must be a string")
		} else if len(c.Name) > 255 {
			errs = append(errs, "name must be less than 255 characters")
		} else if len(c.Name) == 0 {
			errs = append(errs, "name must be greater than 0 characters")
		}
	} else {
		errs = append(errs, "name is required")
	}

	//check if name is unique
	var count int64
	PGDB.Model(&Client{}).Where("name = ?", c.Name).Count(&count)
	if count > 0 {
		errs = append(errs, "client already exists with this name")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

type ClientUpdateAvatarRequest struct {
	File     multipart.File `json:"file"`
	FileType string         `json:"file_type"`
}

func (c *ClientUpdateAvatarRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errors []string

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errors = append(errors, err.Error())
		return errors
	}

	file, fileHeader, _ := r.FormFile("file")
	//check if there is a file
	if fileHeader != nil {
		//check if file is a valid image
		if fileHeader.Header.Get("Content-Type") != "image/jpeg" && fileHeader.Header.Get("Content-Type") != "image/png" {
			errors = append(errors, "file must be a valid image")
		}
		defer file.Close()
		c.File = file
		c.FileType = fileHeader.Header.Get("Content-Type")
	}
	if len(errors) > 0 {
		return errors
	}

	return nil
}

type ClientUpdateRequest struct {
	Name string `json:"name"`
}

func (c *ClientUpdateRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errors []string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Name json.RawMessage `json:"name"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Name != nil {
		if err := json.Unmarshal(aux.Name, &c.Name); err != nil {
			errors = append(errors, "name must be a string")
		} else if len(c.Name) > 255 {
			errors = append(errors, "name must be less than 255 characters")
		} else if len(c.Name) == 0 {
			errors = append(errors, "name must be greater than 0 characters")
		}
	} else {
		errors = append(errors, "name is required")
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
