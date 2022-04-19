package api

type Client struct {
	Host   string
	Region string

	RefreshToken string
	AccessToken  string

	UserAgent string
	BaseURL   string
}

// Meta holds the status of the request informations
type Meta struct {
	Meta struct {
		Status  int    `json:"status_code"`
		Message string `json:"error_message,omitempty"`
	} `json:"meta,omitempty"`
}
