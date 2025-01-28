package handler

import (
	"encoding/json"
	"net/http"

	"github.com/harness/git-connector-cgi/common"
)

func NewErrorResponse(err error, message string, code int) common.ValidationResponse {
	return common.ValidationResponse{
		Status: common.Failure,
		Errors: []common.ErrorDetail{
			{
				Reason:  err.Error(),
				Message: message,
				Code:    code,
			},
		},
		ErrorSummary: message,
	}
}

func SendErrorResponse(w http.ResponseWriter, err error, message string, status int) {
	errResp := NewErrorResponse(err, message, status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(errResp)
}
