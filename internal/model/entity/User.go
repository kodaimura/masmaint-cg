package entity


type User struct {
	UserId int `db:"user_id" json:"user_id"`
	UserName string `db:"user_name" json:"user_name"`
	Password string `db:"password" json:"password"`
	CreateAt string `db:"create_at" json:"create_at"`
	UpdateAt string `db:"update_at" json:"update_at"`
}