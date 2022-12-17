package http

import (
	"encoding/json"
	"net/http"
)

func respondWithJson(response http.ResponseWriter, request *http.Request, code int, payload interface{}) {
	output, err := json.MarshalIndent(payload, "", " ")

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}

	response.WriteHeader(code)
	response.Write(output)
}
