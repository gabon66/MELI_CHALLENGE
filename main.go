package main

import (
	"CHALLENGE_MELI/spaceModels"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func topsecret(w http.ResponseWriter, r *http.Request) {

	satelitesRequest := []spaceModels.Satellites{}

	reqBody := json.NewDecoder(r.Body)

	err := reqBody.Decode(&satelitesRequest)

	if err != nil {
		panic(err)
	}

	fmt.Println(satelitesRequest)

	json.NewEncoder(w).Encode(satelitesRequest)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Meli Challenge")
}

// INICIO API REST
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/topsecret", topsecret)
	log.Fatal(http.ListenAndServe(":3000", router))
}
