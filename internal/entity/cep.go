package entity

import (
	"errors"
	"log/slog"
	"regexp"

	"github.com/goexpert/lab-observabilidade-service-a/internal/dto"
)

type cepEntity struct {
	cep string
}

func NewCep(cep string) (*dto.CepDto, error) {

	var eCep = &cepEntity{
		cep: cep,
	}

	err := eCep.IsValid()
	if err != nil {
		slog.Error("[invalid cep]", "error", err.Error())
		return nil, err
	}

	return &dto.CepDto{
		Cep: eCep.cep,
	}, nil
}

func (z *cepEntity) IsValid() error {

	var re = regexp.MustCompile(`^[0-9]{8}$`)

	if !re.MatchString(z.cep) {
		return errors.New("invalid cep")
	}
	return nil
}
