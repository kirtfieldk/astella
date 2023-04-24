package structures

type Message struct {
	UUID      string  `json:"uuid"`
	Content   string  `json:"content"`
	Created   string  `json:"created"`
	UserId    string  `json:"user_id"`
	ParentId  string  `json:"parent_id"`
	EventId   string  `json:"event_id"`
	Upvotes   int16   `json:"up_votes"`
	Pinned    bool    `json:"pinned"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
