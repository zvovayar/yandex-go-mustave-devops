package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zvovayar/yandex-go-mustave-devops/cmd/server/handlers"
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
	r.Get("/", handlers.GetAllMetrics)
	r.Get("/value/gauge/{GMname}", handlers.GetGMvalue)
	r.Get("/value/counter/{CMname}", handlers.GetCMvalue)

	// POST requests update
	r.Post("/update/{type}/", handlers.NotImplemented)
	r.Post("/update/{type}/{vname}/", http.NotFound)
	r.Post("/update/{type}/{vname}/{value}", handlers.NotImplemented)
	r.Post("/update/gauge/", http.NotFound)
	r.Post("/update/counter/", http.NotFound)
	r.Post("/update/gauge/{GMname}/{GMvalue}", handlers.UpdateGaugeMetric)
	r.Post("/update/counter/{CMname}/{CMvalue}", handlers.UpdateCounterMetric)

	// запуск сервера с адресом localhost, порт 8080
	//log.Fatal(http.ListenAndServe(":8080", r))
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
