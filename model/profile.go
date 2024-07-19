package model

type Profile struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Gender   int    `json:"gender"`
	Avatar   string `json:"avatar"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	UserId   int    `json:"user_id"`
	NickName string `json:"nick_name"`
}

func (Profile) TableName() string {
	return "profile"
}
