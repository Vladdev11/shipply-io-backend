package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type PurchaseOrderStatus struct {
	ID          int
	ClientID    int
	Name        string
	StatusColor string
	TextColor   string

	Client *Client
}

type PurchaseOrderStatusCreateRequest struct {
	ClientID    int    `json:"client_id"`
	Name        string `json:"name"`
	StatusColor string `json:"status_color"`
	TextColor   string `json:"text_color"`
}

type PurchaseOrderStatusReturnJSON struct {
	ID          int    `json:"id"`
	ClientID    int    `json:"client_id,omitempty"`
	Name        string `json:"name"`
	StatusColor string `json:"status_color"`
	TextColor   string `json:"text_color"`
}

type PurchaseOrderStatusListRequest struct {
	ClientID       int    `json:"client_id"`
	OrganizationID int    `json:"organization_id"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	OrderBy        string `json:"order_by"`
	OrderByColumn  string `json:"order_by_column"`
	SearchValue    string `json:"search_value"`
}

type PurchaseOrderStatusUpdateRequest struct {
	Name        string `json:"name"`
	StatusColor string `json:"status_color"`
	TextColor   string `json:"text_color"`
}

func (pos *PurchaseOrderStatus) IsInUse(ctx context.Context) bool {
	var count int64
	err := util.DBFromContext(ctx).Model(&PurchaseOrder{}).Where("status = ?", pos.ID).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return false
	}
	return count > 0
}

func GetPurchaseOrderStatusByClientAndName(ctx context.Context, clientID int, name string) (*PurchaseOrderStatus, error) {
	var purchaseOrderStatus PurchaseOrderStatus
	err := util.DBFromContext(ctx).Where("client_id = ?", clientID).Where("name = ?", name).Find(&purchaseOrderStatus).Error
	if err != nil {
		return nil, err
	}
	return &purchaseOrderStatus, nil
}

func CreatePurchaseOrderStatus(ctx context.Context, pot *PurchaseOrderStatus) (*PurchaseOrderStatus, error) {
	err := util.DBFromContext(ctx).Create(pot).Error
	if err != nil {
		return nil, err
	}

	//return purchase order
	return pot, nil
}

func GetPurchaseOrderStatusByID(ctx context.Context, id int) (*PurchaseOrderStatus, error) {
	var purchaseOrderStatus PurchaseOrderStatus
	err := util.DBFromContext(ctx).Where("id = ?", id).First(&purchaseOrderStatus).Error
	if err != nil {
		return nil, err
	}
	return &purchaseOrderStatus, nil
}

func (poscr *PurchaseOrderStatus) Create(ctx context.Context) error {
	err := util.DBFromContext(ctx).Create(poscr).Error
	if err != nil {
		return err
	}
	return nil
}

func (potcr *PurchaseOrderStatusCreateRequest) ParseAndValidateRequest(r *http.Request) []string {

	errs := []string{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		ClientID    json.RawMessage `json:"client_id"`
		Name        json.RawMessage `json:"name"`
		StatusColor json.RawMessage `json:"status_color"`
		TextColor   json.RawMessage `json:"text_color"`
	}{}

	if err := json.Unmarshal(body, &aux); err != nil {
		errs = append(errs, ("invalid json"))
	}

	if aux.ClientID != nil {
		if err := json.Unmarshal(aux.ClientID, &potcr.ClientID); err != nil {
			errs = append(errs, ("client_id must be int"))
		} else {
			if _, err := GetClientByID(r.Context(), potcr.ClientID); err != nil {
				errs = append(errs, ("client_id does not exist"))
			}
		}
	} else {
		errs = append(errs, ("client_id is required"))
	}

	if aux.Name != nil {
		if err := json.Unmarshal(aux.Name, &potcr.Name); err != nil {
			errs = append(errs, ("name must be string"))
		}
	} else {
		errs = append(errs, ("name is required"))
	}

	if aux.StatusColor != nil {
		if err := json.Unmarshal(aux.StatusColor, &potcr.StatusColor); err != nil {
			errs = append(errs, ("status_color must be string"))
		}
	} else {
		errs = append(errs, ("status_color is required"))
	}

	if aux.TextColor != nil {
		if err := json.Unmarshal(aux.TextColor, &potcr.TextColor); err != nil {
			errs = append(errs, ("text_color must be string"))
		}
	} else {
		errs = append(errs, ("text_color is required"))
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func (potcr *PurchaseOrderStatusCreateRequest) ConvertToStatus() PurchaseOrderStatus {
	return PurchaseOrderStatus{
		ClientID:    potcr.ClientID,
		Name:        potcr.Name,
		StatusColor: potcr.StatusColor,
		TextColor:   potcr.TextColor,
	}
}

func (potcr *PurchaseOrderStatus) ConvertToReturnJSON() *PurchaseOrderStatusReturnJSON {

	if potcr == nil {
		return &PurchaseOrderStatusReturnJSON{}
	}

	return &PurchaseOrderStatusReturnJSON{
		ID:          potcr.ID,
		ClientID:    potcr.ClientID,
		Name:        potcr.Name,
		StatusColor: potcr.StatusColor,
		TextColor:   potcr.TextColor,
	}
}

func (r *PurchaseOrderStatusListRequest) ParseAndValidateRequest(req *http.Request) []string {
	errors := []string{}

	clientID, err := util.GetIntQueryParam(req, "client_id")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "client_id must be an integer")
		} else {
			r.ClientID = clientID
		}
	}

	r.Limit = 100
	limit, err := util.GetIntQueryParam(req, "limit")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "limit must be an integer")

		} else if limit < 0 || limit > 100 {
			errors = append(errors, "limit must be between 0 and 100")
		} else {
			r.Limit = limit
		}
	}

	r.Offset = 0
	offset, err := util.GetIntQueryParam(req, "offset")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "offset must be an integer")
		} else if offset < 0 {
			errors = append(errors, "offset must be greater than 0")
		} else {
			r.Offset = offset
		}
	}

	r.OrderBy = "asc"
	orderBy, err := util.GetStringQueryParam(req, "order_by")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "order_by must be a string")
		} else if orderBy != "asc" && orderBy != "desc" {
			errors = append(errors, "order_by must be either asc or desc")
		} else {
			r.OrderBy = orderBy
		}
	}

	searchValue, err := util.GetStringQueryParam(req, "search_value")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "search_value must be a string")
		} else {
			r.SearchValue = searchValue
		}
	}

	orderByColumn, err := util.GetStringQueryParam(req, "order_by_column")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "order_by_column must be a string")
		} else {
			r.OrderByColumn = orderByColumn
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (poslr *PurchaseOrderStatusListRequest) ConvertToOrganizationQuery(ctx context.Context) *gorm.DB {

	query := util.DBFromContext(ctx).Model(&PurchaseOrderStatus{}).
		Select("DISTINCT purchase_order_statuses.*").
		Joins("LEFT JOIN clients ON clients.id = purchase_order_statuses.client_id").
		Joins("LEFT JOIN organizations ON organizations.id = clients.organization_id")

	query = query.Where("organizations.id = ?", poslr.OrganizationID)

	if poslr.ClientID != 0 {
		query = query.Where("purchase_order_statuses.client_id = ?", poslr.ClientID)
	}

	if poslr.SearchValue != "" {
		query = query.Where("purchase_order_statuses.name ILIKE ?", "%"+poslr.SearchValue+"%")
	}

	if poslr.OrderByColumn != "" {
		query = query.Order(poslr.OrderByColumn + " " + poslr.OrderBy)
	} else {
		query = query.Order("purchase_order_statuses.id " + poslr.OrderBy)
	}

	return query
}

func (poslr *PurchaseOrderStatusListRequest) ConvertToClientQuery(ctx context.Context) *gorm.DB {

	query := util.DBFromContext(ctx).Model(&PurchaseOrderStatus{}).
		Select("DISTINCT purchase_order_statuses.*")

	query = query.Where("purchase_order_statuses.client_id = ?", poslr.ClientID)

	if poslr.SearchValue != "" {
		query = query.Where("purchase_order_statuses.name ILIKE ?", "%"+poslr.SearchValue+"%")
	}

	if poslr.OrderByColumn != "" {
		query = query.Order(poslr.OrderByColumn + " " + poslr.OrderBy)
	} else {
		query = query.Order("purchase_order_statuses.id " + poslr.OrderBy)
	}

	return query
}

func ConvertPurchaseOrderStatusesToSearchResults(matchingStatuses []PurchaseOrderStatus, total int, count int) (*SearchResults, error) {
	var purchaseOrderStatuses []*PurchaseOrderStatusReturnJSON

	for _, purchaseOrderStatus := range matchingStatuses {
		purchaseOrderStatuses = append(purchaseOrderStatuses, purchaseOrderStatus.ConvertToReturnJSON())
	}

	results, err := util.ConvertStructsToInterfaces(purchaseOrderStatuses)
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

func (posur *PurchaseOrderStatusUpdateRequest) ParseAndValidateRequest(r *http.Request) []string {
	ctx := r.Context()
	var errs []string

	poStatusID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		return []string{"id must be an integer"}
	}

	poStatus, err := GetPurchaseOrderStatusByID(ctx, poStatusID)
	if err != nil {
		return []string{"id must be a valid purchase order status id"}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Name        json.RawMessage `json:"name"`
		StatusColor json.RawMessage `json:"status_color"`
		TextColor   json.RawMessage `json:"text_color"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Name != nil {
		if err := json.Unmarshal(aux.Name, &posur.Name); err != nil {
			errs = append(errs, "name must be a string")
		} else if len(posur.Name) > 255 {
			errs = append(errs, "name must be less than 255 characters")
		} else {
			if poStatus.Name != posur.Name {
				poStatuses, err := GetPurchaseOrderStatusesByClientID(ctx, poStatus.ClientID)
				if err != nil {
					errs = append(errs, "error getting purchase order statuses")
				}

				for _, poStatus := range poStatuses {
					if poStatus.Name == posur.Name {
						errs = append(errs, "name already exists")
						break
					}
				}
			}
		}
	}

	if aux.StatusColor != nil {
		if err := json.Unmarshal(aux.StatusColor, &posur.StatusColor); err != nil {
			errs = append(errs, "status_color must be a string")
		} else if len(posur.StatusColor) != 6 {
			errs = append(errs, "status_color must be 6 characters")
		}
	}

	if aux.TextColor != nil {
		if err := json.Unmarshal(aux.TextColor, &posur.TextColor); err != nil {
			errs = append(errs, "text_color must be a string")
		} else if len(posur.TextColor) != 6 {
			errs = append(errs, "text_color must be 6 characters")
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func GetPurchaseOrderStatusesByClientID(ctx context.Context, clientID int) ([]PurchaseOrderStatus, error) {
	var purchaseOrderStatuses []PurchaseOrderStatus

	if err := util.DBFromContext(ctx).Where("client_id = ?", clientID).Find(&purchaseOrderStatuses).Error; err != nil {
		return nil, err
	}

	return purchaseOrderStatuses, nil
}

func (pos *PurchaseOrderStatus) UpdateWithRequest(ctx context.Context, request *PurchaseOrderStatusUpdateRequest) error {

	if request.Name != "" {
		pos.Name = request.Name
	}

	if request.StatusColor != "" {
		pos.StatusColor = request.StatusColor
	}

	if request.TextColor != "" {
		pos.TextColor = request.TextColor
	}

	if err := util.DBFromContext(ctx).Save(pos).Error; err != nil {
		return err
	}

	return nil

}

func (pos *PurchaseOrderStatus) Delete(ctx context.Context) error {
	if err := util.DBFromContext(ctx).Delete(pos).Error; err != nil {
		return err
	}

	return nil
}
