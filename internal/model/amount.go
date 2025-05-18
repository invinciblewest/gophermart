package model

import (
	"encoding/json"
	"math"
)

type Amount int

func (a *Amount) MarshalJSON() ([]byte, error) {
	if a == nil {
		return []byte("null"), nil
	}
	return json.Marshal(float64(*a) / 100)
}

func (a *Amount) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	*a = Amount(math.Round(f * 100))
	return nil
}
