package structures

type Event struct {
	UUID         string       `json:"uuid"`
	Name         string       `json:"name"`
	Public       bool         `json:"public"`
	Duration     int16        `json:"duration"`
	Created      string       `json:"created"`
	Description  string       `json:"description"`
	Code         string       `json:"code"`
	LocationInfo LocationInfo `json:"location_info"`
	Expired      bool         `json:"expired"`
	EndTime      string       `json:"end_time"`
}
