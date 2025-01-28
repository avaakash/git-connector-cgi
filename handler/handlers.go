package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/harness/git-connector-cgi/common"
	gitconnector "github.com/harness/git-connector-cgi/git"
	"github.com/sirupsen/logrus"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	requestData := new(common.RequestData)
	if err := json.NewDecoder(r.Body).Decode(requestData); err != nil {
		SendErrorResponse(w, err, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	logrus.Infof("Connector Params: %+v, HTTP Params: %+v, SSH Params: %+v", requestData.ConnectorParams, requestData.ConnectorParams.HTTPAuth, requestData.ConnectorParams.SSHAuth)

	if requestData.ConnectorParams.Repo == "" {
		SendErrorResponse(w, errors.New("empty validation repository url"), "Validation repository URL is missing", http.StatusBadRequest)
		return
	}

	connector := gitconnector.New(requestData.ConnectorParams)

	authType := strings.ToLower(requestData.ConnectorParams.AuthType)
	operation := strings.ToLower(requestData.ConnectorOperation)

	var result interface{}

	switch operation {
	case "validate":
		result = connector.ValidateAcces(authType)
	default:
		SendErrorResponse(w, errors.New("invalid action"), fmt.Sprintf("The specified action %s is not supported", operation), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
