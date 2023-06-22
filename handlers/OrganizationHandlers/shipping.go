package OrganizationHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	ShipengineHandlers "github.com/shipply-io/shipply-io-backend/api/shipengine/handlers"
	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

// TODO add logging through whole file

func ShippingScanTote(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	request := models.ShippingScanToteRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	toteLocationID, err := strconv.Atoi(request.Barcode)
	if err != nil {
		util.ErrorResponse(w, "invalid tote barcode", http.StatusBadRequest)
		return
	}

	tote, err := models.GetLocationByID(toteLocationID)
	if err != nil {
		util.ErrorResponse(w, "failed to get tote", http.StatusBadRequest)
		return
	}

	if !tote.IsTote {
		util.ErrorResponse(w, "location is not a tote", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsWarehouseOwner(tote.WarehouseID) {
		util.ErrorResponse(w, "user does not have access to warehouse", http.StatusForbidden)
		return
	}

	pickSessionOrder, err := tote.GetActivePickSessionOrder()
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session order", http.StatusBadRequest)
		return
	}

	if !pickSessionOrder.Picked {
		util.ErrorResponse(w, "pick session order is not fully picked", http.StatusBadRequest)
		return
	}

	order, err := pickSessionOrder.GetOrder()
	if err != nil {
		util.ErrorResponse(w, "failed to get order", http.StatusBadRequest)
		return
	}

	err = order.GetShippingMethod()
	if err != nil {
		util.ErrorResponse(w, "failed to get shipping method", http.StatusBadRequest)
		return
	}

	if !order.ShippingMethod.Mapped {
		util.ErrorResponse(w, "shipping method is not mapped", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, models.ShippingScanToteResponse{
		PickSessionOrderID: pickSessionOrder.ID,
	}, http.StatusOK)

}

func ShippingGetPickSessionOrder(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	pickSessionOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid pick session order id", http.StatusBadRequest)
		return
	}

	pickSessionOrder, err := models.GetPickSessionOrderByID(pickSessionOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session order", http.StatusBadRequest)
		return
	}

	err = pickSessionOrder.GetPickSession()
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsWarehouseOwner(pickSessionOrder.PickSession.WarehouseID) {
		util.ErrorResponse(w, "user does not have access to warehouse", http.StatusForbidden)
		return
	}

	if !pickSessionOrder.Picked {
		util.ErrorResponse(w, "pick session order is not fully picked", http.StatusBadRequest)
		return
	}

	if pickSessionOrder.Shipped {
		util.ErrorResponse(w, "pick session order is already shipped", http.StatusBadRequest)
		return
	}

	if !pickSessionOrder.PickSession.Completed {
		util.ErrorResponse(w, "pick session is not completed", http.StatusBadRequest)
		return
	}

	order, err := pickSessionOrder.GetOrder()
	if err != nil {
		util.ErrorResponse(w, "failed to get order", http.StatusBadRequest)
		return
	}

	err = order.GetShippingMethod()
	if err != nil {
		util.ErrorResponse(w, "failed to get shipping method", http.StatusBadRequest)
		return
	}

	if !order.ShippingMethod.Mapped {
		util.ErrorResponse(w, "shipping method is not mapped", http.StatusBadRequest)
		return
	}

	err = pickSessionOrder.GetPickSessionOrderItems()
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session order items", http.StatusBadRequest)
		return
	}

	// TODO suggested boxes (Brennan)
	boxes, err := user.Organization.GetBoxes()
	if err != nil {
		util.ErrorResponse(w, "failed to get boxes", http.StatusBadRequest)
		return
	}

	response := models.ShippingGetPickSessionOrderResponse{
		PickSessionOrderID:      pickSessionOrder.ID,
		OrderID:                 order.ID,
		ToteID:                  *pickSessionOrder.LocationID,
		HasError:                pickSessionOrder.HasError,
		PickSessionOrderErrorID: pickSessionOrder.PickSessionOrderErrorID,
		Order:                   *order.ConvertToReturnJSON(),
	}

	if pickSessionOrder.ShippingRateID != 0 {
		err := pickSessionOrder.GetShippingRate()
		if err != nil {
			util.ErrorResponse(w, "failed to get rate", http.StatusBadRequest)
			return
		}

		response.Rate = pickSessionOrder.ShippingRate
	}

	for _, pickSessionOrderItem := range pickSessionOrder.PickSessionOrderItems {
		pickSessionOrderItemJSON := pickSessionOrderItem.ConvertToReturnJSON()
		response.PickSessionOrderItems = append(response.PickSessionOrderItems, pickSessionOrderItemJSON)
	}

	for _, box := range boxes {
		boxJSON := box.ConvertToReturnJSON()
		response.Boxes = append(response.Boxes, *boxJSON)
	}

	util.JSONResponse(w, response, http.StatusOK)

}

func ShippingShopRates(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	pickSessionOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid pick session order id", http.StatusBadRequest)
		return
	}

	pickSessionOrder, err := models.GetPickSessionOrderByID(pickSessionOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session order", http.StatusBadRequest)
		return
	}

	err = pickSessionOrder.GetPickSession()
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsWarehouseOwner(pickSessionOrder.PickSession.WarehouseID) {
		util.ErrorResponse(w, "user does not have access to warehouse", http.StatusForbidden)
		return
	}

	if !pickSessionOrder.Picked {
		util.ErrorResponse(w, "pick session order is not fully picked", http.StatusBadRequest)
		return
	}

	request := models.ShippingGetRatesRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	shipengineRateShopRequest, err := ShipengineModels.ConstructRateShopRequest(pickSessionOrder, request.BoxID, request.Weight, true)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	shipEngineRateShopResponse, err := ShipengineHandlers.ShopRates(*shipengineRateShopRequest)
	if err != nil {
		fmt.Println(err)
		util.ErrorResponse(w, "failed to shop rates", http.StatusBadRequest)
		return
	}

	response := models.ShippingGetRatesResponse{
		CarrierConnections: []models.CarrierConnectionRatesResponse{},
	}

	// Get a list of carrier connections that are active for response structure
	order, err := pickSessionOrder.GetOrder()
	if err != nil {
		util.ErrorResponse(w, "failed to get order", http.StatusBadRequest)
	}

	err = order.GetShippingMethod()
	if err != nil {
		util.ErrorResponse(w, "failed to get shipping method", http.StatusBadRequest)
	}

	if !order.ShippingMethod.Mapped {
		util.ErrorResponse(w, "shipping method is not mapped", http.StatusBadRequest)
		return
	}

	shippingMethodCarriers := []models.ShippingMethodCarrier{}
	err = json.Unmarshal(order.ShippingMethod.Carriers, &shippingMethodCarriers)
	if err != nil {
		util.ErrorResponse(w, "failed to unmarshal shipping method carriers", http.StatusBadRequest)
		return
	}

	carrierConnections := []models.CarrierConnection{}
	for _, carrier := range shippingMethodCarriers {
		carrierConnection, err := models.GetCarrierConnectionByID(carrier.CarrierConnectionID)
		if err != nil {
			// TODO error log
			continue
		}

		if !carrierConnection.Active {
			continue
		}

		// skip cheapest (cheapest is only for selecting a rate automatically after shopping rates)
		if carrierConnection.ShipengineCarrierID != "" {
			carrierConnections = append(carrierConnections, *carrierConnection)
		}

	}

	for _, carrierConnection := range carrierConnections {
		carrierConnectionRatesResponse := models.CarrierConnectionRatesResponse{
			CarrierConnectionID:       carrierConnection.ID,
			CarrierConnectionNickName: carrierConnection.ShipengineNickname,
			Rates:                     []models.ShippingRate{},
		}

		for _, shipengineRate := range shipEngineRateShopResponse.RateResponse.Rates {

			if shipengineRate.CarrierID != carrierConnection.ShipengineCarrierID {
				continue
			}

			totalShippingAmount := float64(0)
			if &shipengineRate.ShippingAmount != nil {
				totalShippingAmount += shipengineRate.ShippingAmount.Amount
			}

			if &shipengineRate.TaxAmount != nil {
				totalShippingAmount += shipengineRate.TaxAmount.Amount
			}

			if &shipengineRate.InsuranceAmount != nil {
				totalShippingAmount += shipengineRate.InsuranceAmount.Amount
			}

			if &shipengineRate.OtherAmount != nil {
				totalShippingAmount += shipengineRate.OtherAmount.Amount
			}

			// TODO if setting is on to add box cost to shipping amount, add it here

			shippingRate := models.ShippingRate{
				ShipengineRateID:          shipengineRate.RateID,
				ShipengineServiceType:     shipengineRate.ServiceType,
				ShipengineServiceCode:     shipengineRate.ServiceCode,
				Amount:                    totalShippingAmount,
				Currency:                  shipengineRate.ShippingAmount.Currency,
				CarrierConnectionID:       carrierConnection.ID,
				CarrierConnectionNickName: carrierConnection.ShipengineNickname,
				BoxID:                     request.BoxID,
				Weight:                    request.Weight,
				PickSessionOrderID:        pickSessionOrder.ID,
				DeliveryDays:              shipengineRate.DeliveryDays,
			}
			err = shippingRate.Create()
			if err != nil {
				util.ErrorResponse(w, "failed to create shipping rate", http.StatusBadRequest)
				return
			}

			carrierConnectionRatesResponse.Rates = append(carrierConnectionRatesResponse.Rates, models.ShippingRate{
				ID:                    shippingRate.ID,
				ShipengineServiceType: shippingRate.ShipengineServiceType,
				ShipengineServiceCode: shippingRate.ShipengineServiceCode,
				Amount:                shippingRate.Amount,
				Currency:              shippingRate.Currency,
			})
		}

		response.CarrierConnections = append(response.CarrierConnections, carrierConnectionRatesResponse)

	}

	util.JSONResponse(w, response, http.StatusOK)
}

// TODO store selected box so that we can display it on the frontend so packers know which box the selected rate is for
func ShippingSelectRate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusBadRequest)
		return
	}

	pickSessionOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid pick session order id", http.StatusBadRequest)
		return
	}

	pickSessionOrder, err := models.GetPickSessionOrderByID(pickSessionOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session order", http.StatusBadRequest)
		return
	}

	err = pickSessionOrder.GetPickSession()
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsWarehouseOwner(pickSessionOrder.PickSession.WarehouseID) {
		util.ErrorResponse(w, "cannot access pick session order", http.StatusForbidden)
		return
	}

	request := models.ShippingSelectRateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	shippingRate, err := models.GetShippingRateByID(request.ShippingRateID)
	if err != nil {
		util.ErrorResponse(w, "failed to get shipping rate", http.StatusBadRequest)
		return
	}

	if shippingRate.PickSessionOrderID != pickSessionOrder.ID {
		util.ErrorResponse(w, "Invalid Shipping Rate Selection", http.StatusBadRequest)
		return
	}

	carrierConnection, err := models.GetCarrierConnectionByID(shippingRate.CarrierConnectionID)
	if err != nil {
		util.ErrorResponse(w, "failed to get carrier connection", http.StatusBadRequest)
		return
	}

	if !carrierConnection.Active {
		util.ErrorResponse(w, "carrier connection is not active", http.StatusBadRequest)
		return
	}

	pickSessionOrder.ShippingRateID = shippingRate.ID
	err = pickSessionOrder.Update()
	if err != nil {
		util.ErrorResponse(w, "failed to update pick session order", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func ShippingPurchaseLabel(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusBadRequest)
		return
	}

	pickSessionOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid pick session order id", http.StatusBadRequest)
		return
	}

	pickSessionOrder, err := models.GetPickSessionOrderByID(pickSessionOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session order", http.StatusBadRequest)
		return
	}

	err = pickSessionOrder.GetPickSession()
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsWarehouseOwner(pickSessionOrder.PickSession.WarehouseID) {
		util.ErrorResponse(w, "cannot access pick session order", http.StatusForbidden)
		return
	}

	if !pickSessionOrder.Picked {
		util.ErrorResponse(w, "pick session order is not fully picked", http.StatusBadRequest)
		return
	}

	if pickSessionOrder.Shipped {
		util.ErrorResponse(w, "pick session order is already shipped", http.StatusBadRequest)
		return
	}

	if !pickSessionOrder.PickSession.Completed {
		util.ErrorResponse(w, "pick session is not completed", http.StatusBadRequest)
		return
	}

	order, err := pickSessionOrder.GetOrder()
	if err != nil {
		util.ErrorResponse(w, "failed to get order", http.StatusBadRequest)
		return
	}

	err = order.GetShippingMethod()
	if err != nil {
		util.ErrorResponse(w, "failed to get shipping method", http.StatusBadRequest)
		return
	}

	if !order.ShippingMethod.Mapped {
		util.ErrorResponse(w, "shipping method is not mapped", http.StatusBadRequest)
		return
	}

	err = pickSessionOrder.GetPickSessionOrderItems()
	if err != nil {
		util.ErrorResponse(w, "failed to get pick session order items", http.StatusBadRequest)
		return
	}

	request := models.ShippingPurchaseLabelRequest{}

	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	purchaseLabelFromPreviousRate := true
	if pickSessionOrder.ShippingRateID == 0 {
		purchaseLabelFromPreviousRate = false
	} else {
		err = pickSessionOrder.GetShippingRate()
		if err != nil {
			util.ErrorResponse(w, "failed to get shipping rate", http.StatusBadRequest)
			return
		}

		if pickSessionOrder.ShippingRate.BoxID != request.BoxID {
			purchaseLabelFromPreviousRate = false
		}

		if pickSessionOrder.ShippingRate.Weight != request.Weight {
			purchaseLabelFromPreviousRate = false
		}
	}

	if purchaseLabelFromPreviousRate {

		// TODO add settings for label format and layout
		shipenginePurchaseLabelFromRateResponse, err := ShipengineHandlers.PurchaseLabelFromRate(ShipengineModels.PurchaseLabelFromRateRequest{
			LabelFormat: "pdf",
			LabelLayout: "4x6",
		}, pickSessionOrder.ShippingRate.ShipengineRateID)
		if err != nil {
			util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		shipment := models.Shipment{
			QuotedCost:          pickSessionOrder.ShippingRate.Amount,
			ShipmentCost:        shipenginePurchaseLabelFromRateResponse.ShipmentCost.Amount,
			InsuranceCost:       shipenginePurchaseLabelFromRateResponse.InsuranceCost.Amount,
			TrackingNumber:      shipenginePurchaseLabelFromRateResponse.TrackingNumber,
			PackageCode:         shipenginePurchaseLabelFromRateResponse.PackageCode,
			IsReturnLabel:       false,
			MarketplaceNotified: false,
			ShippingRateID:      pickSessionOrder.ShippingRate.ID,
			CreatedByUserID:     user.ID,
			PickSessionOrderID:  pickSessionOrder.ID,
			// TODO change to s3 url
			LabelPDFURL: shipenginePurchaseLabelFromRateResponse.LabelDownload.PDF,
		}
		err = shipment.Create()
		if err != nil {
			util.ErrorResponse(w, "failed to create shipment", http.StatusBadRequest)
			return
		}

		pickSessionOrder.Shipped = true
		err = pickSessionOrder.Update()
		if err != nil {
			util.ErrorResponse(w, "failed to update pick session order", http.StatusBadRequest)
			return
		}

		// TODO create label response struct and return it
		util.JSONResponse(w, shipenginePurchaseLabelFromRateResponse, http.StatusOK)
		return

	}

	// check if the mapping is a "cheapest" option (check if cheapestServiceCode is an empty string), or if it is a set of service codes
	var cheapestServiceCode string

	shippingMethodCarriers := []models.ShippingMethodCarrier{}
	err = json.Unmarshal(order.ShippingMethod.Carriers, &shippingMethodCarriers)
	if err != nil {
		util.ErrorResponse(w, "failed to unmarshal shipping method carriers", http.StatusBadRequest)
		return
	}

	carrierConnections := []models.CarrierConnection{}
	for _, carrier := range shippingMethodCarriers {
		carrierConnection, err := models.GetCarrierConnectionByID(carrier.CarrierConnectionID)
		if err != nil {
			// TODO error log
			continue
		}

		if !carrierConnection.Active {
			continue
		}

		// check if a cheapest service code is enabled (if it is, no other service codes should be enabled, and we should shop all rates)
		if carrierConnection.ID == util.CheapestCarrierConnectionID {
			for _, service := range carrier.Services {
				if service.Enabled {
					cheapestServiceCode = service.ServiceCode
					break
				}
			}
		}

		carrierConnections = append(carrierConnections, *carrierConnection)

	}

	// Rate shop and select rate based on mapping (if cheapest, shop all rates -- determined based on if cheapestServiceCode is an empty string)
	shipengineRateShopRequest, err := ShipengineModels.ConstructRateShopRequest(pickSessionOrder, request.BoxID, request.Weight, cheapestServiceCode != "")
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	shipEngineRateShopResponse, err := ShipengineHandlers.ShopRates(*shipengineRateShopRequest)
	if err != nil {
		util.ErrorResponse(w, "failed to shop rates", http.StatusBadRequest)
		return
	}

	// Store Rates in Database and in shoppedRates variable for use in selecting cheapest rate
	shoppedRates := []models.ShippingRate{}
	for _, carrierConnection := range carrierConnections {
		for _, shipengineRate := range shipEngineRateShopResponse.RateResponse.Rates {

			if shipengineRate.CarrierID != carrierConnection.ShipengineCarrierID {
				continue
			}

			totalShippingAmount := float64(0)
			if &shipengineRate.ShippingAmount != nil {
				totalShippingAmount += shipengineRate.ShippingAmount.Amount
			}

			if &shipengineRate.TaxAmount != nil {
				totalShippingAmount += shipengineRate.TaxAmount.Amount
			}

			if &shipengineRate.InsuranceAmount != nil {
				totalShippingAmount += shipengineRate.InsuranceAmount.Amount
			}

			if &shipengineRate.OtherAmount != nil {
				totalShippingAmount += shipengineRate.OtherAmount.Amount
			}

			// TODO if setting is on to add box cost to shipping amount, add it here

			shippingRate := models.ShippingRate{
				ShipengineRateID:          shipengineRate.RateID,
				ShipengineServiceType:     shipengineRate.ServiceType,
				ShipengineServiceCode:     shipengineRate.ServiceCode,
				Amount:                    totalShippingAmount,
				Currency:                  shipengineRate.ShippingAmount.Currency,
				CarrierConnectionID:       carrierConnection.ID,
				CarrierConnectionNickName: carrierConnection.ShipengineNickname,
				BoxID:                     request.BoxID,
				Weight:                    request.Weight,
				PickSessionOrderID:        pickSessionOrder.ID,
				DeliveryDays:              shipengineRate.DeliveryDays,
			}
			err = shippingRate.Create()
			if err != nil {
				util.ErrorResponse(w, "failed to create shipping rate", http.StatusBadRequest)
				return
			}

			shoppedRates = append(shoppedRates, shippingRate)

		}

	}

	// TODO find more reliable way to read this
	var maxDeliveryDays int
	if cheapestServiceCode != "" {
		switch cheapestServiceCode {
		case "cheapest_1_day":
			maxDeliveryDays = 1
		case "cheapest_2_day":
			maxDeliveryDays = 2
		case "cheapest_3_day":
			maxDeliveryDays = 3
		case "cheapest_4_day":
			maxDeliveryDays = 4
		case "cheapest_5_day":
			maxDeliveryDays = 5
		case "cheapest_6_day":
			maxDeliveryDays = 6
		// default covers cheapest_ever as well
		default:
			maxDeliveryDays = 0
		}
	}

	var shipengineRateForLabelPurchase models.ShippingRate
	for _, shoppedRate := range shoppedRates {

		if maxDeliveryDays != 0 && shoppedRate.DeliveryDays > maxDeliveryDays {
			continue
		}

		if shipengineRateForLabelPurchase.ID == 0 {
			shipengineRateForLabelPurchase = shoppedRate
		}

		if shoppedRate.Amount < shipengineRateForLabelPurchase.Amount {
			shipengineRateForLabelPurchase = shoppedRate
		}

	}

	if shipengineRateForLabelPurchase.ID == 0 || shipengineRateForLabelPurchase.ShipengineRateID == "" {
		util.ErrorResponse(w, "failed to find rate for label purchase", http.StatusBadRequest)
		return
	}

	// TODO add settings for label format and layout
	shipenginePurchaseLabelFromRateResponse, err := ShipengineHandlers.PurchaseLabelFromRate(ShipengineModels.PurchaseLabelFromRateRequest{
		LabelFormat: "pdf",
		LabelLayout: "4x6",
	}, shipengineRateForLabelPurchase.ShipengineRateID)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	shipment := models.Shipment{
		QuotedCost:          shipengineRateForLabelPurchase.Amount,
		ShipmentCost:        shipenginePurchaseLabelFromRateResponse.ShipmentCost.Amount,
		InsuranceCost:       shipenginePurchaseLabelFromRateResponse.InsuranceCost.Amount,
		TrackingNumber:      shipenginePurchaseLabelFromRateResponse.TrackingNumber,
		PackageCode:         shipenginePurchaseLabelFromRateResponse.PackageCode,
		IsReturnLabel:       false,
		MarketplaceNotified: false,
		ShippingRateID:      shipengineRateForLabelPurchase.ID,
		CreatedByUserID:     user.ID,
		PickSessionOrderID:  pickSessionOrder.ID,
		// TODO change to s3 url
		LabelPDFURL: shipenginePurchaseLabelFromRateResponse.LabelDownload.PDF,
	}
	err = shipment.Create()
	if err != nil {
		util.ErrorResponse(w, "failed to create shipment", http.StatusBadRequest)
		return
	}

	pickSessionOrder.Shipped = true
	err = pickSessionOrder.Update()
	if err != nil {
		util.ErrorResponse(w, "failed to update pick session order", http.StatusBadRequest)
		return
	}

	// TODO create label response struct and return it
	util.JSONResponse(w, shipenginePurchaseLabelFromRateResponse, http.StatusOK)

}
