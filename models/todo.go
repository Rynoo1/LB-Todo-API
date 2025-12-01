package models

import "time"

type Todo struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"name" gorm:"not null"`
	Description string `json:"surname" gorm:"not null"`
	Status      string `json:"username" gorm:"not null"`
	Created_at  time.Time
	Updated_at  time.Time

	UserId uint `json:"userId" gorm:"not null"` //fk

	User User `json:"user" gorm:"foreignKey:UserId;references:ID"`
}
