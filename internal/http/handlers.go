package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/zvovayar/yandex-go-mustave-devops/internal/crypt"
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

// Ping Storage
func PingStorage(w http.ResponseWriter, r *http.Request) {

	if errp := inst.StoreMonitor.PingSQLserver(r.Context()); errp == nil {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("<h1>Ping database OK</h1>DSN=" + inst.DatabaseDSN))
		if err != nil {
			log.Fatal(err)
		}

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("<h1>Ping database fail</h1>DSN=" + inst.DatabaseDSN + "<br>" + errp.Error()))
		if err != nil {
			log.Fatal(err)
		}
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
	// сохранять значение метрики
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
	// сохранять значение метрики
	//

	sm.SetCMvalue(cmname, inst.Counter(cm))

	log.Printf("Store %v = %d", cmname, cm)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("<h1>Counter metric</h1>" + ss[3] + ss[4]))
	if err != nil {
		log.Fatal(err)
	}
}

// Вернуть все метрики
func GetAllMetrics(w http.ResponseWriter, r *http.Request) {

	htmlText := "<table border=\"1\">"
	for key, element := range inst.Gmetricnames {
		htmlText += fmt.Sprintf("<tr><td>type gauge</td><td> %v</td><td> #%v =</td><td> %f </td></tr>",
			key, element, sm.GetGMvalue(key))
	}

	for key, element := range inst.Cmetricnames {
		htmlText += fmt.Sprintf("<tr><td>type counter</td><td> %v</td><td> #%v =</td><td> %d</td></tr>",
			key, element, sm.GetCMvalue(key))
	}

	htmlText += "</table>"
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(htmlText))
	if err != nil {
		log.Fatal(err)
	}

}

// Вернуть одну метрику Gauge
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

// Вернуть одну метрику Counter
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

// Сохранение одной метрики из JSON
func UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	var v inst.Metrics
	var cmname string
	var gmname string
	var cm inst.Counter
	var gm inst.Gauge

	// decode json form r.Body and init v
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	// log.Println(v)

	//
	// check hash if key exist
	//
	if inst.Key != "" {
		var mc crypt.MetricsCrypt

		mc.M = v
		if !mc.ControlHashMetrics(inst.Key) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("<h1>Bad hash</h1>" + v.MType))

			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Bad hash actual=%v expected=%v", v.Hash, mc.M.Hash)
			return
		}
		// log.Printf("Good hash actual=%v expected=%v", v.Hash, mc.M.Hash)
	}

	//
	// сохранять значение метрики
	//
	if v.MType == "gauge" {

		// log.Printf("*v.Value=%f", *v.Value)
		gmname = v.ID
		gm = inst.Gauge(*v.Value)

		sm.SetGMvalue(gmname, inst.Gauge(gm))

		// log.Printf("Store %v = %f", gmname, gm)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("<h1>Gauge metric</h1>" + gmname))
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if v.MType == "counter" {

		// log.Printf("*v.Delta=%d", *v.Delta)
		cmname = v.ID
		cm = inst.Counter(*v.Delta)

		sm.SetCMvalue(cmname, inst.Counter(cm))

		// log.Printf("Store %v = %d", cmname, cm)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("<h1>Counter metric</h1>" + cmname))

		if err != nil {
			log.Fatal(err)
		}
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte("<h1>Unknown metric type</h1>" + v.MType))

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Unknown metric type %v", v.MType)

}

// Вернуть одну метрику в JSON
func GetMvalueJSON(w http.ResponseWriter, r *http.Request) {

	var v inst.Metrics
	var cmname string
	var gmname string

	// log.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	log.Println(v)

	cmname, gmname = v.ID, v.ID

	if v.MType == "gauge" {

		gm := float64(sm.GetGMvalue(gmname))
		v.Value = &gm
	} else if v.MType == "counter" {
		cm := int64(sm.GetCMvalue(cmname))
		v.Delta = &cm
	} else {
		log.Printf("Error unknown metric type %v", v.MType)
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("<h1>404 metric type not found</h1>"))
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	//
	// compute hash if key exist
	//
	if inst.Key != "" {
		var mc crypt.MetricsCrypt

		mc.M = v
		v.Hash = mc.MakeHashMetrics(inst.Key)
	}

	log.Printf("GetMValueJSON v=%v", v)

	buf, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf)
	if err != nil {
		log.Fatal(err)
	}

}

// Сохранение пакета метрик
func UpdateMetricBatch(w http.ResponseWriter, r *http.Request) {
	//
	// like UpdateMetricJSON, but with slice json
	//
	var mbatch []inst.Metrics

	// decode json form r.Body and init v
	if err := json.NewDecoder(r.Body).Decode(&mbatch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	log.Printf("UpdateMetricBatch mbatch=%v", mbatch)

	// save batch
	if err := sm.SaveBatch(r.Context(), mbatch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	w.WriteHeader(http.StatusOK)

}
