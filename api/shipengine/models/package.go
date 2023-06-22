package ShipengineModels

type Package struct {
	PackageID          string      `json:"package_id,omitempty"`
	PackageCode        string      `json:"package_code,omitempty"`
	Name               string      `json:"name,omitempty"`
	ContentDescription string      `json:"content_description,omitempty"`
	Dimensions         *Dimensions `json:"dimensions,omitempty"`
	Weight             *Weight     `json:"weight,omitempty"`
	InsuredValue       *Money      `json:"insured_value,omitempty"`
}

type PackageResponse struct {
	PackageID   string     `json:"package_id"`
	PackageCode string     `json:"package_code"`
	Name        string     `json:"name"`
	Dimensions  Dimensions `json:"dimensions"`
	Description string     `json:"description"`
}
