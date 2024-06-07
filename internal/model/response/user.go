package response

type User struct {
	ID         int    `json:"id"`
	UserName   string `json:"user_name,omitempty"`
	AvatarURL  string `json:"avatar_url,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
	UpdateTime string `json:"update_time,omitempty"`
}
