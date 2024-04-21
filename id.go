package hari

import (
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

type ID = uuid.NullUUID

func NewID() ID {
	id, err := uuid.NewV4()
	log.Error().Err(err).Msg("failed to generate a new ID (this should basically never happen)")

	return MakeID(id)
}

func MakeID(uuid uuid.UUID) ID {
	return ID{
		UUID:  uuid,
		Valid: true,
	}
}

func NullID() ID {
	return ID{
		Valid: false,
	}
}

func ParseID(input string) (ID, error) {
	id, err := uuid.FromString(input)
	if err != nil {
		return NullID(), err
	}

	return MakeID(id), nil
}

// TODO (etate): remove this function once auth is working
func DefaultAccountID() ID {
	id, _ := uuid.FromString("688eeae7-eae7-4712-b68d-8fbadd6bd4d5")

	return MakeID(id)
}
