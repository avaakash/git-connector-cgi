package validate

import (
	"fmt"

	"github.com/harness/git-connector-cgi/common"
)

func HandleValidate(provider common.Provider, config *common.GitConnectorParams) common.ValidationResponse {
	if err := handleApiAccessValidation(provider, config.APIAccess); err != nil {
		return common.ValidationResponse{
			Status:       common.Failure,
			Errors:       []common.ErrorDetail{{Message: err.Error()}},
			ErrorSummary: err.Error(),
		}
	}
	if err := handleRepoAccessValidation(config); err != nil {
		return common.ValidationResponse{
			Status:       common.Failure,
			Errors:       []common.ErrorDetail{{Message: err.Error()}},
			ErrorSummary: err.Error(),
		}
	}
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
	return fmt.Errorf("Auth type %v is not supported", authType)
}
