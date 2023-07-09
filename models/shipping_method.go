package models

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
)

type ShippingMethod struct {
	ID       int
	StoreID  int
	Name     string
	Mapped   bool
	MappedBy int
	Priority int
	Carriers json.RawMessage `gorm:"type:jsonb"` //This is []ShippingMethodCarrier

	CreatedAt time.Time
	UpdatedAt time.Time

	Store Store
}

type ShippingMethodUpdateRequest struct {
	Priority int                        `json:"priority"`
	Carriers map[string]map[string]bool `json:"carriers"`
}

func (smur *ShippingMethodUpdateRequest) ParseAndValidateUpdateRequest(r *http.Request) []string {
	errors := []string{}

	aux := &struct {
		Priority json.RawMessage            `json:"priority"`
		Carriers map[string]map[string]bool `json:"carriers"`
	}{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Priority == nil {
		errors = append(errors, "priority is required")
	} else {
		if err := json.Unmarshal(aux.Priority, &smur.Priority); err != nil {
			errors = append(errors, "priority must be an integer")
		}
	}

	if aux.Carriers == nil {
		errors = append(errors, "carriers is required")
	}

	//loop through keys in aux.Carriers and make sure they are valid
	for carrierIdKey := range aux.Carriers {
		//parse carrierIdKey to int
		_, err := strconv.Atoi(carrierIdKey)
		if err != nil {
			errors = append(errors, "carrierID must be an integer")
		}

	}

	smur.Carriers = aux.Carriers

	if len(errors) > 0 {
		return errors
	}

	return nil
}

type ShippingMethodCarrierService struct {
	ServiceName string `json:"service_name"`
	ServiceCode string `json:"service_code"`
	Enabled     bool   `json:"enabled"`
}

type ShippingMethodCarrier struct {
	CarrierConnectionID   int                            `json:"carrier_connection_id"`
	CarrierConnectionName string                         `json:"carrier_connection_name,omitempty"`
	CarrierName           string                         `json:"carrier_name,omitempty"`
	CarrierID             int                            `json:"carrier_id,omitempty"`
	CarrierThumbnailURL   string                         `json:"carrier_thumbnail_url,omitempty"`
	Services              []ShippingMethodCarrierService `json:"carrier_services"`
}

type ShippingMethodReturnJSON struct {
	ID       int                     `json:"id"`
	StoreID  int                     `json:"store_id"`
	Name     string                  `json:"name"`
	Mapped   bool                    `json:"mapped"`
	MappedBy int                     `json:"mapped_by"`
	Priority int                     `json:"priority"`
	Carriers []ShippingMethodCarrier `json:"carriers"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (sm *ShippingMethod) ConvertToReturnJSON() ShippingMethodReturnJSON {
	ShippingMethodReturnJSON := ShippingMethodReturnJSON{
		ID:       sm.ID,
		StoreID:  sm.StoreID,
		Name:     sm.Name,
		Mapped:   sm.Mapped,
		MappedBy: sm.MappedBy,
		Priority: sm.Priority,

		CreatedAt: sm.CreatedAt,
		UpdatedAt: sm.UpdatedAt,
	}

	if sm.Carriers != nil {
		var carriers []ShippingMethodCarrier
		if err := json.Unmarshal(sm.Carriers, &carriers); err != nil {
			return ShippingMethodReturnJSON
		}
		ShippingMethodReturnJSON.Carriers = carriers
	}

	return ShippingMethodReturnJSON
}

type ShippingMethodListRequest struct {
	ClientID       *int  `json:"client_id"`
	OrganizationID int   `json:"organization_id"`
	Mapped         *bool `json:"mapped"`
}

func (sm *ShippingMethodListRequest) ParseAndValidateRequest(r *http.Request) []string {
	errors := []string{}

	clientID, err := util.GetIntQueryParam(r, "client_id")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "client_id must be an integer")
		} else {
			sm.ClientID = &clientID
		}
	}

	mapped, err := util.GetBoolQueryParam(r, "mapped")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "mapped must be a boolean")
		} else {
			sm.Mapped = &mapped
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (sm *ShippingMethod) Create(ctx context.Context) error {
	return util.DBFromContext(ctx).Create(sm).Error
}

func (sm *ShippingMethod) Update(ctx context.Context) error {
	return util.DBFromContext(ctx).Save(sm).Error
}

func (sm *ShippingMethod) Delete(ctx context.Context) error {
	return util.DBFromContext(ctx).Delete(sm).Error
}

func GetShippingMethodByStoreAndName(ctx context.Context, storeID int, name string) (*ShippingMethod, error) {
	var sm ShippingMethod
	err := util.DBFromContext(ctx).Where("store_id = ? AND name = ?", storeID, name).First(&sm).Error
	if err != nil {
		return nil, err
	}
	return &sm, nil
}
func GetShippingMethodByID(ctx context.Context, id int) (*ShippingMethod, error) {
	sm := &ShippingMethod{}
	err := util.DBFromContext(ctx).Where("id = ?", id).First(sm).Error
	return sm, err
}

func (smur *ShippingMethodUpdateRequest) ConvertUpdateRequestToShippingMethod(sm *ShippingMethod) (*ShippingMethod, error) {

	mapped := false

	sm.Priority = smur.Priority

	var shippingMethodCarriers []ShippingMethodCarrier
	err := json.Unmarshal(sm.Carriers, &shippingMethodCarriers)
	if err != nil {
		return nil, err
	}

	// Create a map of maps for the services in the request

	for key := range smur.Carriers {

		// Parse key to int
		carrierConnectionID, err := strconv.Atoi(key)
		if err != nil {
			continue
		}

		// Create a map of the services in the request so we only set values for services that are in the request
		serviceMap := make(map[string]bool)
		for service := range smur.Carriers[key] {
			serviceMap[service] = true
		}

		//get the carrier connection from the shipping method carriers
		for i, carrier := range shippingMethodCarriers {
			if carrier.CarrierConnectionID == carrierConnectionID {
				for j, service := range carrier.Services {
					// If the service is in the request, set the value to the request value
					if _, ok := serviceMap[service.ServiceCode]; ok {
						shippingMethodCarriers[i].Services[j].Enabled = smur.Carriers[key][service.ServiceCode]
						if smur.Carriers[key][service.ServiceCode] {
							mapped = true
						}
					}

				}
			}
		}

	}

	sm.Mapped = mapped
	sm.Carriers, err = json.Marshal(shippingMethodCarriers)
	if err != nil {
		return nil, err
	}

	return sm, nil
}
