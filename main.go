package main

import (
	"crud/server"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	//CRUD - CREATE, READ, UPDATE, DELETE

	// CREATE - POST
	// READ - GET
	// UPDATE - PUT
	// DELETE - DELETE

	router := mux.NewRouter()
	router.HandleFunc("/usuarios", server.CreateUsers).Methods(http.MethodPost)
	router.HandleFunc("/usuarios", server.GetUsers).Methods(http.MethodGet)
	router.HandleFunc("/usuarios/{id}", server.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/usuarios/{id}", server.UpdateUsers).Methods(http.MethodPut)
	router.HandleFunc("/usuarios/{id}", server.DeleteUsers).Methods(http.MethodDelete)

	fmt.Println("Escutando na porta 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

}
