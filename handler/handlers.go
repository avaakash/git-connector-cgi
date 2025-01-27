package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/harness/github-connector-cgi/common"
	gitconnector "github.com/harness/github-connector-cgi/git"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	in := new(common.Input)

	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		SendErrorResponse(w, err, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if in.ConnectorParams.Config == nil {
		SendErrorResponse(w, errors.New("empty config"), "Configuration is missing", http.StatusBadRequest)
		return
	}

	connector := gitconnector.New(in.ConnectorParams.Config)

	authMethod := strings.ToLower(in.ConnectorParams.AuthMethod)
	operation := strings.ToLower(in.ConnectorParams.Action)

	var result interface{}

	switch operation {
	case "validate":
		result = connector.ValidateAcces(authMethod)
	default:
		SendErrorResponse(w, errors.New("invalid action"), fmt.Sprintf("The specified action %s is not supported", operation), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
