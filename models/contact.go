package models

import "github.com/jinzhu/gorm"

type Contact struct {
	gorm.Model
	Name string `json: "Name"`
	Phone string `json: "Phone"`
	UserId uint `json: "user_id"`
}
