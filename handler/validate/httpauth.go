package validate

import (
	"context"
	"errors"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/harness/git-connector-cgi/common"
	"github.com/harness/git-connector-cgi/gitclient"
)

func handleRepoAccessHttpAuthValidation(repo string, config *common.HTTPAuth) error {
	if config == nil {
		return errors.New("HTTP Auth is missing")
	}
	if err := validateHttpAuthConfig(config); err != nil {
		return err
	}
	if config.AuthMethod == common.HTTPAuthPassword || config.AuthMethod == common.HTTPAuthToken || config.AuthMethod == common.HTTPAuthAnonymous {
		gitClient := gitclient.NewHttp(repo, config)
		return gitClient.ValidateWithHttp()
	}
	return handleProviderSpecificAuthValidation(config)
}

func handleProviderSpecificAuthValidation(config *common.HTTPAuth) error {

	if config.AuthMethod == common.HTTPAuthGithubApp {
		if config.GithubApp == nil {
			return errors.New("Github App details not provided")
		}
		return validateGithubApp(context.Background(), &common.APIAccess{
			AccessType: common.APIAccessGithubApp,
			GithubApp:  config.GithubApp,
		})

	}
	return errors.New("Invalid HTTP Auth method")
}

func validateGithubApp(ctx context.Context, config *common.APIAccess) error {
	client, err := gitclient.GetGitClient(common.Github, config)
	if err != nil {
		return err
	}
	_, response, err := client.Repositories.(*github.RepositoryService).ListByInstallation(ctx, scm.ListOptions{Page: 1, Size: 1})

	if err != nil {
		return err
	}
	if response == nil || response.Status > 300 {
		return errors.New("Access validation failed")
	}
	return nil
}

func validateHttpAuthConfig(config *common.HTTPAuth) error {
	if config.AuthMethod == "" {
		return errors.New("HTTP Auth method is missing")
	}
	if config.AuthMethod == common.HTTPAuthToken {
		if config.Token == "" {
			return errors.New("HTTP Auth token is missing")
		}
	}
	if config.AuthMethod == common.HTTPAuthGithubApp {
		if config.GithubApp == nil {
			return errors.New("Github App is missing")
		}
		if config.GithubApp.AppId == "" {
			return errors.New("Github App ID is missing")
		}
		if config.GithubApp.AppInstallationId == "" {
			return errors.New("Github App Installation ID is missing")
		}
		if len(config.GithubApp.PrivateKey) == 0 {
			return errors.New("Github App Private Key is missing")
		}
	}
	return nil
}
