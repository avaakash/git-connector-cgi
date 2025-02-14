package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/harness/git-connector-cgi/common"
	"github.com/harness/git-connector-cgi/handler/validate"
	"github.com/sirupsen/logrus"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	request := new(common.RequestData)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		logrus.Errorf("Failed to decode request body: %v", err)
		SendErrorResponse(w, err, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if request.Params.Repo == "" {
		logrus.Error("Validation repository URL is missing")
		SendErrorResponse(w, errors.New("empty validation repository url"), "Validation repository URL is missing", http.StatusBadRequest)
		return
	}

	operation := strings.ToLower(request.Operation)

	var result interface{}

	switch operation {
	case "validate":
		result = validate.HandleValidate(request.Provider, request.Params)
	default:
		logrus.Errorf("The specified action %s is not supported", operation)
		SendErrorResponse(w, errors.New("invalid action"), fmt.Sprintf("The specified action %s is not supported", operation), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
