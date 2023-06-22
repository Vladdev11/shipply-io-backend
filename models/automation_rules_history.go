package models

import (
	"time"

	"gorm.io/gorm"
)

type AutomationRulesHistory struct {
	ID               int
	AutomationRuleID int
	Note             string
	CreatedBy        int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	AutomationRule *AutomationRule
	CreatedByUser  *User `gorm:":CreatedBy"`
}

type AutomationRulesHistoryReturnJSON struct {
	ID               int            `json:"id"`
	AutomationRuleID int            `json:"automation_rule_id"`
	Note             string         `json:"note"`
	CreatedBy        UserReturnJSON `json:"created_by"`
	CreatedAt        time.Time      `json:"created_at"`
}

func (arh *AutomationRulesHistory) ConvertToAutomationRulesHistoryReturnJSON() *AutomationRulesHistoryReturnJSON {
	return &AutomationRulesHistoryReturnJSON{
		ID:               arh.ID,
		AutomationRuleID: arh.AutomationRuleID,
		Note:             arh.Note,
		CreatedBy:        *arh.CreatedByUser.ConvertToReturnJSON(),
		CreatedAt:        arh.CreatedAt,
	}
}

func (arh *AutomationRulesHistory) Create() error {
	result := PGDB.Create(&arh)
	return result.Error
}
