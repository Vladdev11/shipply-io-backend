package ShipengineModels

type Dimensions struct {
	Unit   string  `json:"unit"`
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Weight struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type Money struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}
