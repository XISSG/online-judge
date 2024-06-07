package entity

type User struct {
	ID           int    `gorm:"column:id;primaryKey;comment:用户id" json:"id"`                        // 用户id
	UserName     string `gorm:"column:user_name;not null;comment:用户名不允许重复" json:"user_name"`        // 用户名不允许重复
	AvatarURL    string `gorm:"column:avatar_url;comment:用户头像" json:"avatar_url"`                   // 用户头像
	UserPassword string `gorm:"column:user_password;not null;comment:用户密码" json:"user_password"`    // 用户密码
	CreateTime   string `gorm:"column:create_time;comment:创建时间" json:"create_time"`                 // 创建时间
	UpdateTime   string `gorm:"column:update_time;comment:更新时间" json:"update_time"`                 // 更新时间
	IsDelete     int32  `gorm:"column:is_delete;not null;comment:是否删除,0为不删除，1为删除" json:"is_delete"` // 是否删除,0为不删除，1为删除
	UserRole     string `gorm:"column:user_role;comment:用户类型，有user,admin,ban" json:"user_role"`
}

func (u User) TableName() string {
	return "user"
}
