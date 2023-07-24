package ClientHandlers

import "errors"

var (
	// ErrProductDoesNotBelongToClient Returned if the user's client doesn't own the requested product
	ErrProductDoesNotBelongToClient = errors.New("product does not belong to client")
	//ErrOrderDoesNotBelongToClient
	ErrOrderDoesNotBelongToClient = errors.New("order does not belong to client")
	//ErrGetClient
	ErrGetClient = errors.New("failed to get client")
	//ErrGetCarrier
	ErrGetCarrier = errors.New("failed to get carrier")
	//ErrGetCarriers
	ErrGetCarriers = errors.New("failed to get carriers")
	//ErrGetStore
	ErrGetStore = errors.New("failed to get store")
	//ErrStoreDoesNotBelongToClient
	ErrStoreDoesNotBelongToClient = errors.New("store does not belong to client")
	//ErrCarrierConnectionDoesNotBelongToClient
	ErrCarrierConnectionDoesNotBelongToClient = errors.New("carrier connection does not belong to this client")
	//ErrShipEngineDeleteCarrier
	ErrShipEngineDeleteCarrier = errors.New("failed to delete carrier from shipengine")
	//ErrShipEngineConnectCarrier
	ErrShipEngineConnectCarrier = errors.New("failed to create carrier connection in shipengine")
	//ErrCarrierConnectionBelongsToOrganization
	ErrCarrierConnectionBelongsToOrganization = errors.New("carrier connection belongs to organization, as a client user you are unable to delete")
	//ErrMarshalJSON
	ErrMarshalJSON = errors.New("failed to marshal json")
	//ErrCarrierConnectionOptions
	ErrCarrierConnectionOptions = errors.New("failed to get carrier connection options")
	//ErrGetOrders
	ErrGetOrders = errors.New("failed to get orders")
	//ErrInvalidOrderID
	ErrInvalidOrderID = errors.New("invalid order id")
	//ErrGetOrder
	ErrGetOrder = errors.New("failed to get order")
	//ErrGetOrderItems
	ErrGetOrderItems = errors.New("failed to get order items")
	//ErrGetBillToAddress
	ErrGetBillToAddress = errors.New("failed to get bill to address")
	//ErrGetShipToAddress
	ErrGetShipToAddress = errors.New("failed to get ship to address")
	//ErrGetWarehouse
	ErrGetWarehouse = errors.New("failed to get warehouse")
	//ErrGetShippingMethod
	ErrGetShippingMethod = errors.New("failed to get shipping method")
	//ErrGetOrderTags
	ErrGetOrderTags = errors.New("failed to get order tags")
	//ErrGetOrderStatus
	ErrGetOrderStatus = errors.New("failed to get order status")
)
