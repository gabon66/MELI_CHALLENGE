package main

import (
	"fmt"
	"net/http"

	"net/http/httptest"
	"testing"
)

func TestSaveTopSecretSplit(t *testing.T) {

	req, err := http.NewRequest("GET", "/topsecret_split", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(index)
	handler.ServeHTTP(rr, req)
	fmt.Println(rr.Code)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	fmt.Println(rr.Body.String())

	expected := `{"id":4,"first_name":"xyz","last_name":"pqr","email_address":"xyz@pqr.com","phone_number":"1234567890"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
