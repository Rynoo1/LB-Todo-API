package models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null"`
	Surname  string `json:"surname" gorm:"not null"`
	Username string `json:"username" gorm:"uniqueIndex;not null"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Password string `json:"-" gorm:"not null"`

	Todos []Todo `json:"todos" gorm:"foreignKey:UserId"`
}
