package ShipengineModels

import "github.com/shipply-io/shipply-io-backend/models"

type CarrierConnect struct {
	Carrier interface{} `json:"carrier"`
}

type CarrierConnectResponse struct {
	CarrierID string `json:"carrier_id"`
}

type CarrierOptionsResponse struct {
	Options []CarrierOption `json:"options"`
}

type CarrierOption struct {
	Name         string      `json:"name"`
	DefaultValue interface{} `json:"default_value"`
	Description  string      `json:"description"`
}

type CarrierPackageTypesResponse struct {
	Packages []CarrierPackageType `json:"packages"`
}

type CarrierPackageType struct {
	PackageID   string `json:"package_id"`
	PackageCode string `json:"package_code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CarrierServicesResponse struct {
	Services []CarrierService `json:"services"`
}

type CarrierService struct {
	CarrierID               string `json:"carrier_id"`
	CarrierCode             string `json:"carrier_code"`
	ServiceCode             string `json:"service_code"`
	Name                    string `json:"name"`
	Domestic                bool   `json:"domestic"`
	International           bool   `json:"international"`
	IsMultiPackageSupported bool   `json:"is_multi_package_supported"`
}

func (secor *CarrierOptionsResponse) ConvertToCarrierConnectionOptions() []models.CarrierOption {

	var carrierConnectionOptions []models.CarrierOption

	for _, option := range secor.Options {

		carrierConnectionOption := models.CarrierOption{
			Name:         option.Name,
			DefaultValue: option.DefaultValue,
			Description:  option.Description,
		}

		carrierConnectionOptions = append(carrierConnectionOptions, carrierConnectionOption)

	}

	return carrierConnectionOptions

}

func (secptr *CarrierPackageTypesResponse) ConvertToCarrierConnectionPackageTypes() []models.CarrierPackageType {

	var carrierConnectionPackageTypes []models.CarrierPackageType

	for _, packageType := range secptr.Packages {

		carrierConnectionPackageType := models.CarrierPackageType{
			PackageID:   packageType.PackageID,
			PackageCode: packageType.PackageCode,
			Name:        packageType.Name,
			Description: packageType.Description,
		}

		carrierConnectionPackageTypes = append(carrierConnectionPackageTypes, carrierConnectionPackageType)

	}

	return carrierConnectionPackageTypes

}

func (secsr *CarrierServicesResponse) ConvertToCarrierConnectionServices() []models.CarrierService {

	var carrierConnectionServices []models.CarrierService

	for _, service := range secsr.Services {

		carrierConnectionService := models.CarrierService{
			CarrierID:               service.CarrierID,
			CarrierCode:             service.CarrierCode,
			ServiceCode:             service.ServiceCode,
			Name:                    service.Name,
			Domestic:                service.Domestic,
			International:           service.International,
			IsMultiPackageSupported: service.IsMultiPackageSupported,
		}

		carrierConnectionServices = append(carrierConnectionServices, carrierConnectionService)

	}

	return carrierConnectionServices

}
