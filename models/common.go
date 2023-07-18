package models

/* --------------------------------- Weight --------------------------------- */
type Weight struct {
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}

/* -------------------------------- Dimesions ------------------------------- */
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}
