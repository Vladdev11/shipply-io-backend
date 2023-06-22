package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Use AutomationRule struct to store automation rules in the database
type AutomationRule struct {
	ID               int64
	CreatedBy        int64
	RuleName         string
	CriteraMatchType CriteriaMatchType
	Actions          json.RawMessage `gorm:"type:jsonb"` //[]AutomationRuleAction

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type OrderClassifier struct {
	ID        int64
	CreatedBy int
	Criteria  json.RawMessage `gorm:"type:jsonb"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type CriteriaMatchType int

const (
	EveryOrder            CriteriaMatchType = 1
	MatchOrderClassifiers CriteriaMatchType = 2
)

type AutomationRuleCriteria struct{}

type AutomationRuleAction struct{}
