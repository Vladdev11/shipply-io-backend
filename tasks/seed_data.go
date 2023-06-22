package tasks

import (
	"encoding/json"

	"github.com/shipply-io/shipply-io-backend/models"
	"gorm.io/gorm"
)

func EnsureSeedData() {
	EnsureCarrierSeedData()
	EnsureCarrierConnectionSeedData()
}

type CarrierSeed struct {
	ID                int
	ShipEngineID      string `json:"ShipEngineID"`
	Name              string `json:"Name"`
	ThumbnailURL      string `json:"ThumbnailURL"`
	SmallThumbnailURL string `json:"SmallThumbnailURL"`
}

// EnsureCarrierSeedData populates the database with carrier seed data if carriers are missing
func EnsureCarrierSeedData() {

	carrierSeeds := []CarrierSeed{
		{
			ID:                1,
			ShipEngineID:      "asendia",
			Name:              "Asendia",
			ThumbnailURL:      "https://pdimage.petparty.co/warehance/asendia.png",
			SmallThumbnailURL: "https://pdimage.petparty.co/warehance/asendia-150px.png",
		},
		{
			ID:                2,
			ShipEngineID:      "dhl_ecommerce",
			Name:              "DHL eCommerce",
			ThumbnailURL:      "https://pdimage.petparty.co/warehance/dhl_ecommerce.png",
			SmallThumbnailURL: "https://pdimage.petparty.co/warehance/dhl_ecommerce-150px.png",
		},
		{
			ID:                3,
			ShipEngineID:      "dhl_express",
			Name:              "DHL Express",
			ThumbnailURL:      "https://pdimage.petparty.co/warehance/dhl_express.png",
			SmallThumbnailURL: "https://pdimage.petparty.co/warehance/dhl_express-150px.png",
		},
		{
			ID:                4,
			ShipEngineID:      "stamps_com",
			Name:              "Stamps.com",
			ThumbnailURL:      "https://pdimage.petparty.co/warehance/stamps.png",
			SmallThumbnailURL: "https://pdimage.petparty.co/warehance/stamps-150px.png",
		},
		{
			ID:                5,
			ShipEngineID:      "endicia",
			Name:              "Endicia",
			ThumbnailURL:      "https://pdimage.petparty.co/warehance/endicia.png",
			SmallThumbnailURL: "https://pdimage.petparty.co/warehance/endicia-150px.png",
		},
		{
			ID:                6,
			ShipEngineID:      "ontrac",
			Name:              "OnTrac",
			ThumbnailURL:      "https://pdimage.petparty.co/warehance/ontrac.png",
			SmallThumbnailURL: "https://pdimage.petparty.co/warehance/ontrac-150px.png",
		},
		{
			ID:                7,
			ShipEngineID:      "ups",
			Name:              "UPS",
			ThumbnailURL:      "https://pdimage.petparty.co/warehance/ups.png",
			SmallThumbnailURL: "https://pdimage.petparty.co/warehance/ups-150px.png",
		},
		{
			ID:                8,
			ShipEngineID:      "fedex",
			Name:              "Fedex",
			ThumbnailURL:      "",
			SmallThumbnailURL: "",
		},
		{
			ID:                1000,
			ShipEngineID:      "cheapest",
			Name:              "Cheapest",
			ThumbnailURL:      "",
			SmallThumbnailURL: "",
		},
	}

	currentCarriers, err := models.GetCarriers()
	if err != nil {
		return
	}

	for _, carrierSeed := range carrierSeeds {
		//check if carrier exists in db
		carrierExists := false
		for _, carrier := range currentCarriers {
			if carrier.ShipEngineID == carrierSeed.ShipEngineID {
				carrierExists = true
				break
			}
		}

		if !carrierExists {
			newCarrier := models.Carrier{
				ID:                carrierSeed.ID,
				ShipEngineID:      carrierSeed.ShipEngineID,
				Name:              carrierSeed.Name,
				ThumbnailURL:      carrierSeed.ThumbnailURL,
				SmallThumbnailURL: carrierSeed.SmallThumbnailURL,
			}
			err = newCarrier.Create()
			if err != nil {
				return
			}
		}
	}

}

type CarrierConnectionSeed struct {
	ID                 int
	CarrierID          int
	OwnerID            int
	OwnerType          int
	Active             bool
	ShipengineNickname string
	CarrierServices    json.RawMessage
}

func EnsureCarrierConnectionSeedData() {

	cheapestCarrierServicesJSON := `[
  {
    "name": "1 Day",
    "domestic": true,
    "carrier_id": "",
    "carrier_code": "cheapest",
    "service_code": "cheapest_1_day",
    "international": false,
    "is_multi_package_supported": false
  },
  {
    "name": "2 Day",
    "domestic": true,
    "carrier_id": "",
    "carrier_code": "cheapest",
    "service_code": "cheapest_2_day",
    "international": false,
    "is_multi_package_supported": false
  },
  {
    "name": "3 Day",
    "domestic": true,
    "carrier_id": "",
    "carrier_code": "cheapest",
    "service_code": "cheapest_3_day",
    "international": false,
    "is_multi_package_supported": false
  },
  {
    "name": "4 Day",
    "domestic": true,
    "carrier_id": "",
    "carrier_code": "cheapest",
    "service_code": "cheapest_4_day",
    "international": false,
    "is_multi_package_supported": false
  },
  {
    "name": "5 Day",
    "domestic": true,
    "carrier_id": "",
    "carrier_code": "cheapest",
    "service_code": "cheapest_5_day",
    "international": false,
    "is_multi_package_supported": false
  },
  {
    "name": "6 Day",
    "domestic": true,
    "carrier_id": "",
    "carrier_code": "cheapest",
    "service_code": "cheapest_6_day",
    "international": false,
    "is_multi_package_supported": false
  },
  {
    "name": "Ever",
    "domestic": true,
    "carrier_id": "",
    "carrier_code": "cheapest",
    "service_code": "cheapest_ever",
    "international": false,
    "is_multi_package_supported": false
  }
]
`

	carrierConnectionSeeds := []CarrierConnectionSeed{{
		ID:                 1000,
		CarrierID:          1000,
		OwnerID:            0,
		OwnerType:          0,
		Active:             true,
		ShipengineNickname: "Cheapest",
		CarrierServices:    json.RawMessage(cheapestCarrierServicesJSON),
	}}

	for _, ccSeed := range carrierConnectionSeeds {
		carrierConnectionExists, err := models.GetCarrierConnectionByID(ccSeed.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			continue
		}

		if carrierConnectionExists == nil {
			newCarrierConnection := models.CarrierConnection{
				ID:              ccSeed.ID,
				CarrierID:       ccSeed.CarrierID,
				OwnerID:         ccSeed.OwnerID,
				OwnerType:       ccSeed.OwnerType,
				Active:          ccSeed.Active,
				CarrierServices: ccSeed.CarrierServices,
			}

			newCarrierConnection.Create()

		}

	}

}
