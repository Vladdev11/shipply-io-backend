package OrganizationHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ListShippingMethods(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to find organization", http.StatusBadRequest)
		return
	}

	request := models.ShippingMethodListRequest{}
	request.OrganizationID = user.Organization.ID
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	shippingMethods, err := user.Organization.GetShippingMethods(request)
	if err != nil {
		util.ErrorResponse(w, "failed to find shipping methods", http.StatusBadRequest)
		return
	}

	shippingMethodsJSON := []models.ShippingMethodReturnJSON{}

	for _, shippingMethod := range shippingMethods {
		shippingMethodsJSON = append(shippingMethodsJSON, shippingMethod.ConvertToReturnJSON())
	}

	util.JSONResponse(w, shippingMethodsJSON, http.StatusOK)

}

func GetShippingMethod(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to find organization", http.StatusBadRequest)
		return
	}

	shippingMethodID, err := util.GetIntFromPath(r, "shipping_method_id")
	if err != nil {
		util.ErrorResponse(w, "failed to get shipping method id", http.StatusBadRequest)
		return
	}

	shippingMethod, err := models.GetShippingMethodByID(shippingMethodID)
	if err != nil {
		util.ErrorResponse(w, "failed to find shipping method", http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(shippingMethod.StoreID)
	if err != nil {
		util.ErrorResponse(w, "failed to find store", http.StatusBadRequest)
		return
	}

	client, err := models.GetClientByID(store.ClientID)
	if err != nil {
		util.ErrorResponse(w, "failed to find client", http.StatusBadRequest)
		return
	}

	if client.OrganizationID != user.Organization.ID {
		util.ErrorResponse(w, "shipping method does not belong to this organization", http.StatusBadRequest)
		return
	}

	carrierConnections, err := models.GetCarrierConnectionsByClientID(client.ID, user.Organization.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to find carrier connections", http.StatusBadRequest)
		return
	}

	cheapestCarrier, err := models.GetCheapestCarrierConnection()
	if err != nil {
		util.ErrorResponse(w, "failed to find carrier connections", http.StatusBadRequest)
		return
	}

	carrierConnections = append(carrierConnections, *cheapestCarrier)

	//make map of carrierConnectionID to carrierConnection
	carrierConnectionMap := make(map[int]models.CarrierConnection)
	for _, carrierConnection := range carrierConnections {
		carrierConnectionMap[carrierConnection.ID] = carrierConnection
	}

	baseCarriers, err := models.GetCarriers()
	if err != nil {
		util.ErrorResponse(w, "failed to find carriers", http.StatusBadRequest)
		return
	}

	//make of map of carrier id to carrier
	carrierMap := make(map[int]models.Carrier)
	for _, carrier := range baseCarriers {
		carrierMap[carrier.ID] = carrier
	}

	var shippingMethodCarriers []models.ShippingMethodCarrier

	//decide if we need to init a new blank shipping method carrier settings
	if shippingMethod.Carriers == nil {

		for _, carrierConnection := range carrierConnections {

			//init the new carrier
			shippingMethodCarrier := models.ShippingMethodCarrier{
				CarrierConnectionID: carrierConnection.ID,
			}

			carrierConnectionCarrierServices, _ := carrierConnection.GetCarrierServices()

			for _, carrierConnectionService := range carrierConnectionCarrierServices {
				//init a new carrier service
				shippingMethodCarrierService := models.ShippingMethodCarrierService{
					ServiceName: carrierConnectionService.Name,
					ServiceCode: carrierConnectionService.ServiceCode,
					Enabled:     false,
				}

				//append to shippingMethodCarrier
				shippingMethodCarrier.Services = append(shippingMethodCarrier.Services, shippingMethodCarrierService)
			}

			//append shippingMethodCarrier to shippingMethod
			shippingMethodCarriers = append(shippingMethodCarriers, shippingMethodCarrier)
		}

	} else {
		//means there is already settings in here and we need to confirm all the carrier connections are in here

		//get a var that holds the shipping method carriers
		json.Unmarshal(shippingMethod.Carriers, &shippingMethodCarriers)
		if shippingMethodCarriers == nil {
			util.ErrorResponse(w, "failed to unmarshal shipping method carriers", http.StatusBadRequest)
			return
		}

		//create a map of carrier id to true/false
		shippingMethodCarrierMap := make(map[int]map[string]bool)
		for _, shippingMethodCarrier := range shippingMethodCarriers {
			for _, shippingMethodCarrierService := range shippingMethodCarrier.Services {
				if shippingMethodCarrierMap[shippingMethodCarrier.CarrierConnectionID] == nil {
					shippingMethodCarrierMap[shippingMethodCarrier.CarrierConnectionID] = make(map[string]bool)
				}
				shippingMethodCarrierMap[shippingMethodCarrier.CarrierConnectionID][shippingMethodCarrierService.ServiceCode] = true
			}
		}

		//loop through all the carrier connections
		for i, carrier := range carrierConnections {
			if len(shippingMethodCarrierMap[carrier.ID]) > 0 {

				//get the carrier services
				carrierServices, err := carrier.GetCarrierServices()
				if err != nil {
					util.ErrorResponse(w, "failed to find carrier services", http.StatusBadRequest)
					return
				}

				//loop through all the carrier services
				for _, carrierService := range carrierServices {

					//if the service code is not in the map, then we need to init a new carrier service
					if !shippingMethodCarrierMap[carrier.ID][carrierService.ServiceCode] {

						//init a new carrier service
						shippingMethodCarrierService := models.ShippingMethodCarrierService{
							ServiceName: carrierService.Name,
							ServiceCode: carrierService.ServiceCode,
							Enabled:     false,
						}

						//append to existing carrier
						shippingMethodCarriers[i].Services = append(shippingMethodCarriers[i].Services, shippingMethodCarrierService)
					}
				}

			} else {

				//init the new carrier
				shippingMethodCarrier := models.ShippingMethodCarrier{
					CarrierConnectionID: carrier.ID,
				}

				carrierConnectionCarrierServices, err := carrier.GetCarrierServices()
				if err != nil {
					util.ErrorResponse(w, "failed to find carrier services", http.StatusBadRequest)
					return
				}

				for _, carrierConnectionService := range carrierConnectionCarrierServices {
					//init a new carrier service
					shippingMethodCarrierService := models.ShippingMethodCarrierService{
						ServiceName: carrierConnectionService.Name,
						ServiceCode: carrierConnectionService.ServiceCode,
						Enabled:     false,
					}

					//append to shippingMethodCarrier
					shippingMethodCarrier.Services = append(shippingMethodCarrier.Services, shippingMethodCarrierService)
				}

				//append shippingMethodCarrier to shippingMethod
				shippingMethodCarriers = append(shippingMethodCarriers, shippingMethodCarrier)

			}
		}

	}

	shippingMethod.Carriers, _ = json.Marshal(shippingMethodCarriers)
	shippingMethod.Update()

	//update the shipping method carriers with most updated base carrier info
	for i, shippingMethodCarrier := range shippingMethodCarriers {
		shippingMethodCarriers[i].CarrierConnectionName = carrierConnectionMap[shippingMethodCarrier.CarrierConnectionID].ShipengineNickname
		carrierID := carrierConnectionMap[shippingMethodCarrier.CarrierConnectionID].CarrierID
		shippingMethodCarriers[i].CarrierThumbnailURL = carrierMap[carrierID].ThumbnailURL
	}

	shippingMethod.Carriers, _ = json.Marshal(shippingMethodCarriers)

	shippingMethodJSON := shippingMethod.ConvertToReturnJSON()
	util.JSONResponse(w, shippingMethodJSON, http.StatusOK)

}

func UpdateShippingMethod(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to find organization", http.StatusBadRequest)
		return
	}

	shippingMethodID, err := util.GetIntFromPath(r, "shipping_method_id")
	if err != nil {
		util.ErrorResponse(w, "failed to get shipping method id", http.StatusBadRequest)
		return
	}

	shippingMethod, err := models.GetShippingMethodByID(shippingMethodID)
	if err != nil {
		util.ErrorResponse(w, "failed to find shipping method", http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(shippingMethod.StoreID)
	if err != nil {
		util.ErrorResponse(w, "failed to find store", http.StatusBadRequest)
		return
	}

	client, err := models.GetClientByID(store.ClientID)
	if err != nil {
		util.ErrorResponse(w, "failed to find client", http.StatusBadRequest)
		return
	}

	if client.OrganizationID != user.Organization.ID {
		util.ErrorResponse(w, "shipping method does not belong to this organization", http.StatusBadRequest)
		return
	}

	request := models.ShippingMethodUpdateRequest{}
	errors := request.ParseAndValidateUpdateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	shippingMethod, err = request.ConvertUpdateRequestToShippingMethod(shippingMethod)
	if err != nil {
		util.ErrorResponse(w, "failed to convert update request to shipping method", http.StatusBadRequest)
		return
	}

	err = shippingMethod.Update()
	if err != nil {
		util.ErrorResponse(w, "failed to update shipping method", http.StatusBadRequest)
		return
	}

	//unmarshal the carriers
	var shippingMethodCarriers []models.ShippingMethodCarrier
	err = json.Unmarshal(shippingMethod.Carriers, &shippingMethodCarriers)

	carrierConnections, err := models.GetCarrierConnectionsByClientID(client.ID, user.Organization.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to find carrier connections", http.StatusBadRequest)
		return
	}

	cheapestCarrier, err := models.GetCheapestCarrierConnection()
	if err != nil {
		util.ErrorResponse(w, "failed to find carrier connections", http.StatusBadRequest)
		return
	}

	carrierConnections = append(carrierConnections, *cheapestCarrier)

	//make map of carrierConnectionID to carrierConnection
	carrierConnectionMap := make(map[int]models.CarrierConnection)
	for _, carrierConnection := range carrierConnections {
		carrierConnectionMap[carrierConnection.ID] = carrierConnection
	}

	baseCarriers, err := models.GetCarriers()
	if err != nil {
		util.ErrorResponse(w, "failed to find carriers", http.StatusBadRequest)
		return
	}

	//make of map of carrier id to carrier
	carrierMap := make(map[int]models.Carrier)
	for _, carrier := range baseCarriers {
		carrierMap[carrier.ID] = carrier
	}

	//update the shipping method carriers with most updated base carrier info
	for i, shippingMethodCarrier := range shippingMethodCarriers {
		shippingMethodCarriers[i].CarrierConnectionName = carrierConnectionMap[shippingMethodCarrier.CarrierConnectionID].ShipengineNickname
		carrierID := carrierConnectionMap[shippingMethodCarrier.CarrierConnectionID].CarrierID
		shippingMethodCarriers[i].CarrierThumbnailURL = carrierMap[carrierID].ThumbnailURL
	}

	shippingMethod.Carriers, _ = json.Marshal(shippingMethodCarriers)

	shippingMethodJSON := shippingMethod.ConvertToReturnJSON()

	util.JSONResponse(w, shippingMethodJSON, http.StatusOK)
}
