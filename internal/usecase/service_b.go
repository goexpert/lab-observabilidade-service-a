package usecase

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/goexpert/lab-observabilidade-service-a/internal/dto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	lab "github.com/goexpert/labobservabilidade"
)

func GetWeather(ctx context.Context, tracer trace.Tracer, cli *http.Client, z dto.CepDto) (*lab.LocaleWeatherDto, error) {

	ctx, span := tracer.Start(ctx, "clima")
	defer span.End()

	_url := "http://" + os.Getenv("SB_HOST") + ":" + os.Getenv("SB_PORT") + "/cep/" + z.Cep

	sbRequest, err := lab.NewWebclient(ctx, cli, http.MethodGet, _url, nil)
	if err != nil {
		slog.Error("erro no webclient do serviceB", "error", err.Error())
		return nil, err
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(sbRequest.Request().Header))

	var localeWeather lab.LocaleWeatherDto

	err = sbRequest.Do(func(p []byte) error {
		err = json.Unmarshal(p, &localeWeather)
		if err != nil {
			slog.Error("erro no unmarshal do body", "error", err.Error())
		}
		return err
	})
	if err != nil {
		slog.Error("erro na execução do request no serviceB", "error", err.Error())
		return nil, err
	}

	return &localeWeather, nil
}
