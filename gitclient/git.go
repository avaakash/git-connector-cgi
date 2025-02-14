package gitclient

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	gitHTTP "github.com/go-git/go-git/v5/plumbing/transport/http"
	gitSSH "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/harness/git-connector-cgi/common"
	"github.com/sirupsen/logrus"
)

type GitClient struct {
	Repo     string
	HTTPAuth *common.HTTPAuth
	SSHAuth  *common.SSHAuth
}

func NewHttp(repo string, config *common.HTTPAuth) *GitClient {
	return &GitClient{
		Repo:     repo,
		HTTPAuth: config,
	}
}

func NewSsh(repo string, config *common.SSHAuth) *GitClient {
	return &GitClient{
		Repo:    repo,
		SSHAuth: config,
	}
}

func (gc *GitClient) ValidateWithHttp() error {
	logrus.Info("Validating repository access using HTTP token auth")
	remote := git.NewRemote(nil, &config.RemoteConfig{Name: "origin",
		URLs: []string{gc.Repo},
	})

	token, err := gc.getHttpToken()
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	_, err = remote.List(&git.ListOptions{
		Auth: &gitHTTP.BasicAuth{
			Username: gc.HTTPAuth.Username,
			Password: token,
		},
	})
	return err

}

func (gc *GitClient) ValidateWithSSH() error {
	logrus.Info("Validating repository access using SSH auth")
	remote := git.NewRemote(nil, &config.RemoteConfig{
		URLs: []string{gc.Repo},
	})

	auth, err := gc.getSSHKey()
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	_, err = remote.List(&git.ListOptions{
		Auth: auth,
	})
	return err
}

func (gc *GitClient) getHttpToken() (string, error) {
	if gc.HTTPAuth.AuthMethod == common.HTTPAuthToken {
		return gc.HTTPAuth.Token, nil
	} else if gc.HTTPAuth.AuthMethod == common.HTTPAuthPassword {
		return gc.HTTPAuth.Password, nil
	} else if gc.HTTPAuth.AuthMethod == common.HTTPAuthAnonymous {
		return "", nil
	}
	return "", fmt.Errorf("Token/Password not provided")
}

func (gc *GitClient) getSSHKey() (*gitSSH.PublicKeys, error) {
	if (gc.SSHAuth.SshKeyAuthMethod == common.SSHAuthKeyPath) && (gc.SSHAuth.SshKeyPath != "") {
		return gitSSH.NewPublicKeysFromFile(gc.SSHAuth.Username, gc.SSHAuth.SshKeyPath, gc.SSHAuth.Passphrase)
	} else if (gc.SSHAuth.SshKeyAuthMethod == common.SSHAuthKey) && (gc.SSHAuth.SshKey != nil) {
		return gitSSH.NewPublicKeys(gc.SSHAuth.Username, gc.SSHAuth.SshKey, gc.SSHAuth.Passphrase)
	}
	return nil, fmt.Errorf("SSH key not provided")
}
