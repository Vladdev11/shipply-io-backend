package models

type ColumnSearch struct {
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	ValueType string      `json:"value_type"`
}

type SearchResults struct {
	FilteredCount int           `json:"filtered_count"`
	TotalCount    int           `json:"total_count"`
	Data          []interface{} `json:"data"`
}
