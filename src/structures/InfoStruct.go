package structures

type Info struct {
	Page  int  `json:"page"`
	Total int  `json:"total"`
	Count int  `json:"count"`
	Next  bool `json:"next"`
}
