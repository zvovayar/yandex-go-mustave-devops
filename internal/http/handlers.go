package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

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
	} else if _, ok := inst.Gmetricnames[ss[3]]; !ok {
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
	//storage.StoreMonitor.Gmetrics[Gmetricnames[gmname]] = Gauge(gm)
	s := &inst.StoreMonitor //.GetMonitor()
	s.Gmetrics[inst.Gmetricnames[gmname]] = inst.Gauge(gm)
	log.Printf("Store %v = %f", gmname, gm)

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
	} else if _, ok := inst.Cmetricnames[ss[3]]; !ok {
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
	//storage.StoreMonitor.Cmetrics[Cmetricnames[cmname]] += Counter(cm)
	s := &inst.StoreMonitor //.GetMonitor()
	s.Cmetrics[inst.Cmetricnames[cmname]] += inst.Counter(cm)
	log.Printf("Store %v = %d", cmname, cm)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Gauge metric</h1>" + ss[3] + ss[4]))
}

func GetAllMetrics(w http.ResponseWriter, r *http.Request) {

	htmlText := ""
	for key, element := range inst.Gmetricnames {
		htmlText += fmt.Sprintf("type gauge %v #%v = %f \n", key, element, inst.StoreMonitor.Gmetrics[inst.Gmetricnames[key]])
	}

	for key, element := range inst.Cmetricnames {
		htmlText += fmt.Sprintf("type counter %v #%v = %d \n", key, element, inst.StoreMonitor.Cmetrics[inst.Cmetricnames[key]])
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlText))

}

func GetGMvalue(w http.ResponseWriter, r *http.Request) {

	if _, ok := inst.Gmetricnames[chi.URLParam(r, "GMname")]; !ok {
		// не нашли название метрики, были ошибки
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<h1>404 Gauge metric not found</h1>"))
		return
	}

	htmlText := fmt.Sprint(inst.StoreMonitor.Gmetrics[inst.Gmetricnames[chi.URLParam(r, "GMname")]])
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlText))
}

func GetCMvalue(w http.ResponseWriter, r *http.Request) {

	if _, ok := inst.Cmetricnames[chi.URLParam(r, "CMname")]; !ok {
		// не нашли название метрики, были ошибки
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<h1>404 Counter metric not found</h1>"))
		return
	}

	htmlText := fmt.Sprint(inst.StoreMonitor.Cmetrics[inst.Cmetricnames[chi.URLParam(r, "CMname")]])
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlText))
}
