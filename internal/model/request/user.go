package request

type User struct {
	UserName     string `json:"user_name"`
	AvatarURL    string `json:"avatar_url"`
	UserPassword string `json:"user_password"`
}

type Login struct {
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_password"`
}

type UpdateUser struct {
	Type string `json:"type"`
	Body Body   `json:"data"`
}

type Body struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}
