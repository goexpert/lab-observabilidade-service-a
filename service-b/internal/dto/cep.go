package dto

type CedDtoOut struct {
	Localidade string `json:"localidade"`
}

type CedDtoIn struct {
	Cep string `json:"cep"`
}
