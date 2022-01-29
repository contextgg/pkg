package x

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statusCode int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(obj)
}

func WriteError(w http.ResponseWriter, err error) {
	httpErr, ok := err.(*HTTPError)
	if !ok {
		httpErr = NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Code)
	json.NewEncoder(w).Encode(httpErr)
}

func WriteAggregateId(w http.ResponseWriter, id string) {
	out := struct {
		AggregateId string `json:"aggregate_id"`
	}{
		AggregateId: id,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(out)
}
