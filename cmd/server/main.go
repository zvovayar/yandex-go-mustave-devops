package main

import (
	"flag"
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
	var cfgFromFlags ServerConfig

	log.Println("Server started")

	// load environment variables
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server Config environment:%+v", cfg)

	// load flags
	flag.StringVar(&cfgFromFlags.Address, "a", inst.ServerAddress, "address to listen on")
	flag.DurationVar(&cfgFromFlags.StoreInterval, "i", inst.StoreInterval, "store interval")
	flag.StringVar(&cfgFromFlags.StoreFile, "f", inst.StoreFile, "store file")
	flag.BoolVar(&cfgFromFlags.Restore, "r", inst.Restore, "restore from file on start")
	flag.Parse()
	log.Printf("Server Config flags:%+v", cfgFromFlags)

	// assign work parameters
	if cfg.Address != "" {
		inst.ServerAddress = cfg.Address
	} else {
		inst.ServerAddress = cfgFromFlags.Address
	}
	log.Printf("Server Strated with variables: address=%v", inst.ServerAddress)

	if cfg.StoreInterval > 0 {
		inst.StoreInterval = cfg.StoreInterval
	} else {
		inst.StoreInterval = cfgFromFlags.StoreInterval
	}
	log.Printf("Server Strated with variables: StoreInterval=%v", inst.StoreInterval)

	if len(cfg.StoreFile) > 0 {
		inst.StoreFile = cfg.StoreFile
	} else {
		inst.StoreFile = cfgFromFlags.StoreFile
	}
	log.Printf("Server Strated with variables: StoreFile=%v", inst.StoreFile)

	if cfg.Restore {
		inst.Restore = true
	} else if cfgFromFlags.Restore {
		inst.Restore = true
	} else {
		inst.Restore = false
	}
	// inst.Restore = cfg.Restore && cfgFromFlags.Restore

	log.Printf("Server Strated with variables: Restore=%v", inst.Restore)

	// маршрутизация запросов обработчику
	r := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.AllowContentEncoding("deflate", "gzip"))
	r.Use(middleware.Compress(5, "application/json"))

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

	if inst.Restore {
		inst.StoreMonitor.LoadData()
	}

	// start listen http
	go ListenRutine(r)
	// start data's saver
	go inst.StoreMonitor.NewPersistanceStorage()

	chanOS := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(chanOS, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-chanOS

	log.Printf("INFO got a signal '%v', start shutting down... wait 5 seconds\n", sig) // put breakpoint here
	inst.StoreMonitor.ClosePersistanceStorage()
	<-time.After(time.Second * 5)

	log.Printf("Shutdown complete")
}

func ListenRutine(r *chi.Mux) {
	if err := http.ListenAndServe(inst.ServerAddress, r); err != nil {
		log.Fatal(err)
	}
}
