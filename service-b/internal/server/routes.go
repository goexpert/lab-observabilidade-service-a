package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/goexpert/lab-observabilidade-service-b/internal/dto"
	"github.com/goexpert/lab-observabilidade-service-b/pkg"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	// handleFunc is a replacement for mux.HandleFunc
	// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	handleFunc("/{cep}", s.GetTempByLocation)

	handler := otelhttp.NewHandler(mux, "/")

	return handler
}

func (s *Server) GetTempByLocation(w http.ResponseWriter, r *http.Request) {
	cepIn := r.PathValue("cep")

	if len(cepIn) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	urlTarget := fmt.Sprintf("http://viacep.com.br/ws/%s/json", cepIn)
	response, err := http.Get(urlTarget)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cepResponse dto.CedDtoOut
	err = json.NewDecoder(response.Body).Decode(&cepResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	API_KEY := "d5a4cb1d02924fda961141910240111"
	city := cepResponse.Localidade
	encodedCity := url.QueryEscape(city)
	urlTarget = fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", API_KEY, encodedCity)

	response, err = http.Get(urlTarget)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var weather dto.WeatherDTO
	err = json.NewDecoder(response.Body).Decode(&weather)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()

	resp := pkg.ResponseDTO{
		Celcius:    weather.Temperature.Celcius,
		Fahrenheit: weather.Temperature.Fahrenheit,
		Kelvin:     weather.Temperature.Celcius + float32(273),
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	code, err := w.Write(jsonResp)
	if err != nil {
		log.Println("error >>>>>>>>>>>>", code)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
