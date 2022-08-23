// Package http include handlers and senders HTTP functions
package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/zvovayar/yandex-go-mustave-devops/internal/crypt"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

var sm inst.Storage = &inst.StoreMonitor

// NotImplemented Not implemented stub
func NotImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	_, err := w.Write([]byte("<h1>Not implemented</h1> length="))
	if err != nil {
		log.Fatal(err)
	}

}

// PingStorage ping database Storage
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

// UpdateGaugeMetric save Gauge metric
func UpdateGaugeMetric(w http.ResponseWriter, r *http.Request) {

	ss := strings.Split(r.URL.Path, "/")
	gmnamechi := chi.URLParam(r, "GMname")
	gmvaluechi := chi.URLParam(r, "GMvalue")
	inst.Sugar.Infof("%v count=%v %v = %v", r.URL.Path, len(ss), gmnamechi, gmvaluechi)

	if len(ss) != 5 {
		// мало или много параметров в URL
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("<h1>Gauge metric URL is not valid</h1> length=" + fmt.Sprintf("%d", len(ss))))
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	// gm, err := strconv.ParseFloat(ss[4], 64)
	gm, err := strconv.ParseFloat(gmvaluechi, 64)
	if err != nil {
		// значения метрики нет
		inst.Sugar.Infow(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("<h1>Gauge metric value not found</h1>"))
		if err != nil {
			log.Fatal(err)
		}

		return
	}
	inst.Sugar.Infof("Gauge metric %v = %f", gmnamechi, gm)

	//
	// сохранять значение метрики
	//
	sm.SetGMvalue(gmnamechi, inst.Gauge(gm))

	inst.Sugar.Infof("Store %v = %f", gmnamechi, gm)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("<h1>Gauge metric</h1>" + gmnamechi + gmvaluechi))
	if err != nil {
		log.Fatal(err)
	}
}

// UpdateCounterMetric save Counter metric
func UpdateCounterMetric(w http.ResponseWriter, r *http.Request) {
	ss := strings.Split(r.URL.Path, "/")
	cmnamechi := chi.URLParam(r, "CMname")
	cmvaluechi := chi.URLParam(r, "CMvalue")
	inst.Sugar.Infof("%v count=%v %v = %v", r.URL.Path, len(ss), cmnamechi, cmvaluechi)

	if len(ss) != 5 {
		// мало или много параметров в URL
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("<h1>Counter metric URL is not valid</h1> length=" + fmt.Sprintf("%d", len(ss))))
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	cm, err := strconv.ParseInt(cmvaluechi, 10, 64)

	if err != nil {
		// значения метрики нет
		inst.Sugar.Infow(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("<h1>Counter metric value not found</h1>"))
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	inst.Sugar.Infof("Counter metric %v = %d", cmnamechi, cm)

	//
	// сохранять значение метрики
	//

	sm.SetCMvalue(cmnamechi, inst.Counter(cm))

	inst.Sugar.Infof("Store %v = %d", cmnamechi, cm)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("<h1>Counter metric</h1>" + cmnamechi + cmvaluechi))
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllMetrics(w http.ResponseWriter, r *http.Request) {

	// sort Gnames
	keysG := make([]string, 0, len(inst.Gmetricnames))
	for k := range inst.Gmetricnames {
		keysG = append(keysG, k)
	}
	sort.Strings(keysG)

	// sort Cnames
	keysC := make([]string, 0, len(inst.Cmetricnames))
	for k := range inst.Cmetricnames {
		keysC = append(keysC, k)
	}
	sort.Strings(keysG)
	sort.Strings(keysC)

	htmlText := "<table border=\"1\">"
	for _, key := range keysG {
		htmlText += fmt.Sprintf("<tr><td>type gauge</td><td> %v</td><td> #%v =</td><td> %f </td></tr>",
			key, inst.Gmetricnames[key], sm.GetGMvalue(key))
	}

	for _, key := range keysC {
		htmlText += fmt.Sprintf("<tr><td>type counter</td><td> %v</td><td> #%v =</td><td> %d</td></tr>",
			key, inst.Cmetricnames[key], sm.GetCMvalue(key))
	}

	htmlText += "</table>"
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(htmlText))
	if err != nil {
		log.Fatal(err)
	}

}

// GetGMvalue get one metric gauge
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

	htmlText := fmt.Sprint(sm.GetGMvalue(chi.URLParam(r, "GMname")))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(htmlText))
	if err != nil {
		log.Fatal(err)
	}
}

// GetCMvalue get one metric Counter
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

	htmlText := fmt.Sprint(sm.GetCMvalue(chi.URLParam(r, "CMname")))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(htmlText))
	if err != nil {
		log.Fatal(err)
	}
}

// UpdateMetricJSON save one metric from JSON
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
	// inst.Sugar.Infow(v)

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
			inst.Sugar.Infof("Bad hash actual=%v expected=%v", v.Hash, mc.M.Hash)
			return
		}
		// inst.Sugar.Infof("Good hash actual=%v expected=%v", v.Hash, mc.M.Hash)
	}

	//
	// сохранять значение метрики
	//
	if v.MType == "gauge" {

		// inst.Sugar.Infof("*v.Value=%f", *v.Value)
		gmname = v.ID
		gm = inst.Gauge(*v.Value)

		sm.SetGMvalue(gmname, inst.Gauge(gm))

		// inst.Sugar.Infof("Store %v = %f", gmname, gm)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("<h1>Gauge metric</h1>" + gmname))
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if v.MType == "counter" {

		// inst.Sugar.Infof("*v.Delta=%d", *v.Delta)
		cmname = v.ID
		cm = inst.Counter(*v.Delta)

		sm.SetCMvalue(cmname, inst.Counter(cm))

		// inst.Sugar.Infof("Store %v = %d", cmname, cm)

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
	inst.Sugar.Infof("Unknown metric type %v", v.MType)

}

// GetMvalueJSON get one metric in JSON
func GetMvalueJSON(w http.ResponseWriter, r *http.Request) {

	var v inst.Metrics
	var cmname string
	var gmname string

	// inst.Sugar.Infow(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err)
		return
	}
	inst.Sugar.Infof("v=%v", v)

	cmname, gmname = v.ID, v.ID

	_, okg := inst.Gmetricnames[gmname]
	_, okc := inst.Cmetricnames[cmname]

	if v.MType == "gauge" && okg {
		gm := float64(sm.GetGMvalue(gmname))
		v.Value = &gm
	} else if v.MType == "counter" && okc {
		cm := int64(sm.GetCMvalue(cmname))
		v.Delta = &cm
	} else {
		inst.Sugar.Infof("Error unknown metric type or name %v, %v", v.MType, v.ID)
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("<h1>404 metric type or name not found</h1>"))
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

	inst.Sugar.Infof("GetMValueJSON v=%v", v)

	buf, err := json.Marshal(v)
	if err != nil {
		inst.Sugar.Infof("GetMValueJSON marshal error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf)
	if err != nil {
		inst.Sugar.Infof("GetMValueJSON Write error: %v", err)
	}
	inst.Sugar.Infof("GetMValueJSON string(buf)=%v", string(buf))

}

// UpdateMetricBatch save batch of metrics
// like UpdateMetricJSON, but with slice in json
func UpdateMetricBatch(w http.ResponseWriter, r *http.Request) {
	var mbatch []inst.Metrics

	// decode json form r.Body and init v
	if err := json.NewDecoder(r.Body).Decode(&mbatch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
		return
	}

	// save batch
	if err := sm.SaveBatch(r.Context(), mbatch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		inst.Sugar.Infow(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)

}
