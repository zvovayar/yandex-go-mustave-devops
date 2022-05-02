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

var sm inst.Storage = &inst.StoreMonitor

// Не реализовано
func NotImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	_, err := w.Write([]byte("<h1>Not implemented</h1> length="))
	if err != nil {
		log.Fatal(err)
	}

}

// Сохранение метрики  Gauge
func UpdateGaugeMetric(w http.ResponseWriter, r *http.Request) {

	ss := strings.Split(r.URL.Path, "/")
	gmnamechi := chi.URLParam(r, "GMname")
	gmvaluechi := chi.URLParam(r, "GMvalue")
	log.Printf("%v count=%v %v = %v", r.URL.Path, len(ss), gmnamechi, gmvaluechi)
	log.Println(ss)

	if len(ss) != 5 {
		// мало или много параметров в URL
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("<h1>Gauge metric URL is not valid</h1> length=" + fmt.Sprintf("%d", len(ss))))
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	gm, err := strconv.ParseFloat(ss[4], 64)
	if err != nil {
		// значения метрики нет
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("<h1>Gauge metric value not found</h1>"))
		if err != nil {
			log.Fatal(err)
		}

		return
	} else if _, ok := inst.Gmetricnames[ss[3]]; !ok {
		// не нашли название метрики, были ошибки
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("<h1>Gauge metric not found</h1>"))
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	gmname := ss[3]
	log.Printf("Gauge metric %v = %f", gmname, gm)

	//
	// TODO: здесь сохранять значение метрики
	//
	//storage.StoreMonitor.Gmetrics[Gmetricnames[gmname]] = Gauge(gm)
	// swq := &inst.StoreMonitor //.GetMonitor()
	// s.Gmetrics[inst.Gmetricnames[gmname]] = inst.Gauge(gm)
	sm.SetGMvalue(gmname, inst.Gauge(gm))

	log.Printf("Store %v = %f", gmname, gm)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("<h1>Gauge metric</h1>" + ss[3] + ss[4]))
	if err != nil {
		log.Fatal(err)
	}
}

// Сохранение метрики Counter
func UpdateCounterMetric(w http.ResponseWriter, r *http.Request) {
	ss := strings.Split(r.URL.Path, "/")
	сmnamechi := chi.URLParam(r, "CMname")
	сmvaluechi := chi.URLParam(r, "CMvalue")
	log.Printf("%v count=%v %v = %v", r.URL.Path, len(ss), сmnamechi, сmvaluechi)
	log.Println(ss)

	if len(ss) != 5 {
		// мало или много параметров в URL
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("<h1>Counter metric URL is not valid</h1> length=" + fmt.Sprintf("%d", len(ss))))
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	cm, err := strconv.ParseInt(ss[4], 10, 64)

	if err != nil {
		// значения метрики нет
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("<h1>Counter metric value not found</h1>"))
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if _, ok := inst.Cmetricnames[ss[3]]; !ok {
		// не нашли название метрики, были ошибки
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("<h1>Counter metric not found</h1>"))
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	cmname := ss[3]
	log.Printf("Counter metric %v = %d", cmname, cm)

	//
	// TODO: здесь сохранять значение метрики
	//
	//storage.StoreMonitor.Cmetrics[Cmetricnames[cmname]] += Counter(cm)
	// s := &inst.StoreMonitor //.GetMonitor()
	// s.Cmetrics[inst.Cmetricnames[cmname]] += inst.Counter(cm)
	sm.SetCMvalue(cmname, inst.Counter(cm))

	log.Printf("Store %v = %d", cmname, cm)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("<h1>Gauge metric</h1>" + ss[3] + ss[4]))
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllMetrics(w http.ResponseWriter, r *http.Request) {

	htmlText := ""
	for key, element := range inst.Gmetricnames {
		htmlText += fmt.Sprintf("type gauge %v #%v = %f \n", key, element, sm.GetGMvalue(key)) //inst.StoreMonitor.Gmetrics[inst.Gmetricnames[key]])
	}

	for key, element := range inst.Cmetricnames {
		htmlText += fmt.Sprintf("type counter %v #%v = %d \n", key, element, sm.GetCMvalue(key)) //inst.StoreMonitor.Cmetrics[inst.Cmetricnames[key]])
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(htmlText))
	if err != nil {
		log.Fatal(err)
	}

}

func GetGMvalue(w http.ResponseWriter, r *http.Request) {

	if _, ok := inst.Gmetricnames[chi.URLParam(r, "GMname")]; !ok {
		// не нашли название метрики, были ошибки
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("<h1>404 Gauge metric not found</h1>"))
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// htmlText := fmt.Sprint(inst.StoreMonitor.Gmetrics[inst.Gmetricnames[chi.URLParam(r, "GMname")]])
	htmlText := fmt.Sprint(sm.GetGMvalue(chi.URLParam(r, "GMname")))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(htmlText))
	if err != nil {
		log.Fatal(err)
	}
}

func GetCMvalue(w http.ResponseWriter, r *http.Request) {

	if _, ok := inst.Cmetricnames[chi.URLParam(r, "CMname")]; !ok {
		// не нашли название метрики, были ошибки
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("<h1>404 Counter metric not found</h1>"))
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// htmlText := fmt.Sprint(inst.StoreMonitor.Cmetrics[inst.Cmetricnames[chi.URLParam(r, "CMname")]])
	htmlText := fmt.Sprint(sm.GetCMvalue(chi.URLParam(r, "CMname")))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(htmlText))
	if err != nil {
		log.Fatal(err)
	}
}
