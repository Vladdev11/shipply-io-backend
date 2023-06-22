package util

const (
	HoldNone uint64 = 1 << iota
	FraudHold
	AddressHold
	PaymentHold
	WarehouseHold
)
