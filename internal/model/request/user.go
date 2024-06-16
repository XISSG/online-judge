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
	ID   int    `json:"id" validate:"required"`
	Type string `json:"type" validate:"oneof=password avatar"`
	Data string `json:"data"`
}
