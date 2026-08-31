package utils

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
)

type JsonResponse struct {
	Success bool      `json:"success"`
	Data    any       `json:"data"`
	Error   *APIError `json:"error,omitempty"`
}

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func ReadJson(w http.ResponseWriter, r *http.Request, v any) error {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		return err
	}

	return nil
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJson(w, status, JsonResponse{
		Success: false,
		Data:    nil,
		Error: &APIError{
			Status:  status,
			Message: err.Error(),
		},
	})
}

func GenerateObjectKey(userID int64, filename, folder string) string {
	return path.Join(folder, strconv.Itoa(int(userID)), filename+".jpeg")
}
