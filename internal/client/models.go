package client

type Credential struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Service struct {
	Name        string       `json:"name"`
	URL         string       `json:"url"`
	Credentials []Credential `json:"credentials"`
}

type Environment struct {
	Name     string    `json:"name"`
	Region   string    `json:"region"`
	State    string    `json:"state"`
	IP       string    `json:"ip"`
	DNSName  string    `json:"dnsName"`
	Owner    string    `json:"owner"`
	Type     string    `json:"type"`
	Services []Service `json:"services"`
}

type EnvironmentCreateRequest struct {
	Name     string `json:"name"`
	Region   string `json:"region"`
	Password string `json:"password"`
}
