package rpc

import (
	"time"
)

type JSONTime time.Time

func (j JSONTime) MarshalJSON() ([]byte, error) {
	t := time.Time(j)
	t = t.In(UTC8)
	return []byte(`"` + t.Format(time.RFC3339) + `"`), nil
}
