package spaceModels

type Satellite struct {
	name     string   `json:name`
	message  []string `json:message`
	distance float64  `json:distance`
}
