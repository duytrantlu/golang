package models

import (
	u "app/utils"
	"github.com/jinzhu/gorm"
)

type Contact struct {
	gorm.Model
	Name string `json: "name"`
	Phone string `json: "phone"`
	UserId uint `json: "user_id"` // The user that this contact belongs to
}

func Validate(contact *Contact) (map[string]interface{}, bool) {
	if contact.Name == "" {
		return u.Message(false, "Invalid name"), false
	}

	if contact.Phone == "" {
		return u.Message(false, "Invalid phone"), false
	}

	return u.Message(false, "Validate passed"), true
}

func (contact *Contact) Create() (map[string]interface{}) {
	if resp, ok := Validate(contact); !ok {
		return resp
	}

	GetDB().Create(contact)

	return u.Message(true, "Contact have been created!")
}


var GetOneContact = func(id string) (*Contact) {
	contact := &Contact{}

	err := GetDB().Table("contacts").Where("id =?", id).First(contact).Error

	if err != nil {
		return nil
	}

	return contact
}

var GetAllContact = func() ([]*Contact) {
	contacts := make([]*Contact, 0)

	err := GetDB().Table("contacts").Find(&contacts).Error

	if err != nil {
		return nil
	}

	return contacts
}