package models

func (OrderTags) TableName() string {
	return "orders_tags"
}

type OrderTags struct {
	OrderID int
	TagID   int
}
