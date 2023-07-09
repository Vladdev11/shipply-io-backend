package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type SystemTask struct {
	ID          int64           `json:"id"`
	TaskType    int64           `json:"task_type"`
	Payload     json.RawMessage `json:"payload" gorm:"type:jsonb"`
	Status      string          `json:"status"`
	Progress    int64           `json:"progress"`
	Priority    int64           `json:"priority"`
	Errors      string          `json:"errors"`
	Processing  bool            `json:"processing"`
	HasError    bool            `json:"has_error"`
	StartedAt   time.Time       `json:"started_at"`
	CompletedAt time.Time       `json:"completed_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (st *SystemTask) Create(ctx context.Context) error {
	return util.DBFromContext(ctx).Create(st).Error
}

func (st *SystemTask) Update(ctx context.Context) error {
	return util.DBFromContext(ctx).Save(st).Error
}

func GetAllTasksForProcessing(ctx context.Context) ([]SystemTask, error) {
	var tasks []SystemTask
	err := util.DBFromContext(ctx).Where("processing IS NOT TRUE").Find(&tasks).Error
	return tasks, err
}

func (st *SystemTask) GetParameters() (map[string]interface{}, error) {
	var params map[string]interface{}
	err := json.Unmarshal(st.Payload, &params)
	if err != nil {
		return nil, err
	}

	switch st.TaskType {
	case util.PullShopifyOrdersTaskType:
		err = st.validatePullShopifyOrdersTaskParameters(params)
	case util.PullShopifyProductsTaskType:
		err = st.validatePullShopifyProductsTaskParameters(params)
	}

	return params, err
}

func (st *SystemTask) validatePullShopifyOrdersTaskParameters(params map[string]interface{}) error {

	if params["shop_name"] == nil {
		return fmt.Errorf("shop_name is required")
	} else if _, ok := params["shop_name"].(string); !ok {
		return fmt.Errorf("shop_name must be a string")
	}

	if params["access_token"] == nil {
		return fmt.Errorf("access_token is required")
	} else if _, ok := params["access_token"].(string); !ok {
		return fmt.Errorf("access_token must be a string")
	}

	if params["store_id"] == nil {
		return fmt.Errorf("store_id is required")
	} else if _, ok := params["store_id"].(float64); !ok {
		return fmt.Errorf("store_id must be a number")
	} else {
		params["store_id"] = int(params["store_id"].(float64))
	}

	return nil

}

func (st *SystemTask) validatePullShopifyProductsTaskParameters(params map[string]interface{}) error {

	if params["shop_name"] == nil {
		return fmt.Errorf("shop_name is required")
	} else if _, ok := params["shop_name"].(string); !ok {
		return fmt.Errorf("shop_name must be a string")
	}

	if params["access_token"] == nil {
		return fmt.Errorf("access_token is required")
	} else if _, ok := params["access_token"].(string); !ok {
		return fmt.Errorf("access_token must be a string")
	}

	if params["store_id"] == nil {
		return fmt.Errorf("store_id is required")
	} else if _, ok := params["store_id"].(float64); !ok {
		return fmt.Errorf("store_id must be a number")
	} else {
		params["store_id"] = int(params["store_id"].(float64))
	}

	return nil

}
