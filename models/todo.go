package models

import "time"

type Todo struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"name" gorm:"not null"`
	Description string     `json:"surname" gorm:"not null"`
	Status      TodoStatus `json:"status" gorm:"type:VARCHAR(20)"`
	Created_at  time.Time  `json:"created_at"`
	Updated_at  time.Time  `json:"updated_at"`

	UserId uint `json:"userId" gorm:"not null"` //fk

	User User `json:"user" gorm:"foreignKey:UserId;references:ID"`
}
