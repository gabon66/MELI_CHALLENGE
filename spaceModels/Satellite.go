package spaceModels

type Satellite struct {
	Name     string   `json:"name"`
	Message  []string `json:"message"`
	Distance float64  `json:"distance"`
}
