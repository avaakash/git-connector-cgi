package validate

import (
	"context"
	"errors"
	"fmt"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/harness/git-connector-cgi/common"
	"github.com/sirupsen/logrus"

	"github.com/harness/git-connector-cgi/gitclient"
)

func handleApiAccessValidation(provider common.Provider, config *common.APIAccess) error {
	if config == nil {
		return nil
	}
	logrus.Info("Validating API Access")
	var (
		response    *scm.Response
		ctx         = context.Background()
		isGithubApp = config.AccessType == common.APIAccessGithubApp
		pageNumber  = 1
		pageSize    = 1
	)

	if err := validateAPIAccessConfig(config); err != nil {
		logrus.Errorf("Invalid API access config provided: %v", err)
		return err
	}

	client, err := gitclient.GetGitClient(provider, config)
	if err != nil {
		logrus.Errorf("Failed to create git provider client: %v", err)
		return err
	}
	if isGithubApp {
		_, response, err = client.Repositories.(*github.RepositoryService).ListByInstallation(ctx, scm.ListOptions{Page: pageNumber, Size: pageSize})
	} else {
		_, response, err = client.Repositories.List(ctx, scm.ListOptions{Page: pageNumber, Size: pageSize})
	}

	if err != nil {
		logrus.Errorf("Failed to authenticate due to error: %v", err)
		return err
	}
	if response == nil || response.Status > 300 {
		err := fmt.Sprintf("Received error response from server for authentication, status code: %d", response.Status)
		logrus.Error(err)
		return errors.New(err)
	}
	return nil

}

func validateAPIAccessConfig(config *common.APIAccess) error {
	if config.AccessType == "" {
		return errors.New("API Access type is missing")
	}
	switch config.AccessType {
	case common.APIAccessToken:
		if config.Token == "" {
			return errors.New("API Access token is missing")
		}
		break
	case common.APIAccessGithubApp:
		if config.GithubApp == nil {
			return errors.New("Github App config is missing")
		}
		if config.GithubApp.AppId == "" {
			return errors.New("Github App ID is missing")
		}
		if config.GithubApp.AppInstallationId == "" {
			return errors.New("Github App private key is missing")
		}
		if config.GithubApp.PrivateKey == nil {
			return errors.New("Github App private key is missing")
		}
		break
	}
	return nil
}
