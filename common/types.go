package common

type ResponseStatus string
type SSHAuthMechanism string
type HTTPAuthMethod string
type SSHAuthMethod string

type RequestData struct {
	ConnectorType      string              `json:"connector_type"`
	ConnectorOperation string              `json:"connector_operation"`
	ConnectorParams    *GitConnectorParams `json:"connector_params"`
}

type GitConnectorParams struct {
	AuthType string    `json:"auth_type"`
	Repo     string    `json:"repo"`
	HTTPAuth *HTTPAuth `json:"http_auth"`
	SSHAuth  *SSHAuth  `json:"ssh_auth"`
}

type HTTPAuth struct {
	AuthMethod HTTPAuthMethod `json:"auth_method"`
	Username   string         `json:"username"`
	Token      string         `json:"token"`
	Password   string         `json:"password"`
}

type SSHAuth struct {
	AuthMechanism      SSHAuthMechanism `json:"auth_mechanism"`
	SshKeyAuthMethod   SSHAuthMethod    `json:"ssh_key_auth_method"`
	KerberosAuthMethod SSHAuthMethod    `json:"kerberos_auth_method"`
	SshKey             []byte           `json:"ssh_key"`
	SshKeyPath         string           `json:"ssh_key_path"`
	Password           string           `json:"password"`
	Passphrase         string           `json:"passphrase"`
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

type GitAuthType string
