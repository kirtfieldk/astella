package structures

type EventWithJsonCords struct {
	UUID        string       `json:"uuid"`
	Name        string       `json:"name"`
	Public      bool         `json:"public"`
	Duration    string       `json:"duration"`
	Created     string       `json:"created"`
	Description string       `json:"description"`
	Code        string       `json:"code"`
	Location    LocationInfo `json:"location_info"`
}
