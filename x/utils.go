package x

import (
	"encoding/json"
	"net/http"

	"github.com/contextgg/pkg/es"
)

func WriteJSON(w http.ResponseWriter, statusCode int, obj interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}

func WriteError(w http.ResponseWriter, err error) {
	httpErr, ok := err.(*HTTPError)
	if !ok {
		httpErr = NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(httpErr.Code)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(httpErr)
}

func WriteCommand(w http.ResponseWriter, cmd es.Command) {
	out := struct {
		AggregateId string `json:"aggregate_id"`
	}{
		AggregateId: cmd.GetAggregateId(),
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
	return
}
