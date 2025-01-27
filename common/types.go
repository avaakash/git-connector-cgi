package common

type Input struct {
	Provider        string           `json:"provider"`
	ConnectorParams *ConnectorParams `json:"connector_params"`
}

type ConnectorParams struct {
	Action     string     `json:"connector_operation"`
	AuthMethod string     `json:"auth_method"`
	Config     *GitConfig `json:"git_config"`
}

type GitConfig struct {
	Token    string `json:"token"`
	Owner    string `json:"owner"`
	Username string `json:"username"`
	Repo     string `json:"repo"`
	SSHKey   []byte `json:"ssh_key"`
}

type ValidationResponse struct {
	IsValid bool   `json:"valid"`
	Error   *Error `json:"error"`
}

type OperationStatus string

var (
	OperationStatusSuccess OperationStatus = "SUCCESS"
	OperationStatusFailure OperationStatus = "FAILURE"
)

type OperationResponse struct {
	Name            string          `json:"name"`
	Message         string          `json:"message"`
	Error           *Error          `json:"error"`
	OperationStatus OperationStatus `json:"status"`
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Status  int    `json:"status"`
}
