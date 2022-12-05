package models

import (
	"gorm.io/gorm"
	
)

type User struct {

	gorm.Model
	email    string `json:"email" gorm:"unique"`
	password string `json:"password"`
}
