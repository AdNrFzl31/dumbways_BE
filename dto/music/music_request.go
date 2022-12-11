package musicdto

type MusicRequest struct {
	Title    string `json:"title" form:"title" gorm:"type: varchar(255)"`
	Tumbnail string `json:"tumbnail" form:"tumbnail" gorm:"type: varchar(255)"`
	Year     int    `json:"year" form:"year" gorm:"type: int"`
	ArtistId int    `json:"artistId" form:"artistId" gorm:"type: int"`
	Music    string `json:"music" form:"music" gorm:"type: varchar(255)"`
}
