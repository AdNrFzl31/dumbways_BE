package models

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primary_key:auto_increment"`
	Fullname  string    `json:"fullname" gorm:"type: varchar(255)"`
	Email     string    `json:"email" gorm:"type: varchar(50)"`
	Password  string    `json:"password" gorm:"type: varchar(255)"`
	Status    string    `json:"status" form:"status" gorm:"type: varchar(25)"`
	Gender    string    `json:"gender" gorm:"type: varchar(50)"`
	Phone     string    `json:"phone" gorm:"type: varchar(255)"`
	Address   string    `json:"address" gorm:"type: text"`
	Subscribe string    `json:"subcribe" gorm:"type: varchar(50)"`
	CreateAt  time.Time `json:"-"`
	UpdateAt  time.Time `json:"-"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Role     string `json:"-"`
}

func (UserResponse) TableName() string {
	return "users"
}
