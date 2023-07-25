package util

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

var (
	//ErrMarshalFailed
	ErrMarshalFailed = errors.New("failed to marshal object")
)

func ConvertToJSONRawMessage(object interface{}) (json.RawMessage, error) {

	data, err := json.Marshal(object)
	if err != nil {
		return nil, ErrMarshalFailed
	}

	return json.RawMessage(data), nil

}

type FieldInfo struct {
	Name string
	Type string
}

type FieldInfoJSON struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (f *FieldInfo) ConvertToReturnJSON() FieldInfoJSON {
	return FieldInfoJSON{
		Name: f.Name,
		Type: f.Type,
	}
}

func GetFieldTypes(s interface{}) []FieldInfo {
	var fieldInfos []FieldInfo
	v := reflect.ValueOf(s)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		// Extract the JSON tag and split it by the comma (if present)
		jsonTag := t.Field(i).Tag.Get("json")
		parts := strings.Split(jsonTag, ",")
		fieldName := parts[0]
		fieldType := t.Field(i).Type.Name()
		fieldInfos = append(fieldInfos, FieldInfo{Name: fieldName, Type: fieldType})
	}
	return fieldInfos
}
