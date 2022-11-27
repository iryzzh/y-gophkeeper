package v1

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

// func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
//	if IsJSONRequest(r) {
//		WriteJSON(w, data, status)
//	}
//}
//
// func IsJSONRequest(req *http.Request) bool {
//	return strings.Contains(req.Header.Get("Accept"), "application/json")
//}
//
// func WriteJSONError(w http.ResponseWriter, status int, err error) {
//	data := make(map[string]interface{})
//
//	if err != nil {
//		data["error"] = err.Error()
//	}
//
//	switch status {
//	case http.StatusNotFound:
//		data["message"] = "Not Found"
//	case http.StatusInternalServerError:
//		data["message"] = "Internal Server Error"
//	}
//
//	WriteJSON(w, data, status)
//}

func WriteJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic("WriteJSON: " + err.Error())
	}
}
