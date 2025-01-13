package dto

type LogradouroDto struct {
	Cep        string `json:"cep"`
	Localidade string `json:"localidade"`
	Erro       string `json:"erro"`
}
