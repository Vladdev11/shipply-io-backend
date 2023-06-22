package models

type PickSessionOrderItem struct {
	ID                 int `json:"id" gorm:"primary_key"`
	PickSessionID      int `json:"pick_session_id"`
	PickSessionOrderID int `json:"pick_session_order_id"`
	OrderItemID        int `json:"order_item_id"`
	QuantityToPick     int `json:"quantity_to_pick"`
	QuantityPicked     int `json:"quantity_picked"`
	ProductID          int `json:"product_id"`
	LocationID         int `json:"location_id"`

	PickSession      PickSession      `json:"pick_session"`
	PickSessionOrder PickSessionOrder `json:"pick_session_order"`
	OrderItem        OrderItem        `json:"order_item"`
	Product          Product          `json:"product"`
	Location         Location         `json:"location"`
}

type PickSessionOrderItemReturnJSON struct {
	ID             int `json:"id"`
	OrderItemID    int `json:"order_item_id"`
	QuantityToPick int `json:"quantity_to_pick"`
	QuantityPicked int `json:"quantity_picked"`
	ProductID      int `json:"product_id"`

	OrderItem OrderItemReturnJSON `json:"order_item"`
	Product   *ProductReturnJSON  `json:"product"`
}

func (psoi *PickSessionOrderItem) Create() error {
	err := PGDB.Create(psoi).Error
	if err != nil {
		return err
	}
	return nil
}

func (psoi *PickSessionOrderItem) Pick() error {
	psoi.QuantityPicked++
	err := PGDB.Save(psoi).Error
	if err != nil {
		return err
	}
	return nil
}

func (psoi *PickSessionOrderItem) GetPickSessionOrder() error {

	err := PGDB.Where("id = ?", psoi.PickSessionOrderID).First(&psoi.PickSessionOrder).Error
	if err != nil {
		return err
	}

	return nil
}

func (psoi *PickSessionOrderItem) GetOrderItem() error {
	err := PGDB.Where("id = ?", psoi.OrderItemID).First(&psoi.OrderItem).Error
	if err != nil {
		return err
	}
	return nil
}

func (psoi *PickSessionOrderItem) GetProduct() error {
	err := PGDB.Where("id = ?", psoi.ProductID).First(&psoi.Product).Error
	if err != nil {
		return err
	}
	return nil
}

func (psoi *PickSessionOrderItem) ConvertToReturnJSON() PickSessionOrderItemReturnJSON {

	err := psoi.GetOrderItem()
	if err != nil {
		return PickSessionOrderItemReturnJSON{}
	}

	err = psoi.GetProduct()
	if err != nil {
		return PickSessionOrderItemReturnJSON{}
	}

	returnJSON := PickSessionOrderItemReturnJSON{
		ID:             psoi.ID,
		OrderItemID:    psoi.OrderItemID,
		QuantityToPick: psoi.QuantityToPick,
		QuantityPicked: psoi.QuantityPicked,
		ProductID:      psoi.ProductID,
		OrderItem:      psoi.OrderItem.ConvertToReturnJSON(),
		Product:        psoi.Product.ConvertToReturnJSON(),
	}

	return returnJSON
}

func GetPickSessionOrderItemsByPickSessionID(pickSessionID int) ([]PickSessionOrderItem, error) {
	var pickSessionOrderItems []PickSessionOrderItem
	err := PGDB.Where("pick_session_id = ?", pickSessionID).Find(&pickSessionOrderItems).Error
	if err != nil {
		return nil, err
	}
	return pickSessionOrderItems, nil
}
