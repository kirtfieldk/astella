package structures

type Message struct {
	UUID     string
	Content  string
	Created  string
	UserId   string
	ParentId string
	EventId  string
	Likes    int16
	Pinned   bool
}
