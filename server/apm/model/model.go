package model

// Site holds a site data
type Site struct {
	SourceKey string `json:"sourceKey"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
}

// Asset holds an asset data
type Asset struct {
	SourceKey string `json:"sourceKey"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
}
