package main

/*
	This is the main API file. Here are shown the routes and the methods they call.
*/

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Transfer struct {
	ID string `json:"id"`
	AccountOriginId string `json:"account_origin_id"`
	AccountDestinationId string `json:"account_destination_id"`
	Amount float64 `json:"amount"`
	CreatedAt string `json:"created_at"`
}

type allTranfers []Transfer
var tranfers = allTranfers{}

type Account struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Cpf string `json:"cpf"`
	Balance float64 `json:"balance"`
	CreatedAt string `json:"created_at"`
}

type allAccounts []Account

// I will use a list to make the database function.
var accounts = allAccounts{
	{
		ID: "1",
		Name: "Pedro Bala",
		Cpf: "370.986.547-65",
		Balance: 1000,
		CreatedAt: "2020-01-01",
	},
	{
		ID: "2",
		Name: "Capitu",
		Cpf: "147.258.896-35",
		Balance: 2000,
		CreatedAt: "2020-02-02",
	},
}

/* Get the list of accounts */
func getAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}


/* Create a account */
func createAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var account Account
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(reqBody, &account); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	json.Unmarshal(reqBody, &account)
	json.NewDecoder(r.Body).Decode(&account)
	currentTime := time.Now()

	account.ID = strconv.Itoa(rand.Intn(1000000))
	account.Balance = 0
	account.CreatedAt = currentTime.Format("2006-01-02")

	accounts = append(accounts, account)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&account)
}

/* Get the account balance */
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

/*
	Transfers from one Account to another.
	If the source Account has no balance, return an appropriate error code.
*/
func postTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var transfer Transfer
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(reqBody, &transfer); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	json.Unmarshal(reqBody, &transfer)
	for _, account := range accounts {
		if account.ID == transfer.AccountOriginId {
			if account.Balance - transfer.Amount < 0 {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(400)
				w.Write([]byte("AccountOriginId insufficient funds!"))
				return
			}
		}
	}

	for _, account := range accounts {
		if account.ID == transfer.AccountDestinationId {
			account.Balance = account.Balance + transfer.Amount
		} else if account.ID == transfer.AccountOriginId {
			account.Balance = account.Balance - transfer.Amount
		}
	}
	tranfers = append(tranfers, transfer)
	json.NewEncoder(w).Encode(transfer)
}

/* Get the list of transfers */
func getTranfers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tranfers)
}


func main() {
	router := mux.NewRouter()

	// ROUTES
	router.HandleFunc("/accounts", getAccounts).Methods("GET")
	router.HandleFunc("/accounts/{id}/balance", getAccountBalance).Methods("GET")
	router.HandleFunc("/accounts", createAccount).Methods("POST")
	router.HandleFunc("/tranfers", getTranfers).Methods("GET")
	router.HandleFunc("/tranfers", postTransfer).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}