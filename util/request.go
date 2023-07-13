package util

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	ErrMissingQueryParam = errors.New("query paramter is not present")
)

func GetIntQueryParam(r *http.Request, param string) (int, error) {
	stringValue := r.URL.Query().Get(param)
	if stringValue == "" {
		return 0, ErrMissingQueryParam
	}
	value, err := strconv.Atoi(stringValue)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s: %w", param, err)
	}
	return value, nil
}

func GetStringQueryParam(r *http.Request, param string) (string, error) {
	value := r.URL.Query().Get(param)
	if value == "" {
		return "", ErrMissingQueryParam
	}
	return value, nil
}

func GetIntFromPath(r *http.Request, param string) (int, error) {
	vars := mux.Vars(r)
	stringValue := vars[param]
	if stringValue == "" {
		return 0, fmt.Errorf("%s not found in URL path", param)
	}
	value, err := strconv.Atoi(stringValue)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s: %w", param, err)
	}
	return value, nil
}

func GetStringFromPath(r *http.Request, param string) (string, error) {
	vars := mux.Vars(r)
	value := vars[param]
	if value == "" {
		return "", fmt.Errorf("%s not found in URL path", param)
	}
	return value, nil
}

func GetBoolQueryParam(r *http.Request, param string) (bool, error) {
	stringValue := r.URL.Query().Get(param)
	if stringValue == "" {
		return false, ErrMissingQueryParam
	}
	value, err := strconv.ParseBool(stringValue)
	if err != nil {
		return false, fmt.Errorf("failed to parse %s: %w", param, err)
	}
	return value, nil
}

func GetOptionalBoolQueryParam(r *http.Request, param string) (*bool, error) {
	stringValue := r.URL.Query().Get(param)
	if stringValue == "" {
		return nil, nil
	}
	value, err := strconv.ParseBool(stringValue)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", param, err)
	}
	return &value, nil
}
