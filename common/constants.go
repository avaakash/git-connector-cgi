package common

const (
	Success ResponseStatus = "SUCCESS"
	Failure ResponseStatus = "FAILURE"
	Partial ResponseStatus = "PARTIAL"
	Unknown ResponseStatus = "UNKNOWN"
	Pending ResponseStatus = "PENDING"
)

const (
	Github Provider = "Github"
)

const (
	AuthTypeHttp GitAuthType = "Http"
	AuthTypeSsh  GitAuthType = "Ssh"
)
const (
	SSHKey   SSHAuthMechanism = "SSH_KEY"
	Kerberos SSHAuthMechanism = "KERBEROS"
)

const (
	HTTPAuthPassword  HTTPAuthMethod = "UsernamePassword"
	HTTPAuthToken     HTTPAuthMethod = "UsernameToken"
	HTTPAuthAnonymous HTTPAuthMethod = "Anonymous"
	HTTPAuthGithubApp HTTPAuthMethod = "GithubApp"
)

const (
	SSHAuthPassword SSHAuthMethod = "Password"
	SSHAuthKey      SSHAuthMethod = "KeyReference"
	SSHAuthKeyPath  SSHAuthMethod = "KeyPath"
)

const (
	APIAccessToken     APIAccessType = "Token"
	APIAccessGithubApp APIAccessType = "GithubApp"
)
