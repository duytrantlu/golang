package controllers

import (
	"app/models"
	"encoding/json"
	"net/http"
	u "app/utils"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}

	err := json.NewDecoder(r.Body).Decode(account)

	if err != nil {
		u.Message(false, err.Error())
	}
	 // save to database & return response
	 resp := account.Create()
	u.Response(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}

	err := json.NewDecoder(r.Body).Decode(account)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Response(w, resp)
}
