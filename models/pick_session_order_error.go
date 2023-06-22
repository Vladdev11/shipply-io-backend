package models

type PickSessionOrderError struct {
	ID                 int64  `json:"id"`
	PickSessionOrderID int64  `json:"pick_session_order_id"`
	Error              string `json:"error"`
	CreatedAt          int64  `json:"created_at"`
	CreatedBy          int64  `json:"created_by"`
	UpdatedAt          int64  `json:"updated_at"`
	UpdatedBy          int64  `json:"updated_by"`

	PickSessionOrder PickSessionOrder `json:"pick_session_order"`
}

func (psoe *PickSessionOrderError) Create() error {
	return PGDB.Create(psoe).Error
}
