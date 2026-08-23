// pkg/utils/response.go

package utils

import (
	"encoding/json"
	"net/http"
)

func JSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func Error(
	w http.ResponseWriter,
	status int,
	message string,
) {
	JSON(w, status, map[string]any{
		"success": false,
		"message": message,
	})
}

func Success(
	w http.ResponseWriter,
	status int,
	data any,
) {
	JSON(w, status, map[string]any{
		"success": true,
		"data":    data,
	})
}
