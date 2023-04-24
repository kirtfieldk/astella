package api

type Message struct {
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	Created  string `json:"created"`
	UserId   string `json:"userId"`
	ParentId string `json:"parentId"`
	EventId  string `json:"eventId"`
	Likes    int16  `json:"likes"`
	Pinned   bool   `json:"pinned"`
}
