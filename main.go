package main

import (
	"CHALLENGE_MELI/helpers"
	"CHALLENGE_MELI/spaceModels"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

/*
* POST request
 */
func topsecret(w http.ResponseWriter, r *http.Request) {
	var jsonDatatest spaceModels.Satellites

	err := helpers.DecodeJSONBody(w, r, &jsonDatatest)
	if err != nil {
		helpers.ErroRequest(err, w)
		return
	}

	if len(jsonDatatest.Satellites) >= 3 {

		var distances []float32
		var messages []string

		for index := range jsonDatatest.Satellites {
			distances = append(distances, jsonDatatest.Satellites[index].Distance)
			messages = append(messages, jsonDatatest.Satellites[index].Message...)
		}

		//s1_distance := jsonDatatest.Satellites[0].Distance
		//s2_distance := jsonDatatest.Satellites[1].Distance
		//s3_distance := jsonDatatest.Satellites[2].Distance

		cordsx, cordsy := helpers.GetLocation(distances...)
		msg := helpers.GetMessage(messages)

		if cordsx == 0 && cordsy == 0 {
			// validar que pasa si efectivamente la nave esta en estas coordenadas.
			http.Error(w, "No se puede obtener posicion", 404)
			return
		}

		helpers.ResponseJson(w, helpers.CreateResponseTopSecret(msg, cordsx, cordsy))

	} else {
		http.Error(w, "Satetelites insuficientes para triangular", 404)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Meli Challenge")
}

// INICIO API REST
func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/topsecret", topsecret).Methods("POST")

	log.Fatal(http.ListenAndServe(":4000", router))
}
