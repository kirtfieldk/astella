package structures

type LocationInfo struct {
	Id             string  `json:"id"`
	TopLeftLat     float32 `json:"top_left_lat"`
	TopLeftLon     float32 `json:"top_left_lon"`
	TopRightLat    float32 `json:"top_right_lat"`
	TopRightLon    float32 `json:"top_right_lon"`
	BottomLeftLat  float32 `json:"bottom_left_lat"`
	BottomLeftLon  float32 `json:"bottom_left_lon"`
	BottomRightLat float32 `json:"bottom_right_lat"`
	BottomRightLon float32 `json:"bottom_right_lon"`
	City           string  `json:"city"`
}
