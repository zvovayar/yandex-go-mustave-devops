package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	inhttp "github.com/zvovayar/yandex-go-mustave-devops/internal/http"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

type ServerConfig struct {
	Address string `env:"ADDRESS"`
}

func main() {

	var cfg ServerConfig
	// загрузим переменные среды
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server started")
	log.Printf("Config environment:%+v", cfg)

	if cfg.Address != "" {
		inst.ServerAddress = cfg.Address
	}

	log.Printf("Strated with variables: address=%v", inst.ServerAddress)

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

	if err := http.ListenAndServe(inst.ServerAddress, r); err != nil {
		log.Fatal(err)
	}
}
