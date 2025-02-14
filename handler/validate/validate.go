package validate

import (
	"fmt"

	"github.com/harness/git-connector-cgi/common"
	"github.com/sirupsen/logrus"
)

func HandleValidate(provider common.Provider, config *common.GitConnectorParams) common.ValidationResponse {
	if err := handleApiAccessValidation(provider, config.APIAccess); err != nil {
		logrus.Errorf("Failed validating API access: %v", err)
		return common.ValidationResponse{
			Status:       common.Failure,
			Errors:       []common.ErrorDetail{{Message: err.Error()}},
			ErrorSummary: "Failed validating API access",
		}
	}
	if err := handleRepoAccessValidation(config); err != nil {
		logrus.Errorf("Failed validating repository access: %v", err)
		return common.ValidationResponse{
			Status:       common.Failure,
			Errors:       []common.ErrorDetail{{Message: err.Error()}},
			ErrorSummary: "Failed validating repository access",
		}
	}
	logrus.Info("Validation successful")
	return common.ValidationResponse{
		Status: common.Success,
	}
}

func handleRepoAccessValidation(config *common.GitConnectorParams) error {
	authType := config.AuthType
	switch authType {
	case common.AuthTypeHttp:
		return handleRepoAccessHttpAuthValidation(config.Repo, config.HTTPAuth)
	case common.AuthTypeSsh:
		return handleRepoAccessSshAuthValidation(config.Repo, config.SSHAuth)
	}
	logrus.Errorf("Auth type %v is not supported", authType)
	return fmt.Errorf("Auth type %v is not supported", authType)
}
