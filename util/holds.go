package util

const (
	HoldNone uint64 = 1 << iota
	FraudHold
	AddressHold
	PaymentHold
	WarehouseHold
)

func HasFraudHold(value uint64) bool {
	return value&FraudHold == FraudHold
}

func HasAddressHold(value uint64) bool {
	return value&AddressHold == AddressHold
}

func HasPaymentHold(value uint64) bool {
	return value&PaymentHold == PaymentHold
}

func HasWarehouseHold(value uint64) bool {
	return value&WarehouseHold == WarehouseHold
}
