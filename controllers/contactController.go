package controllers

import (
	"app/models"
	u "app/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var CreateContact = func(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(uint) // Grab the id of the user that send the request
	contact := &models.Contact{}

	err := json.NewDecoder(r.Body).Decode(contact)

	if err != nil {
		u.Message(false, err.Error())
	}
	contact.UserId = userId

	respond := contact.Create()
	u.Response(w, respond)
}

var GetContact = func(w http.ResponseWriter, r *http.Request) {
	idContact, ok := mux.Vars(r)["id"]

	if !ok {
		u.Message(false, "Missing id contact")
		return
	}

	data := models.GetOneContact(idContact)
	fmt.Println("====1", data)

	respond := u.Message(true, "succeed")
	respond["data"] = data
	u.Response(w, respond)
}

var GetAllContact = func(w http.ResponseWriter, r *http.Request) {
	data := models.GetAllContact()
	respond := u.Message(true, "succeed")
	respond["data"] = data
	u.Response(w, respond)
}