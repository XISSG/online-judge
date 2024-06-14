package request

type User struct {
	UserName     string `json:"user_name" validate:"required,max=64"`
	AvatarURL    string `json:"avatar_url" validate:"omitempty"`
	UserPassword string `json:"user_password" validate:"required,max=128"`
}

type Login struct {
	UserName     string `json:"user_name" validate:"required"`
	UserPassword string `json:"user_password" validate:"required"`
}

type UpdateUser struct {
	Type string `json:"type" validate:"oneof=password avatar"`
	Body Body   `json:"data"`
}

type Body struct {
	ID   int    `json:"id" validate:"max=64,required"`
	Data string `json:"data"`
}
