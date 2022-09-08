package stardogapi

// User reflects the content returned by the Stardog API.
type User struct {
	Name        string
	Enabled     bool         `json:"enabled"`
	Superuser   bool         `json:"superuser"`
	Roles       []string     `json:"roles"`
	Permissions []Permission `json:"permissions"`
}

// Permission reflects the content returned and required by the Stardog API.
type Permission struct {
	Action       string   `json:"action"`
	ResourceType string   `json:"resource_type"`
	Resources    []string `json:"resource"`
}
