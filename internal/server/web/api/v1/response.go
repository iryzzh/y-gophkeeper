package v1

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

func WriteJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic("WriteJSON: " + err.Error())
	}
}
