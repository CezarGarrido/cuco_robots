package utils

import (
	"encoding/json"
	"net/http"
)

// respondwithJSON write json response format
func RespondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	_encode(w, code, payload)
}

// respondwithError return error message
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	_encode(w, code, map[string]interface{}{"message": msg})
}

func _encode(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
