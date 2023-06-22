package ShipengineModels

type PurchaseLabelFromRateRequest struct {
	LabelFormat string `json:"label_format"`
	LabelLayout string `json:"label_layout"`
}

type PurchaseLabelFromRateResponse struct {
	LabelID         string  `json:"label_id"`
	Status          string  `json:"status"`
	ShipmentID      string  `json:"shipment_id"`
	ShipDate        string  `json:"ship_date"`
	CreatedAt       string  `json:"created_at"`
	ShipmentCost    Money   `json:"shipment_cost"`
	InsuranceCost   Money   `json:"insurance_cost"`
	TrackingNumber  string  `json:"tracking_number"`
	IsReturnLabel   bool    `json:"is_return_label"`
	IsInternational bool    `json:"is_international"`
	BatchID         string  `json:"batch_id"`
	CarrierID       string  `json:"carrier_id"`
	ServiceCode     string  `json:"service_code"`
	PackageCode     string  `json:"package_code"`
	Voided          bool    `json:"voided"`
	VoidedAt        *string `json:"voided_at"`
	LabelFormat     string  `json:"label_format"`
	LabelLayout     string  `json:"label_layout"`
	Trackable       bool    `json:"trackable"`
	CarrierCode     string  `json:"carrier_code"`
	TrackingStatus  string  `json:"tracking_status"`
	LabelDownload   Label   `json:"label_download"`
	FormDownload    *string `json:"form_download"`
	InsuranceClaim  *string `json:"insurance_claim"`
}

type Label struct {
	PDF  string `json:"pdf"`
	PNG  string `json:"png"`
	ZPL  string `json:"zpl"`
	Href string `json:"href"`
}
