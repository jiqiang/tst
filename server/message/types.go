package message

// Timer is a message type
type Timer struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

// Asset is a message type
type Asset struct {
	Name                       string `json:"name"`
	TimeElapsedSinceLastUpdate int    `json:"time_elapsed_since_last_update"`
}

// Assets is a message type
type Assets struct {
	Type   string  `json:"type"`
	Assets []Asset `json:"assets"`
}
