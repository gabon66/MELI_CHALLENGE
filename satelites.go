

import (
	"GO/geometryModels"
	"fmt"
	"math"
)

/*func main() {
	c1 := NewCircle(-4, 4, 7.5)
	c2 := NewCircle(-3, -2, 7.3)

	//println(getDistance(2, 2, 4, 3))

	fmt.Println(Intersect(c1, c2))

}*/

func GetLocation(distances ...float32) (x, y float32) {
	fmt.Println(distances)
	return 2, 12
}

// Crea un objeto circular
func NewCircle(x, y, r float64) *geometryModels.Circle {
	return &geometryModels.Circle{geometryModels.Point{x, y}, r}
}

// Encuentra la intersección de dos círculos, el número de intersecciones puede ser 0,1,2
func Intersect(a *geometryModels.Circle, b *geometryModels.Circle) (p []geometryModels.Point) {
	dx, dy := b.X-a.X, b.Y-a.Y
	lr := a.R + b.R // suma del radio
	dr := math.Abs(a.R - b.R)
	ab := math.Sqrt(dx*dx + dy*dy) // distancia al centro
	//print(ab)
	if ab <= lr && ab > dr {
		theta1 := math.Atan(dy / dx)
		ef := lr - ab
		ao := a.R - ef/2
		theta2 := math.Acos(ao / a.R)
		theta := theta1 + theta2
		xc := a.X + a.R*math.Cos(theta)
		yc := a.Y + a.R*math.Sin(theta)
		p = append(p, geometryModels.Point{xc, yc})

		if ab < lr { // Dos intersecciones
			theta3 := math.Acos(ao / a.R)
			theta = theta3 - theta1
			xd := a.X + a.R*math.Cos(theta)
			yd := a.Y - a.R*math.Sin(theta)
			p = append(p, geometryModels.Point{xd, yd})
		}
	}
	return p
}

func getDistance(x1 float64, y1 float64, x2 float64, y2 float64) float64 {
	x := x1 - x2
	y := y1 - y2
	return math.Sqrt(x*x + y*y)
}
