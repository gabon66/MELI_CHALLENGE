package spaceModels

type Satellite struct {
	Name     string    `json:"name"`
	Message  []string  `json:"message"`
	Distance float32   `json:"distance"`
	Coords   []float64 `json:"coords"`
}
