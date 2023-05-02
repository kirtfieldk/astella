package responses

import "github.com/kirtfieldk/astella/src/structures"

type EventListResponse struct {
	Info structures.Info    `json:"info"`
	Data []structures.Event `json:"data"`
}
