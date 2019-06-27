package utils

import (
	"time"
)

type NullableTime struct {
	Null  bool
	Value time.Time
}

func (nt NullableTime) MarshalJSON() ([]byte, error) {
	if nt.Null {
		return []byte("null"), nil
	}

	b, err := nt.Value.MarshalJSON()
	return b, err
}
