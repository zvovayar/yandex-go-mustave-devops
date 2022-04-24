package main

import (
	"net/http"
)

func main() {
	// маршрутизация запросов обработчику

	// root
	http.HandleFunc("/", http.NotFound)
	// update
	http.HandleFunc("/update/", UpdateMetric)

	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(":8080", nil)
}

// Сохранение метрики
func UpdateMetric(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Hello, World</h1>"))

}
