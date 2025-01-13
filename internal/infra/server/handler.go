package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/goexpert/lab-observabilidade-service-a/internal/dto"
	"github.com/goexpert/lab-observabilidade-service-a/internal/entity"
	"github.com/goexpert/lab-observabilidade-service-a/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))

	trc := otel.Tracer("weatherByCEP-tracer")

	slog.Debug("body", "r.Body", r.Body)

	var z dto.BodyDto

	err := json.NewDecoder(r.Body).Decode(&z)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&dto.DtoError{Msg: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")

	cepDto, err := entity.NewCep(z.Cep)
	if err != nil {
		code := http.StatusInternalServerError
		message := err.Error()

		if strings.Contains(strings.ToLower(err.Error()), "invalid cep") {
			code = http.StatusNotFound
			message = "cep inválido"
		}

		w.WriteHeader(code)
		json.NewEncoder(w).Encode(&dto.DtoError{Msg: message})
		return
	}

	httpClient := http.DefaultClient
	localeWeatherDto, err := usecase.GetWeather(ctx, trc, httpClient, *cepDto)
	if err != nil {
		code := http.StatusInternalServerError
		message := err.Error()

		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			code = http.StatusNotFound
			message = "cep não encontrado"
		}

		w.WriteHeader(code)
		json.NewEncoder(w).Encode(&dto.DtoError{Msg: message})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(localeWeatherDto)
}
