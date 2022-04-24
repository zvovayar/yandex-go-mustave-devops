package main

import (
	"net/http"
)

func main() {
	// маршрутизация запросов обработчику

	// root
	http.HandleFunc("/", HelloWorld)
	// update
	http.HandleFunc("/update/", UpdateMetric)

	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(":8080", nil)
}

// HelloWorld — обработчик запроса.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello, World</h1>"))

}

// HelloWorld — обработчик запроса.
func UpdateMetric(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello, World</h1>"))

}
