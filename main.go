package main

import (
	"CHALLENGE_MELI/helpers"
	"CHALLENGE_MELI/spaceModels"
	"fmt"
	"log"
	"net/http"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/gorilla/mux"
)

var cache ttlcache.SimpleCache = ttlcache.NewCache()

/*
* GET request
 */

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Meli Challenge")
}

func processTopSecret(w http.ResponseWriter, jsonData spaceModels.Satellites) {
	var distances []float32
	var messages [][]string

	for index := range jsonData.Satellites {
		distances = append(distances, jsonData.Satellites[index].Distance)
		messages = append(messages, jsonData.Satellites[index].Message)
	}

	cordsx, cordsy := helpers.GetLocation(distances...)
	msg := helpers.GetMessage(messages...)

	if cordsx == 0 && cordsy == 0 {

		http.Error(w, "No se puede obtener posicion", 404)
		return
	}

	if len(msg) == 0 {
		http.Error(w, "No se puede obtener obtener mensaje", 404)
		return
	}

	helpers.ResponseJson(w, helpers.CreateResponseTopSecret(msg, cordsx, cordsy))
}

/*
* POST request
 */
func topsecret(w http.ResponseWriter, r *http.Request) {
	var jsonData spaceModels.Satellites

	err := helpers.DecodeJSONBody(w, r, &jsonData)
	if err != nil {
		helpers.ErroRequest(err, w)
		return
	}

	if len(jsonData.Satellites) >= 3 {

		processTopSecret(w, jsonData)

	} else {
		http.Error(w, "Satetelites insuficientes para triangular", 404)
	}
}

/*
* POST - GET request
 */
func topsecret_split(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	satellite_name := vars["satellite_name"]

	if len(satellite_name) > 0 {

		var jsonData spaceModels.Satellite

		err := helpers.DecodeJSONBody(w, r, &jsonData)
		if err != nil {
			helpers.ErroRequest(err, w)
			return
		}

		helpers.SaveAndGetDataBySatelliteName(cache, satellite_name, jsonData)

	} else {
		// si entra por aca es para validar
		var satellites []spaceModels.Satellite
		for i := 0; i < 3; i++ {
			nameSat := helpers.GetNameFromFileByIndex(i)
			currentSat, dataStored := helpers.GetDataBySatelliteName(cache, nameSat)

			// valido si tengo distancia es que esta guardado
			if dataStored {
				satellites = append(satellites, currentSat)
			}
		}

		if len(satellites) == 3 {
			var allSat spaceModels.Satellites
			allSat.Satellites = satellites
			processTopSecret(w, allSat)
		} else {
			http.Error(w, "Satetelites insuficientes para triangular", 404)
		}
	}

}

// INICIO API REST
func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", topsecret).Methods("GET")

	router.HandleFunc("/topsecret", topsecret).Methods("POST")
	router.HandleFunc("/topsecret_split", topsecret_split).Methods("GET")                   // si es get es para validar coords y mensaje
	router.HandleFunc("/topsecret_split/{satellite_name}", topsecret_split).Methods("POST") // dejo solo post para enviar data

	log.Fatal(http.ListenAndServe(":4000", router))
}
