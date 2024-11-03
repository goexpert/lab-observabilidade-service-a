package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/goexpert/lab-observabilidade-service-a/internal/dto"
	"github.com/goexpert/lab-observabilidade-service-a/pkg"
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

	handleFunc("/", s.GetTemp)

	handler := otelhttp.NewHandler(mux, "/")

	return handler
}

func (s *Server) GetTemp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var dto dto.DtoCEP
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = pkg.ValidaCEP(dto.Cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	var cepStr string
	err = json.Unmarshal(dto.Cep, &cepStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	urlTarget := fmt.Sprintf("http://service-b:8080/%s", cepStr)
	response, err := http.Get(urlTarget)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var respb pkg.ResponseDTO
	err = json.NewDecoder(response.Body).Decode(&respb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(respb)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
