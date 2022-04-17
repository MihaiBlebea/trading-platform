package yahoofin

type Chart struct {
	Timestamp int     `json:"timestamp"`
	Open      float64 `json:"open"`
	Close     float64 `json:"close"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Volume    int     `json:"volume"`
}
