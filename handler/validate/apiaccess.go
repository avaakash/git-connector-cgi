package validate

import (
	"context"
	"errors"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/harness/git-connector-cgi/common"

	"github.com/harness/git-connector-cgi/gitclient"
)

func handleApiAccessValidation(provider common.Provider, config *common.APIAccess) error {
	if config == nil {
		return nil
	}
	var (
		response    *scm.Response
		ctx         = context.Background()
		isGithubApp = config.AccessType == common.APIAccessGithubApp
		pageNumber  = 1
		pageSize    = 1
	)

	client, err := gitclient.GetGitClient(provider, config)
	if err != nil {
		return err
	}
	if isGithubApp {
		_, response, err = client.Repositories.(*github.RepositoryService).ListByInstallation(ctx, scm.ListOptions{Page: pageNumber, Size: pageSize})
	} else {
		_, response, err = client.Repositories.List(ctx, scm.ListOptions{Page: pageNumber, Size: pageSize})
	}

	if err != nil {
		return err
	}
	if response == nil || response.Status > 300 {
		return errors.New("API Access validation failed")
	}
	return nil

}

func validateAPIAccessConfig(config *common.APIAccess) error {
	if config.AccessType == "" {
		return errors.New("API Access type is missing")
	}
	if config.Token == "" {
		return errors.New("API Access token is missing")
	}
	return nil
}
