package handler

import (
	svcErr "go_rest_api_with_mysql/usecase/error"
	"net/http"
)

func addDefaultHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")
}

func ConfigureCorsHeader(w http.ResponseWriter, r *http.Request, originHost string, allowedHeaders string) {
	w.Header().Set("Access-Control-Allow-Origin", originHost)
	w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
}

func processResponseErrorStatus(w http.ResponseWriter, err error, defaultErrorCode int) {
	if reqErr, ok := err.(*svcErr.ServiceError); ok {
		switch reqErr.Code {
		case svcErr.UNAUTHORIZED_REQUEST:
			w.WriteHeader(http.StatusForbidden)
		case svcErr.NO_DATA_FOUND:
			w.WriteHeader(http.StatusNotFound)
		case svcErr.INVALID_REQUEST:
			w.WriteHeader(http.StatusBadRequest)
		case svcErr.PROCESSING_ERROR:
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(defaultErrorCode)
		}
		return
	}
	w.WriteHeader(defaultErrorCode)
}
