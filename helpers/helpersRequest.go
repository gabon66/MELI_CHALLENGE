package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
)

type MalformedRequest struct {
	Status int
	Msg    string
}

type ReponsePositionFormat struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type ReponseTopsecretFormat struct {
	Message  string                `json:"message"`
	Position ReponsePositionFormat `json:"position"`
}

func CreateResponseTopSecret(message string, x float32, y float32) ReponseTopsecretFormat {
	var responseFormated ReponseTopsecretFormat
	responseFormated.Message = message
	responseFormated.Position.X = x
	responseFormated.Position.Y = y
	return responseFormated
}

func (mr *MalformedRequest) Error() string {
	return mr.Msg
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			Msg := "Content-Type header is not application/json"
			return &MalformedRequest{Status: http.StatusUnsupportedMediaType, Msg: Msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			Msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			Msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

		case errors.As(err, &unmarshalTypeError):
			Msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			Msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

		case errors.Is(err, io.EOF):
			Msg := "Request body must not be empty"
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

		case err.Error() == "http: request body too large":
			Msg := "Request body must not be larger than 1MB"
			return &MalformedRequest{Status: http.StatusRequestEntityTooLarge, Msg: Msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		Msg := "Request body must only contain a single JSON object"
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}
	}

	return nil
}

func ResponseJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func ErroRequest(err error, w http.ResponseWriter) {
	var mr *MalformedRequest
	if errors.As(err, &mr) {
		http.Error(w, mr.Msg, mr.Status)
	} else {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func ErrorResponseNumberSatellites(w http.ResponseWriter) {
	http.Error(w, "Satetelites insuficientes para triangular", 404)
}
