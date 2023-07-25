package ShipengineHandlers

import "errors"

var (
	//ErrEmptyCarrierInformation
	ErrEmptyCarrierInformation = errors.New("carrier information is empty")
	//ErrDeleteCarrierFailed
	ErrDeleteCarrierFailed = errors.New("failed to delete carrier")
)
