package common

type ResponseStatus string
type Provider string
type GitAuthType string
type SSHAuthMechanism string
type HTTPAuthMethod string
type SSHAuthMethod string
type APIAccessType string

type RequestData struct {
	Provider  Provider            `json:"connector_type"`
	Operation string              `json:"connector_operation"`
	Params    *GitConnectorParams `json:"connector_params"`
}

type GitConnectorParams struct {
	AuthType  GitAuthType `json:"auth_type"`
	Repo      string      `json:"repo"`
	HTTPAuth  *HTTPAuth   `json:"http_auth"`
	SSHAuth   *SSHAuth    `json:"ssh_auth"`
	APIAccess *APIAccess  `json:"api_access"`
}

type HTTPAuth struct {
	AuthMethod HTTPAuthMethod `json:"auth_method"`
	Username   string         `json:"username"`
	Token      string         `json:"token"`
	Password   string         `json:"password"`
	GithubApp  *GithubApp     `json:"github_app_auth"`
}

type GithubApp struct {
	AppInstallationId string `json:"app_installation_id"`
	AppId             string `json:"app_id"`
	PrivateKey        []byte `json:"private_key"`
	GithubUrl         string `json:"github_url"`
}

type SSHAuth struct {
	AuthMechanism      SSHAuthMechanism `json:"auth_mechanism"`
	SshKeyAuthMethod   SSHAuthMethod    `json:"ssh_key_auth_method"`
	KerberosAuthMethod SSHAuthMethod    `json:"kerberos_auth_method"`
	Username           string           `json:"username"`
	SshKey             []byte           `json:"ssh_key"`
	SshKeyPath         string           `json:"ssh_key_path"`
	Password           string           `json:"password"`
	Passphrase         string           `json:"passphrase"`
}

type APIAccess struct {
	AccessType APIAccessType `json:"access_type"`
	Endpoint   string        `json:"endpoint"`
	ProxyURL   string        `json:"proxy_url"`
	Token      string        `json:"token"`
	GithubApp  *GithubApp    `json:"github_app"`
}
type ValidationResponse struct {
	Status       ResponseStatus `json:"status"`
	Errors       []ErrorDetail  `json:"errors"`
	ErrorSummary string         `json:"error_summary"`
}

type ErrorDetail struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
