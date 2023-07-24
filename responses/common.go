package responses

// ChangedByUser represents the expected response body for the ChangedByUser used by multiple endpoints
type ChangedByUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ImageURL  string `json:"image_url"`
}

// Dimensions represents the expected response body for the Dimensions used by multiple endpoints
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Weight represents the expected response body for the Weight used by multiple endpoints
type Weight struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

// Address represents the expected response body for the Address used by multiple endpoints
type Address struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Company      string `json:"company"`
	Street1      string `json:"street1"`
	Street2      string `json:"street2"`
	Street3      string `json:"street3"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	Phone        string `json:"phone"`
	EmailAddress string `json:"email_address"`
}

// FieldProperties represents the response body for the FieldProperties used by multiple endpoints
type FieldProperties struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
