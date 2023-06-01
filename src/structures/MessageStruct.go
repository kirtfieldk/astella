package structures

type Message struct {
	Id            string  `json:"id"`
	Content       string  `json:"content"`
	Created       string  `json:"created"`
	User          User    `json:"user"`
	ParentId      string  `json:"parent_id"`
	EventId       string  `json:"event_id"`
	Upvotes       int16   `json:"up_votes"`
	UpvotedByUser bool    `json:"upvoted_by_user"`
	Pinned        bool    `json:"pinned"`
	PinnedByUser  bool    `json:"pinned_by_user"`
	Latitude      float32 `json:"latitude"`
	Longitude     float32 `json:"longitude"`
	Replies       int     `json:"replies"`
}

type MessageRequestBody struct {
	Id        string  `json:"id"`
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
