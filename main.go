package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)


type Account struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Cpf string `json:"cpf"`
	Balance float64 `json:"balance"`
	CreatedAt string `json:"created_at"`
}


type Transfer struct {
	ID string `json:"id"`
	AccountOriginId string `json:"account_origin_id"`
	AccountDestinationId string `json:"account_destination_id"`
	Amount float64 `json:"amount"`
	CreatedAt string `json:"created_at"`
}


type allAccounts []Account
var accounts = allAccounts{
	{
		ID: "1",
		Name: "Roger Francisco",
		Cpf: "370.986.547-65",
		Balance: 1000,
		CreatedAt: "01-01-2020",
	},
	{
		ID: "2",
		Name: "Pedro Bala",
		Cpf: "147.258.896-35",
		Balance: 2000,
		CreatedAt: "02-02-2020",
	},
}


func getAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}


func createAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var account Account
	_ = json.NewDecoder(r.Body).Decode(&account)
	account.ID = strconv.Itoa(rand.Intn(1000000))
	accounts = append(accounts, account)
	json.NewEncoder(w).Encode(&account)
}


func getAccountBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range accounts {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Account{})
}


func postTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var transfer Transfer
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &transfer)

	for _, account := range accounts {
		if account.ID == transfer.AccountDestinationId {
			account.Balance = account.Balance + transfer.Amount
		} else if account.ID == transfer.AccountOriginId {
			account.Balance = account.Balance - transfer.Amount
		}
	}
	json.NewEncoder(w).Encode(transfer)
}


func main() {
	router := mux.NewRouter()

	router.HandleFunc("/accounts", getAccounts).Methods("GET")
	router.HandleFunc("/accounts/{id}/balance", getAccountBalance).Methods("GET")
	router.HandleFunc("/accounts", createAccount).Methods("POST")

	router.HandleFunc("/tranfers", postTransfer).Methods("POST")

	http.ListenAndServe(":8000", router)
}