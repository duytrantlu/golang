package main

import (
	"app/auth"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main()  {
	// initial router with mux
	router := mux.NewRouter()

	// attach jwt auth
	router.Use(auth.JwtAuthentication)

	// Get port
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080" // localhost
	}

	// listen server on port
	err := http.ListenAndServe(":" + port, router)
	if err != nil {
		panic(err)
	}

}

