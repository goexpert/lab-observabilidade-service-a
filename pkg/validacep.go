package pkg

import (
	"encoding/json"
	"errors"
)

func ValidaCEP(cepInput json.RawMessage) error {
	var cepNum float64
	if !json.Valid(cepInput) {
		return errors.New("invalid json for cep")
	}
	if err := json.Unmarshal(cepInput, &cepNum); err == nil {
		return errors.New("invalid zipcode")
	}
	return nil
}
