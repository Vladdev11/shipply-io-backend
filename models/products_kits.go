package models

// override table name for inventory
func (ProductKit) TableName() string {
	return "products_kits"
}

type ProductKit struct {
	ID        int
	ProductID int
	KitID     int
	Quantity  int
}
