package detection

type BidRequest struct {
	Site   Site   `json:"site,omitempty"`
	Device Device `json:"device,omitempty"`
}

type Site struct {
	Page string `json:"page,omitempty"` // URL of the page
}

type Device struct {
	UA string `json:"ua,omitempty"` // User agent
	IP string `json:"ip,omitempty"` // IPv4
}
