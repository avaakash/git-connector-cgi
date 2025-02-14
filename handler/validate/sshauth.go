package validate

import (
	"errors"
	"fmt"

	"github.com/harness/git-connector-cgi/common"
	"github.com/harness/git-connector-cgi/gitclient"
)

func handleRepoAccessSshAuthValidation(repo string, config *common.SSHAuth) error {
	if config == nil {
		return errors.New("SSH Auth is missing")
	}
	if err := validateSshAuthConfig(config); err != nil {
		return err
	}

	gitClient := gitclient.NewSsh(repo, config)
	return gitClient.ValidateWithSSH()

}

func validateSshAuthConfig(config *common.SSHAuth) error {
	if config.AuthMechanism == common.SSHKey {
		if config.Username == "" {
			return errors.New("SSH Auth username is missing")
		}
		if config.SshKeyAuthMethod == common.SSHAuthKey && config.SshKey == nil {
			return errors.New("SSH Auth private key is missing")
		} else if config.SshKeyAuthMethod == common.SSHAuthKeyPath && config.SshKeyPath == "" {
			return errors.New("SSH Auth private key path is missing")
		} else if config.SshKeyAuthMethod == common.SSHAuthPassword && config.Password == "" {
			return errors.New("SSH Auth password is missing")
		}
	} else {
		return errors.New(fmt.Sprintf("SSH Auth mechanism %v is not supported", config.AuthMechanism))
	}

	return nil
}
