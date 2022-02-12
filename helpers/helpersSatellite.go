package helpers

import (
	"CHALLENGE_MELI/geometryModels"
	"CHALLENGE_MELI/spaceModels"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
)

func GetLocation(distances ...float32) (x, y float32) {
	// armo primer radio de satelite 1
	coordsS1 := getKnownPositionByIndex(0)
	c1 := NewCircle(coordsS1[0], coordsS1[1], float64(distances[0]))

	// armo segundo radio de satelite 2
	coordsS2 := getKnownPositionByIndex(1)
	c2 := NewCircle(coordsS2[0], coordsS2[1], float64(distances[1]))

	//obtengo los puntos de interseccion en base a 2 circulos (radios de satelites en base a distancia de señal)
	pointsIntersection := Intersect(c1, c2)

	// traigo las coordenadas del tercer salite
	coordsS3 := getKnownPositionByIndex(2)

	for index := range pointsIntersection {
		// valido y redondeo por que entre las convesiones de floats algunos decimales me los cambia.
		if math.Round(getDistance(coordsS3[0], coordsS3[1], pointsIntersection[index].X, pointsIntersection[index].Y)) == math.Round(float64(distances[2])) {
			// SI ENTRA ACA ES QUE TENGO LA INTERSECTION DE LOS 3 SATELITES , O SEA LA POSICION DE LA NAVE :)
			log.Printf("coordenadas encontradas x= %#v y= %#v ", pointsIntersection[index].X, pointsIntersection[index].Y)
			return float32(pointsIntersection[index].X), float32(pointsIntersection[index].Y)
		}
	}
	return
}

func GetMessage(messages ...[]string) (msg string) {

	return ""
}

// Crea un objeto circular
func NewCircle(x, y, r float64) *geometryModels.Circle {
	return &geometryModels.Circle{geometryModels.Coord{x, y}, r}
}

// Encuentra la intersección de dos círculos, el número de intersecciones puede ser 0,1,2
func Intersect(a *geometryModels.Circle, b *geometryModels.Circle) (p []geometryModels.Coord) {
	dx, dy := b.X-a.X, b.Y-a.Y
	lr := a.R + b.R // suma del radio
	dr := math.Abs(a.R - b.R)
	ab := math.Sqrt(dx*dx + dy*dy) // distancia al centro

	if ab <= lr && ab > dr {
		theta1 := math.Atan(dy / dx)
		ef := lr - ab
		ao := a.R - ef/2
		theta2 := math.Acos(ao / a.R)
		theta := theta1 + theta2
		xc := a.X + a.R*math.Cos(theta)
		yc := a.Y + a.R*math.Sin(theta)
		p = append(p, geometryModels.Coord{xc, yc})

		if ab < lr { // Dos intersecciones
			theta3 := math.Acos(ao / a.R)
			theta = theta3 - theta1
			xd := a.X + a.R*math.Cos(theta)
			yd := a.Y - a.R*math.Sin(theta)
			p = append(p, geometryModels.Coord{xd, yd})
		}
	}
	return p
}

//obtengo la distancia entre 2 coordenadas
func getDistance(x1 float64, y1 float64, x2 float64, y2 float64) float64 {
	x := x1 - x2
	y := y1 - y2
	return math.Sqrt(x*x + y*y)
}

//obtengo ultima posicion conocida en base a index de aray
func getKnownPositionByIndex(index int) []float64 {
	return loadFileLastKnowPostition().Satellites[index].Coords
}

//obtengo data de satelites pre cargada, esto se podria traer de una db  (con mas tiempo.. )
func loadFileLastKnowPostition() spaceModels.Satellites {
	jsonFile, err := os.Open("data/last_known_position.json")
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var jsonDataFile spaceModels.Satellites
	json.Unmarshal([]byte(byteValue), &jsonDataFile)

	return jsonDataFile
}
