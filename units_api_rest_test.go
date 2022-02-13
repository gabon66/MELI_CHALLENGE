package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"net/http/httptest"
	"testing"
)

//test para validar que responda si no tengo datos
func TestTopSecretWithoutDataSplit(t *testing.T) {

	req, err := http.NewRequest("GET", "/topsecret_split", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(TopSecret_Split_Check)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

// testa para validar top secret  con datos coherentes
func TestTopSecret(t *testing.T) {

	var jsonStrKenobi = []byte(`{ 
		"satellites":[
			{"name":"kenobi",
			"distance":424.26,
			"message":["", "", "", "buen", ""]
			},
			{"name":"skywalker",
			"distance":360.55,
			"message":["", "es", "", "", "test"]
			},
			{"name":"sato",
			"distance":700,
			"message":["este", "", "un", "", ""]
			}
		]
	}`)

	req1, err1 := http.NewRequest("POST", "/topsecret", bytes.NewBuffer(jsonStrKenobi))
	if err1 != nil {
		t.Fatal(err1)
	}
	req1.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TopSecret)
	handler.ServeHTTP(rr, req1)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected := `{"message": "este es  un  mensaje secretoooo","position": {"x": -199.9986,"y": 99.99285}}`

	if !strings.Contains(rr.Body.String(), "este es un buen test") {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// busco coordenaada 1 x
	if !strings.Contains(rr.Body.String(), "-199.9") {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// busco coordenaada 2 y
	if !strings.Contains(rr.Body.String(), "99.99285") {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestTopSecretFailMessage(t *testing.T) {

	var jsonStrKenobi = []byte(`{ 
		"satellites":[
			{"name":"kenobi",
			"distance":424.26,
			"message":["", "", "", "buen", ""]
			},
			{"name":"skywalker",
			"distance":360.55,
			"message":["", "", "", "", ""]
			},
			{"name":"sato",
			"distance":700,
			"message":["este", "", "un", "", ""]
			}
		]
	}`)

	req1, err1 := http.NewRequest("POST", "/topsecret", bytes.NewBuffer(jsonStrKenobi))
	if err1 != nil {
		t.Fatal(err1)
	}
	req1.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TopSecret)
	handler.ServeHTTP(rr, req1)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	expected := "No se puede obtener obtener mensaje"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	fmt.Println(rr.Body.String())
}

func TestTopSecretFailCoords(t *testing.T) {

	var jsonStrKenobi = []byte(`{ 
		"satellites":[
			{"name":"kenobi",
			"distance":1,
			"message":["", "", "", "buen", ""]
			},
			{"name":"skywalker",
			"distance":2,
			"message":["", "dd", "", "", "fff"]
			},
			{"name":"sato",
			"distance":3,
			"message":["este", "", "un", "", ""]
			}
		]
	}`)

	req1, err1 := http.NewRequest("POST", "/topsecret", bytes.NewBuffer(jsonStrKenobi))
	if err1 != nil {
		t.Fatal(err1)
	}
	req1.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TopSecret)
	handler.ServeHTTP(rr, req1)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	expected := "No se puede obtener posicion"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
