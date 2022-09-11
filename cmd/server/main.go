// main start point for server
// environment variables:
// ADDRESS=0.0.0.0:44751
// REPORT_INTERVAL=10s
// POLL_INTERVAL=2s
// RESTORE=false
// SHUTDOWN_TIMEOUT=5s
// STORE_INTERVAL=10s
// STORE_FILE=./devops-metrics-db.json
// RESTORE=true
// KEY=secret
// DATABASE_DSN=postgres://postgres:password@localhost:5432/dbname?sslmode=disable
// flags:
// -a "address to listen on"
// -k "key for hash calculate"
// -i "store interval"
// -f "store file"
// -r "restore from file on start"
// -d "Database DSN"
package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zvovayar/yandex-go-mustave-devops/internal/config"
	inhttp "github.com/zvovayar/yandex-go-mustave-devops/internal/http"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
	"go.uber.org/zap"
)

// profiler port
const pprofAddr = ":8083"

// @Title BucketStorage API
// @Description Сервис хранения данных bucket-ов.
// @Version 1.0

// @Contact.email support@ultimatestore.io

// @BasePath /api/v1
// @Host ultimatestore.io:8080

// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name authorization

// @Tag.name Info
// @Tag.description "Группа запросов состояния сервиса"

// @Tag.name Storage
// @Tag.description "Группа для работы с данными внутри bucket-ов"

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {

	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version:%s\n", buildVersion)
	fmt.Printf("Build date:%s\n", buildDate)
	fmt.Printf("Build commit:%s\n", buildCommit)

	inst.Sugar = zap.NewExample().Sugar()
	defer inst.Sugar.Sync()

	config.ConfigServerInit()
	// маршрутизация запросов обработчику
	r := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// so mach information, switch ON only for debug
	if inst.LogHTTP {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.AllowContentEncoding("deflate", "gzip"))
	r.Use(middleware.Compress(5, "application/json", "html/text", "text/plain", "text/html"))

	// GET requests
	r.Get("/", inhttp.GetAllMetrics)
	r.Get("/value/gauge/{GMname}", inhttp.GetGMvalue)
	r.Get("/value/counter/{CMname}", inhttp.GetCMvalue)
	r.Get("/ping", inhttp.PingStorage)

	// POST requests update, get
	r.Post("/value", inhttp.GetMvalueJSON)
	r.Post("/update", inhttp.UpdateMetricJSON)
	r.Post("/updates", inhttp.UpdateMetricBatch)
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
	// start data's saver persistance
	go inst.StoreMonitor.NewPersistanceStorage()

	// start profiler
	// go http.ListenAndServe(pprofAddr, nil)
	go func() {
		log.Println(http.ListenAndServe(pprofAddr, nil))
	}()

	chanOS := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(chanOS, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-chanOS

	inst.Sugar.Infof("INFO got a signal '%v', start shutting down... wait 5 seconds\n", sig) // put breakpoint here
	inst.StoreMonitor.ClosePersistanceStorage()
	<-time.After(time.Second * 5)

	inst.Sugar.Infof("Shutdown complete")
}

func ListenRutine(r *chi.Mux) {
	if err := http.ListenAndServe(inst.ServerAddress, r); err != nil {
		log.Fatal(err)
	}
}
