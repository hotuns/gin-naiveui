package model

type Profile struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Gender   int    `json:"gender"`
	Avatar   string `json:"avatar"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	UserId   int    `json:"userId" gorm:"column:user_id"`
	NickName string `json:"nickname"`
}

func (Profile) TableName() string {
	return "profile"
}
