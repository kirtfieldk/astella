package uuidtransform

import (
	"log"

	"github.com/google/uuid"
)

func StringToUuidTransform(id string) (uuid.UUID, error) {
	uuid, err := uuid.ParseBytes([]byte(id))
	if err != nil {
		log.Println(err)
		return uuid, err
	}
	return uuid, nil
}
