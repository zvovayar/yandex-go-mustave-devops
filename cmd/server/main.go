package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	inhttp "github.com/zvovayar/yandex-go-mustave-devops/internal/http"
)

func main() {
	// маршрутизация запросов обработчику
	r := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	// GET requests
	//http.HandleFunc("/", http.NotFound)
	r.Get("/", inhttp.GetAllMetrics)
	r.Get("/value/gauge/{GMname}", inhttp.GetGMvalue)
	r.Get("/value/counter/{CMname}", inhttp.GetCMvalue)

	// POST requests update, get
	r.Post("/value", inhttp.GetMvalueJSON)
	r.Post("/update", inhttp.UpdateMetricJSON)
	r.Post("/update/{type}/", inhttp.NotImplemented)
	r.Post("/update/{type}/{vname}/", http.NotFound)
	r.Post("/update/{type}/{vname}/{value}", inhttp.NotImplemented)
	r.Post("/update/gauge/", http.NotFound)
	r.Post("/update/counter/", http.NotFound)
	r.Post("/update/gauge/{GMname}/{GMvalue}", inhttp.UpdateGaugeMetric)
	r.Post("/update/counter/{CMname}/{CMvalue}", inhttp.UpdateCounterMetric)

	// запуск сервера с адресом localhost, порт 8080
	//log.Fatal(http.ListenAndServe(":8080", r))
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
