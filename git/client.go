package gitconnector

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/harness/github-connector-cgi/common"
)

type GitConnector struct {
	Config *common.GitConfig
}

func New(config *common.GitConfig) *GitConnector {
	return &GitConnector{
		Config: config,
	}
}

func (gc *GitConnector) ValidateAcces(authMethod string) common.ValidationResponse {
	switch authMethod {
	case "token":
		return gc.validateAccesWithToken()
	case "ssh":
		return gc.validateAccesWithSSH()
	}
	return common.ValidationResponse{
		IsValid: false,
		Error: &common.Error{
			Type:    "InvalidAuthMethod",
			Message: "Invalid auth method provided",
			Reason:  "Provided auth method is invalid, only Token or SSH based auth is supported",
		},
	}
}

func (gc *GitConnector) validateAccesWithToken() common.ValidationResponse {
	remote := git.NewRemote(nil, &config.RemoteConfig{Name: "origin",
		URLs: []string{gc.Config.Repo},
	})

	if _, err := remote.List(&git.ListOptions{
		Auth: &http.BasicAuth{
			Username: gc.Config.Username,
			Password: gc.Config.Token,
		},
	}); err != nil {
		return common.ValidationResponse{
			IsValid: false,
			Error: &common.Error{
				Type:    "AccessError",
				Message: "Unable to access the specified repository",
				Reason:  err.Error(),
			},
		}
	}

	return common.ValidationResponse{
		IsValid: true,
		Error:   nil,
	}
}

func (gc *GitConnector) validateAccesWithSSH() common.ValidationResponse {
	remote := git.NewRemote(nil, &config.RemoteConfig{
		URLs: []string{gc.Config.Repo},
	})

	auth, err := ssh.NewPublicKeys(gc.Config.Username, gc.Config.SSHKey, "")
	if err != nil {
		return common.ValidationResponse{
			IsValid: false,
			Error: &common.Error{
				Type:    "SSHKeyError",
				Message: "Unable to use SSH Key",
				Reason:  err.Error(),
			},
		}
	}

	if _, err := remote.List(&git.ListOptions{
		Auth: auth,
	}); err != nil {
		return common.ValidationResponse{
			IsValid: false,
			Error: &common.Error{
				Type:    "AccessError",
				Message: "Unable to access the specified repository",
				Reason:  err.Error(),
			},
		}
	}
	return common.ValidationResponse{
		IsValid: true,
		Error:   nil,
	}
}
