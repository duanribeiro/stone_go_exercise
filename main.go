package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Account struct {
	id          string `json:"id"`
	name        string `json:"name"`
	cpf 		string `json:"cpf"`
	balance		float64 `json:"balance"`
	created_at	string `json:"created_at"`
}

type allAccount []Account

func main() {
	router := mux.NewRouter().StrictSlash(true)

	log.Fatal(http.ListenAndServe(":8080", router))
}