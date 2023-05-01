package structures

type Event struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	IsPublic     bool         `json:"is_public"`
	Duration     int16        `json:"duration"`
	Created      string       `json:"created"`
	Description  string       `json:"description"`
	Code         string       `json:"code"`
	LocationInfo LocationInfo `json:"location_info"`
	Expired      bool         `json:"expired"`
	EndTime      string       `json:"end_time"`
}
