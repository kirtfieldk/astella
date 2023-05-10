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

func ParseThreeIds(idOne string, idTwo string, idThree string) (uuid.UUID, uuid.UUID, uuid.UUID, error) {
	var one, two, three uuid.UUID
	one, err := StringToUuidTransform(idOne)
	if err != nil {
		log.Printf("Failed to be UUID: " + idOne)
		return one, two, three, err
	}
	two, err = StringToUuidTransform(idTwo)
	if err != nil {
		log.Printf("Failed to be UUID: " + idTwo)
		return one, two, three, err
	}
	three, err = StringToUuidTransform(idThree)
	if err != nil {
		log.Printf("Failed to be UUID: " + idThree)
		return one, two, three, err
	}
	return one, two, three, nil
}

func ParseTwoIds(idOne string, idTwo string) (uuid.UUID, uuid.UUID, error) {
	var one, two uuid.UUID
	one, err := StringToUuidTransform(idOne)
	if err != nil {
		log.Printf("Failed to be UUID: " + idOne)
		return one, two, err
	}
	two, err = StringToUuidTransform(idTwo)
	if err != nil {
		log.Printf("Failed to be UUID: " + idTwo)
		return one, two, err
	}
	return one, two, nil
}
