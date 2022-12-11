package musicdto

import "dumbsound/models"

type MusicResponse struct {
	ID        int           `json:"id"`
	Title     string        `json:"title"`
	Year      string        `json:"year"`
	Thumbnail string        `json:"thumbnail"`
	Music     string        `json:"music"`
	Artis     models.Artist `json:"artis"`
}

type MusicResponseDelete struct {
	ID int `json:"id"`
}
