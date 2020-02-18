package models

import (
	u "app/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

type Account struct {
	gorm.Model
	Email string `json:"email"`
	Password string `json:"password"`
	Token string `json:"token";sql:"-"`
}

//Validate incoming user details...
func(account *Account) Validate()(map[string]interface{}, bool)  {
	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is require"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is require"), false
	}

	// Email must be unique
	temp := &Account{}

	// Check for errors and duplicate email
	err := GetDB().Table("accounts").Where("email =?", account.Email).First(temp).Error

	if err != nil || err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connected error. Please retry"), false
	}

	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user"), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (account *Account) Create() (map[string] interface{}) {
	if resp, ok := account.Validate(); !ok {
		return resp
	}

	// hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <=0 {
		return u.Message(false, "Failed to create account, connection error!")
	}

	// Create new JWT token for the newly registered account
	tk := &Token{UserId:account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	// token -> string. Only server knows this secret.
	tokenstring, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenstring

	account.Password ="" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func Login(email, password string) (map[string]interface{})  {
	account := &Account{}

	err := GetDB().Table("accounts").Where("email =?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Failed to Login, please try again!")
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))

	if err != nil || err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false,  "Invalid login credentials. Please try again")
	}

	// Worrked! Logged in
	account.Password =""

	// Create jwt token
	tk := &Token{UserId:account.ID}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	// token -> string. Only server knows this secret.
	tokenstring, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenstring // Store in response

	resp := u.Message(true, "Logged in")
	resp["account"] = account
	return resp
}
