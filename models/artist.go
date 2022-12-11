package models

import "time"

type Artist struct {
	ID       int       `json:"id" gorm:"primary_key:auto_increment"`
	Name     string    `json:"name" gorm:"type: varchar(255)"`
	Old      int       `json:"old" gorm:"tipe: int"`
	Artist   string    `json:"artist" gorm:"type: varchar(255)"`
	Career   string    `json:"career" gorm:"type: varchar(255)"`
	CreateAt time.Time `json:"-"`
	UpdateAt time.Time `json:"-"`
}

type ArtistResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Old    int    `json:"old"`
	// Artist string `json:"artist"`
	// Career string `json:"career"`
}

func (ArtistResponse) TableName() string {
	return "Artist"
}
