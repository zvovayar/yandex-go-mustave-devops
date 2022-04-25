package main

import (
	"net/http"

	"github.com/zvovayar/yandex-go-mustave-devops/cmd/server/handlers"
)

func main() {
	// маршрутизация запросов обработчику

	// root
	http.HandleFunc("/", http.NotFound)
	// update
	http.HandleFunc("/update/", handlers.NotImplemented)

	http.HandleFunc("/update/gauge/", handlers.UpdateGaugeMetric)
	http.HandleFunc("/update/counter/", handlers.UpdateCounterMetric)

	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(":8080", nil)
}
