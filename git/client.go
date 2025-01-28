package gitconnector

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	gitHTTP "github.com/go-git/go-git/v5/plumbing/transport/http"
	gitSSH "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/harness/git-connector-cgi/common"
)

const (
	tokenAuth common.GitAuthType = "http"
	sshAuth   common.GitAuthType = "ssh"
)

type GitConnector struct {
	Repo     string
	HTTPAuth *common.HTTPAuth
	SSHAuth  *common.SSHAuth
}

func New(config *common.GitConnectorParams) *GitConnector {
	return &GitConnector{
		HTTPAuth: config.HTTPAuth,
		SSHAuth:  config.SSHAuth,
		Repo:     config.Repo,
	}
}

func (gc *GitConnector) ValidateAcces(authType string) common.ValidationResponse {
	switch authType {
	case string(tokenAuth):
		return gc.validateAccesWithHttp()
	case string(sshAuth):
		return gc.validateAccesWithSSH()
	}
	return common.ValidationResponse{
		Status: common.Failure,
		Errors: []common.ErrorDetail{
			{
				Reason:  "Provided auth type is invalid, only Token or SSH based auth is supported",
				Message: "InvalidAuthType",
				Code:    http.StatusBadRequest,
			},
		},
		ErrorSummary: "Invalid auth type provided",
	}
}

func (gc *GitConnector) validateAccesWithHttp() common.ValidationResponse {
	remote := git.NewRemote(nil, &config.RemoteConfig{Name: "origin",
		URLs: []string{gc.Repo},
	})

	token, err := gc.getHttpToken()
	if err != nil {
		return common.ValidationResponse{
			Status: common.Failure,
			Errors: []common.ErrorDetail{
				{
					Reason:  err.Error(),
					Message: "TokenReadError",
					Code:    http.StatusBadRequest,
				},
			},
			ErrorSummary: "Unable to read Token",
		}
	}
	if _, err := remote.List(&git.ListOptions{
		Auth: &gitHTTP.BasicAuth{
			Username: gc.HTTPAuth.Username,
			Password: token,
		},
	}); err != nil {
		return common.ValidationResponse{
			Status: common.Failure,
			Errors: []common.ErrorDetail{
				{
					Reason:  err.Error(),
					Message: "AccessError",
					Code:    http.StatusUnauthorized,
				},
			},
			ErrorSummary: "Unable to access the specified repository",
		}
	}

	return common.ValidationResponse{
		Status: common.Success,
	}
}

func (gc *GitConnector) validateAccesWithSSH() common.ValidationResponse {
	remote := git.NewRemote(nil, &config.RemoteConfig{
		URLs: []string{gc.Repo},
	})

	auth, err := gc.getSSHKey(getSshUsername(gc.Repo))
	if err != nil {
		return common.ValidationResponse{
			Status: common.Failure,
			Errors: []common.ErrorDetail{
				{
					Reason:  err.Error(),
					Message: "SSHKeyError",
					Code:    http.StatusBadRequest,
				},
			},
			ErrorSummary: "Unable to use SSH Key",
		}
	}

	if _, err := remote.List(&git.ListOptions{
		Auth: auth,
	}); err != nil {
		return common.ValidationResponse{
			Status: common.Failure,
			Errors: []common.ErrorDetail{
				{
					Reason:  err.Error(),
					Message: "AccessError",
					Code:    http.StatusUnauthorized,
				},
			},
			ErrorSummary: "Unable to access the specified repository",
		}
	}
	return common.ValidationResponse{
		Status: common.Success,
	}
}

func (gc *GitConnector) getHttpToken() (string, error) {
	if gc.HTTPAuth.AuthMethod == common.HTTPAuthToken {
		return gc.HTTPAuth.Token, nil
	} else if gc.HTTPAuth.AuthMethod == common.HTTPAuthPassword {
		return gc.HTTPAuth.Password, nil
	} else if gc.HTTPAuth.AuthMethod == common.HTTPAuthAnonymous {
		return "", nil
	}
	return "", fmt.Errorf("Token/Password not provided")
}

func (gc *GitConnector) getSSHKey(username string) (*gitSSH.PublicKeys, error) {
	if (gc.SSHAuth.SshKeyAuthMethod == common.SSHAuthKeyPath) && (gc.SSHAuth.SshKeyPath != "") {
		return gitSSH.NewPublicKeysFromFile(username, gc.SSHAuth.SshKeyPath, gc.SSHAuth.Passphrase)
	} else if (gc.SSHAuth.SshKeyAuthMethod == common.SSHAuthKey) && (gc.SSHAuth.SshKey != nil) {
		return gitSSH.NewPublicKeys(username, gc.SSHAuth.SshKey, gc.SSHAuth.Passphrase)
	}
	return nil, fmt.Errorf("SSH key not provided")
}

func getSshUsername(repoUrl string) string {
	parts := strings.Split(repoUrl, "@")
	if len(parts) > 1 {
		return parts[0]
	}
	return "git"
}
