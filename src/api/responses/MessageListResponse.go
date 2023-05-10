package responses

import "github.com/kirtfieldk/astella/src/structures"

type MessageListResponse struct {
	Success bool                 `json:"success,omitempty"`
	Info    structures.Info      `json:"info"`
	Data    []structures.Message `json:"data"`
}
