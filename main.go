package main

import (
	"app/auth"
	"app/controllers"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main()  {
	// initial router with mux
	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/api/contacts", controllers.GetAllContact).Methods("GET")
	router.HandleFunc("/api/contacts/{id}", controllers.GetContact).Methods("GET")

	// attach jwt auth
	router.Use(auth.JwtAuthentication)

	// Get port
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080" // localhost
	}
	fmt.Println(port)
	// listen server on port
	err := http.ListenAndServe(":" + port, router)
	if err != nil {
		panic(err)
	}

}

