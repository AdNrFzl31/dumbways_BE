package models

import "time"

type Music struct {
	ID       int       `json:"id" `
	Title    string    `json:"title" gorm:"tipe: varchar(255)"`
	Tumbnail string    `json:"tumbnail" gorm:"tipe: varchar(255)"`
	Year     int       `json:"year" gorm:"tipe: int"`
	ArtistId int       `json:"artistId"`
	Artist   Artist    `json:"artist" gorm:"constraint:OnUpdate:CASCADE, OnDelete:CASCADE"`
	Music    string    `json:"music" gorm:"tipe: varchar(255)"`
	CreateAt time.Time `json:"-"`
	UpdateAt time.Time `json:"-"`
}
