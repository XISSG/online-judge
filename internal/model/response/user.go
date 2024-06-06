package response

type User struct {
	ID         int64  `json:"id"`
	UserName   string `json:"user_name"`
	AvatarURL  string `json:"avatar_url"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}
