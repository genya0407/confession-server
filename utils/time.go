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

	return nt.Value.MarshalJSON(), nil
}
