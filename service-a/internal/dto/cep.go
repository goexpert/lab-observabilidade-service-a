package dto

import "encoding/json"

type DtoCEP struct {
	Cep json.RawMessage `json:"cep"`
}
