package models

import (
	"github.com/cespare/xxhash"
)

type Tag string

func (t Tag) ID() uint64 {
	return xxhash.Sum64String(string(t))
}
