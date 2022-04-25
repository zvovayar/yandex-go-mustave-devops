package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/zvovayar/yandex-go-mustave-devops/internal"
)

func main() {
	// маршрутизация запросов обработчику

	// root
	http.HandleFunc("/", http.NotFound)
	// update
	http.HandleFunc("/update/", NotImplemented)

	http.HandleFunc("/update/gauge/", UpdateGaugeMetric)
	http.HandleFunc("/update/counter/", UpdateCounterMetric)

	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(":8080", nil)
}

// Не реализовано
func NotImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("<h1>Not implemented</h1> length="))

}

// Сохранение метрики  Gauge
func UpdateGaugeMetric(w http.ResponseWriter, r *http.Request) {

	ss := strings.Split(r.URL.Path, "/")
	log.Printf("%v count=%v", r.URL.Path, len(ss))
	log.Println(ss)

	if len(ss) != 5 {
		// мало или много параметров в URL
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<h1>Gauge metric URL is not valid</h1> length=" + fmt.Sprintf("%d", len(ss))))
		return
	}
	gm, err := strconv.ParseFloat(ss[4], 64)
	if err != nil {
		// значения метрики нет
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<h1>Gauge metric value not found</h1>"))
		return
	} else if _, ok := internal.Gmetricnames[ss[3]]; !ok {
		// не нашли название метрики, были ошибки
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<h1>Gauge metric not found</h1>"))
		return
	}

	gmname := ss[3]
	log.Printf("Gauge metric %v = %f", gmname, gm)

	//
	// TODO: здесь сохранять значение метрики
	//

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Gauge metric</h1>" + ss[3] + ss[4]))
}

// Сохранение метрики Counter
func UpdateCounterMetric(w http.ResponseWriter, r *http.Request) {
	ss := strings.Split(r.URL.Path, "/")
	log.Printf("%v count=%v", r.URL.Path, len(ss))
	log.Println(ss)

	if len(ss) != 5 {
		// мало или много параметров в URL
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<h1>Counter metric URL is not valid</h1> length=" + fmt.Sprintf("%d", len(ss))))
		return
	}

	cm, err := strconv.ParseInt(ss[4], 10, 64)

	if err != nil {
		// значения метрики нет
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<h1>Counter metric value not found</h1>"))
		return
	} else if _, ok := internal.Cmetricnames[ss[3]]; !ok {
		// не нашли название метрики, были ошибки
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<h1>Counter metric not found</h1>"))
		return
	}

	cmname := ss[3]
	log.Printf("Counter metric %v = %d", cmname, cm)

	//
	// TODO: здесь сохранять значение метрики
	//

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Gauge metric</h1>" + ss[3] + ss[4]))
}
