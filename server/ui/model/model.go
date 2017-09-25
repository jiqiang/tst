package model

// Asset holds an UI asset data.
type Asset struct {
	EnterpriseSourceKey        string  `json:"enterprise_source_key"`
	SiteName                   string  `json:"site_name"`
	TimeElapsedSinceLastUpdate float64 `json:"time_elapsed_since_last_update"`
}

// Assets holds UI assets data.
type Assets struct {
	Type   string  `json:"type"`
	Assets []Asset `json:"assets"`
}
