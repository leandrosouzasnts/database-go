package main

import (
	"database-go/server"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/users", server.CreatedUser).Methods(http.MethodPost)

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
