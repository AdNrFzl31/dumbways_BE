package authdto

type RegisterRequest struct {
	Fullname string `gorm:"type: varchar(255)" json:"fullname" validate:"required"`
	Email    string `gorm:"type: varchar(255)" json:"email" validate:"required"`
	Password string `gorm:"type: varchar(255)" json:"password" validate:"required"`
	Gender   string `gorm:"type: varchar(255)" json:"gender"`
	Status   string `gorm:"type: varchar(255)" json:"status"`
	Phone    string `gorm:"type: varchar(255)" json:"phone"`
	Address  string `gorm:"type: varchar(255)" json:"address"`
}

type LoginRequest struct {
	Email    string `gorm:"type: varchar(255)" json:"email" validate:"required"`
	Password string `gorm:"type: varchar(255)" json:"password" validate:"required"`
}
