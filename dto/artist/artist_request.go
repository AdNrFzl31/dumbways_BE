package artistdto

type ArtistRequest struct {
	Name   string `json:"name" validate:"required"`
	Old    int    `json:"old"  validate:"required"`
	Artist string `json:"artist" validate:"required"`
	Career string `json:"Career"  validate:"required"`
}
