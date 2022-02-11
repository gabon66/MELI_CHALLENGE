package main

import (
	"CHALLENGE_MELI/helpers"
	"CHALLENGE_MELI/spaceModels"

	"errors"
	"fmt"
	"log"
	"net/http"
)

type jsonData struct {
	satellites []spaceModels.Satellites
}

type sate struct {
	Name string `json:"name"`
}

type ReqTopSecretPayLoad struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Satellites []sate `json:"satellites"`
}

func topsecret(w http.ResponseWriter, r *http.Request) {

	//var jsonData jsonData
	var jsonDatatest spaceModels.Satellites

	err := helpers.DecodeJSONBody(w, r, &jsonDatatest)
	if err != nil {
		erroRequest(err, w)
		return
	}

	fmt.Println("data", jsonDatatest.Satellites[0].Message[0])

}

func erroRequest(err error, w http.ResponseWriter) {
	var mr *helpers.MalformedRequest
	if errors.As(err, &mr) {
		http.Error(w, mr.Msg, mr.Status)
	} else {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Meli Challenge")
}

// INICIO API REST
func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/topsecret", topsecret)

	log.Println("Starting server on :4000...")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
