package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/accounts", getAccounts).Methods("GET")
	router.HandleFunc("/accounts/{id}/balance", getAccountBalance).Methods("GET")
	router.HandleFunc("/accounts", createAccount).Methods("POST")
	router.HandleFunc("/tranfers", getTranfers).Methods("GET")
	router.HandleFunc("/tranfers", postTransfer).Methods("POST")

	return router
}

func TestGetAccountsEndpoint(t *testing.T) {
	request, _ := http.NewRequest("GET", "/accounts", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestGetAccountBalanceEndpoint(t *testing.T) {
	request, _ := http.NewRequest("GET", "/accounts/1/balance", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestCreateAccountEndpoint(t *testing.T) {
	var body = []byte(`{"name": "Harry Potter"}`)

	request, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 201, response.Code, "OK response is expected")
}

func TestWrongJsonInput(t *testing.T) {
	var body = []byte(`"!@#$%"`)

	request, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 422, response.Code, "Error message is expected")
}

