package util

import (
	"fmt"
	"reflect"
)

func SliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func FieldExistsInSlice(slice interface{}, field string, value interface{}) bool {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < s.Len(); i++ {
		v := s.Index(i)
		f := v.FieldByName(field)
		if f.IsValid() && f.Interface() == value {
			return true
		}
	}
	return false
}

func SliceContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ConvertStructsToInterfaces(inputs interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(inputs)
	if s.Kind() != reflect.Slice {
		return nil, fmt.Errorf("ConvertStructsToInterfaces() given a non-slice type")
	}

	result := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		result[i] = s.Index(i).Interface()
	}
	return result, nil
}
