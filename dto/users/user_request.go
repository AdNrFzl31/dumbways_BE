package usersdto

type CreateUserRequest struct {
	Fullname string `json:"fullname" form:"fullname" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
