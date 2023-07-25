package util

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error  string   `json:"error,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

func JSONError(text string) map[string]string {
	return (map[string]string{"error": text})
}

func JSONErrors(errors []string) map[string][]string {
	return (map[string][]string{"errors": errors})
}

func JSONSuccess() map[string]string {
	return (map[string]string{"success": "true"})
}

func ErrResponse(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	// Handle Joined errors
	if errs, ok := err.(interface{ Unwrap() []error }); ok {
		errs := errs.Unwrap()
		strs := make([]string, len(errs))
		for i := range errs {
			strs[i] = errs[i].Error()
		}
		json.NewEncoder(w).Encode(errorResponse{Errors: strs})
		return
	}
	json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
}

func ErrorResponse(w http.ResponseWriter, err string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(JSONError(err))
}

func ErrorsResponse(w http.ResponseWriter, errs []string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(JSONErrors(errs))
}

func SuccessResponse(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(JSONSuccess())
}

func JSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func RedirectResponse(w http.ResponseWriter, r *http.Request, url string, statusCode int) {
	http.Redirect(w, r, url, statusCode)
}
