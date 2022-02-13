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

// inicio cache
var cache ttlcache.SimpleCache = ttlcache.NewCache()

/*
* GET request
 */
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Meli Challenge -  Operacion ​ Fuego de Quasar")
}

// Valida en base a array de satelites coordenadas y mensaje de nave
// en caso de no poder calcular coordenadas o recuperar mensaje devuelve error.
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
func TopSecret(w http.ResponseWriter, r *http.Request) {
	var jsonData spaceModels.Satellites

	err := helpers.DecodeJSONBody(w, r, &jsonData)
	if err != nil {
		helpers.ErroRequest(err, w)
		return
	}

	if len(jsonData.Satellites) >= 3 {
		processTopSecret(w, jsonData)
	} else {
		helpers.ErrorResponseNumberSatellites(w)
	}
}

/*
* GET request
 */
func TopSecret_Split_Check(w http.ResponseWriter, r *http.Request) {
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
		helpers.ErrorResponseNumberSatellites(w)
	}
}

/*
* POST
 */
func TopSecret_Split_Save(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Nombre de satelite no provisto", 404)
	}
}

// INICIO API REST
func main() {

	//fmt.Println(helpers.GetDistance(5, 5, 5, -2))

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/topsecret", TopSecret).Methods("POST")
	router.HandleFunc("/topsecret_split", TopSecret_Split_Check).Methods("GET")                  // si es get es para validar coords y mensaje
	router.HandleFunc("/topsecret_split/{satellite_name}", TopSecret_Split_Save).Methods("POST") // dejo solo post para enviar data

	log.Fatal(http.ListenAndServe(":4000", router))
}
