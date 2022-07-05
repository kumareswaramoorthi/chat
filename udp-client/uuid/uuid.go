package uuid

import (
	"errors"

	"github.com/sony/sonyflake"
)

// package uuid provides an interface for the sonyflake's
// uuid generator library.


type UuidGen interface {
	GetUUID() (int, error)
}

type uuidGen struct {
}

func NewUuidGen() UuidGen {
	return &uuidGen{}
}


//GetUUID returns a new uuid.
func (u *uuidGen) GetUUID() (int, error) {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		return 0, errors.New("failed to create uuid")
	}
	return int(id), nil
}
