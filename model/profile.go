package model

type Profile struct {
	ID       int    `json:"id"`
	Gender   int    `json:"gender"`
	Avatar   string `json:"avatar"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	UserId   int    `gorm:"column:user_id"`
	NickName string `gorm:"column:nick_name"`
}

func (Profile) TableName() string {
	return "profile"
}
