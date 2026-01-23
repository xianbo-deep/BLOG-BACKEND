package response

type Event struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
}
