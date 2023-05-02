package responses

import "github.com/kirtfieldk/astella/src/structures"

type UserListResponse struct {
	Info structures.Info   `json:"info"`
	Data []structures.User `json:"data"`
}
