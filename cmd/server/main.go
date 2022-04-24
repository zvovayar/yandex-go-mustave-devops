package main

import (
	"net/http"
)

func main() {
	// маршрутизация запросов обработчику

	// root
	http.HandleFunc("/", http.NotFound)
	// update
	http.HandleFunc("/update/gauge/", UpdateGaugeMetric)
	http.HandleFunc("/update/counter/", UpdateCounterMetric)

	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(":8080", nil)
}

// Сохранение метрики Gauge
func UpdateGaugeMetric(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Gauge metric</h1>"))

}

// Сохранение метрики Counter
func UpdateCounterMetric(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Counter metric</h1>"))

}
