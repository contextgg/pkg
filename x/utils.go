package x

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statusCode int, obj interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}

func WriteError(w http.ResponseWriter, err error) {
	httpErr := err.(*HTTPError)
	if httpErr != nil {
		w.WriteHeader(httpErr.Code)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httpErr)
		return
	}
}
