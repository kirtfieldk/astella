package structures

type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Created     string `json:"created"`
	Ig          string `json:"ig"`
	Description string `json:"description"`
	Twitter     string `json:"twitter"`
	TikTok      string `json:"tiktok"`
	AvatarUrl   string `json:"avatar_url"`
	ImgOne      string `json:"img_one"`
	ImgTwo      string `json:"img_two"`
	ImgThree    string `json:"img_three"`
}
