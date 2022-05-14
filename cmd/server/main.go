package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	inhttp "github.com/zvovayar/yandex-go-mustave-devops/internal/http"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

type ServerConfig struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

func main() {

	var cfg ServerConfig
	// загрузим переменные среды
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server started")
	log.Printf("Server Config environment:%+v", cfg)

	if cfg.Address != "" {
		inst.ServerAddress = cfg.Address
	}
	log.Printf("Server Strated with variables: address=%v", inst.ServerAddress)

	if cfg.StoreInterval > 0 {
		inst.StoreInterval = cfg.StoreInterval
	}
	log.Printf("Server Strated with variables: StoreInterval=%v", inst.StoreInterval)

	if len(cfg.StoreFile) > 0 {
		inst.StoreFile = cfg.StoreFile
	}
	log.Printf("Server Strated with variables: StoreFile=%v", inst.StoreFile)

	inst.Restore = cfg.Restore
	log.Printf("Server Strated with variables: Restore=%v", inst.Restore)

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

	go ListenRutine(r)

	chanOS := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(chanOS, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-chanOS
	log.Printf("INFO got a signal '%v', start shutting down...\n", sig) // put breakpoint here
	log.Printf("Shutdown complete")
}

func ListenRutine(r *chi.Mux) {
	if err := http.ListenAndServe(inst.ServerAddress, r); err != nil {
		log.Fatal(err)
	}
}
