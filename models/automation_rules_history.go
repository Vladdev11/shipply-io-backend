package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
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

func (arh *AutomationRulesHistory) ConvertToAutomationRulesHistoryReturnJSON(ctx context.Context) *AutomationRulesHistoryReturnJSON {
	return &AutomationRulesHistoryReturnJSON{
		ID:               arh.ID,
		AutomationRuleID: arh.AutomationRuleID,
		Note:             arh.Note,
		CreatedBy:        *arh.CreatedByUser.ConvertToReturnJSON(ctx),
		CreatedAt:        arh.CreatedAt,
	}
}

func (arh *AutomationRulesHistory) Create(ctx context.Context) error {
	result := util.DBFromContext(ctx).Create(&arh)
	return result.Error
}
