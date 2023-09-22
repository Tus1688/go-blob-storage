package render

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}
