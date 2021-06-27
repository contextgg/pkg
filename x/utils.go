package x

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, obj interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}
