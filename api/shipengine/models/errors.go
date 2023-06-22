package ShipengineModels

type ShipengineError struct {
	RequestID string `json:"request_id"`
	Errors    []struct {
		ErrorSource string `json:"error_source"`
		ErrorType   string `json:"error_type"`
		ErrorCode   string `json:"error_code"`
		Message     string `json:"message"`
	} `json:"errors"`
}
