package authdto

type LoginResponse struct {
	Email     string `gorm:"type: varchar(255)" json:"email"`
	Status    string `gorm:"type: int" json:"status"`
	Token     string `gorm:"type: varchar(255)" json:"token"`
	Subscribe string `gorm:"type: varchar(50)" json:"subscribe"`
}

type RegisterResponse struct {
	Fullname string `gorm:"type: varchar(255)" json:"fullname"`
	// Email    string `gorm:"type: varchar(255)" json:"email"`
	// Status   string `gorm:"type: varchar(1)" json:"status"`
	// Token string `gorm:"type: varchar(255)" json:"token"`
}

type CheckAuthResponse struct {
	Id        int    `gorm:"type: int" json:"id"`
	Fullname  string `gorm:"type: varchar(255)" json:"fullname"`
	Email     string `gorm:"type: varchar(255)" json:"email"`
	Gender    string `gorm:"type: varchar(255)" json:"gender"`
	Address   string `gorm:"type: varchar(255)" json:"address"`
	Phone     string `gorm:"type: varchar(255)" json:"phone"`
	Status    string `gorm:"type: varchar(25) "  json:"status"`
	Subscribe string `gorm:"type: varchar(50)" json:"subscribe"`
	Image     string `json:"image" grom:"type: varchar(255)"`
}
