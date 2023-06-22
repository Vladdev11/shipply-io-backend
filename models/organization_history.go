package models

import (
	"time"

	"gorm.io/gorm"
)

type OrganizationHistory struct {
	ID             int
	OrganizationID int
	Note           string
	CreatedBy      int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Organization  *Organization
	CreatedByUser *User `gorm:"foreignKey:CreatedBy"`
}

type OrganizationHistoryReturnJSON struct {
	ID             int            `json:"id"`
	OrganizationID int            `json:"organization_id"`
	Note           string         `json:"note"`
	CreatedBy      UserReturnJSON `json:"created_by"`
	CreatedAt      time.Time      `json:"created_at"`
}

func (oh *OrganizationHistory) ConvertToOrganizationHistoryReturnJSON() *OrganizationHistoryReturnJSON {
	return &OrganizationHistoryReturnJSON{
		ID:             oh.ID,
		OrganizationID: oh.OrganizationID,
		Note:           oh.Note,
		CreatedBy:      *oh.CreatedByUser.ConvertToReturnJSON(),
		CreatedAt:      oh.CreatedAt,
	}
}

func (oh *OrganizationHistory) Create() error {
	result := PGDB.Create(&oh)
	return result.Error
}
