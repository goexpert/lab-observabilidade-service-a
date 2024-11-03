package pkg

import (
	"encoding/json"
	"errors"
)

func ValidaCEP(cepInput json.RawMessage) error {
	if !json.Valid(cepInput) {
		return errors.New("invalid json for cep")
	}
	var cepNum float64
	if err := json.Unmarshal(cepInput, &cepNum); err == nil {
		return errors.New("invalid zipcode")
	}
	var cepStr string
	if err := json.Unmarshal(cepInput, &cepStr); err != nil {
		return errors.New("invalid zipcode")
	}
	if len(cepStr) != 8 {
		return errors.New("invalid zipcode")
	}
	return nil
}
