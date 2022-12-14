package usersdto

type CreateUserRequest struct {
	Fullname string `json:"fullname" form:"fullname" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Fullname string `json:"fullname" form:"fullname"`
	Gender   string `json:"gender" form:"gender"`
	Phone    string `json:"phone" form:"phone"`
	Address  string `json:"address" form:"address"`
	Image    string `json:"image" form:"image"`
}
