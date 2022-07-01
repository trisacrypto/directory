package models

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidOrgID = errors.New("invalid organization id")
)

func (org *Organization) Key() []byte {
	if org.Id == "" {
		return uuid.Nil[:]
	}

	key := uuid.MustParse(org.Id)
	return key[:]
}

func ParseOrgID(orgID interface{}) (uuid.UUID, error) {
	switch t := orgID.(type) {
	case string:
		return uuid.Parse(t)
	case uuid.UUID:
		return t, nil
	case []byte:
		return uuid.FromBytes(t)
	default:
		return uuid.Nil, ErrInvalidOrgID
	}
}
