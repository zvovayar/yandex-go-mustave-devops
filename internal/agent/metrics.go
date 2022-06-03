package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/zvovayar/yandex-go-mustave-devops/internal/crypt"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

// begin collect metrics infinitly and send they to the channel
func NewMonitor(duration time.Duration, chanmonitor chan inst.Monitor) {
	var m inst.Monitor
	var rtm runtime.MemStats

	m.Cmetrics = make([]inst.Counter, len(inst.Cmetricnames))
	m.Gmetrics = make([]inst.Gauge, len(inst.Gmetricnames))

	for {
		<-time.After(duration)

		// Read full memory stats
		runtime.ReadMemStats(&rtm)

		// Collect stats
		m.Gmetrics[inst.Gmetricnames["Alloc"]] = inst.Gauge(rtm.Alloc)
		m.Gmetrics[inst.Gmetricnames["BuckHashSys"]] = inst.Gauge(rtm.BuckHashSys)
		m.Gmetrics[inst.Gmetricnames["Frees"]] = inst.Gauge(rtm.Frees)
		m.Gmetrics[inst.Gmetricnames["GCCPUFraction"]] = inst.Gauge(rtm.GCCPUFraction)
		m.Gmetrics[inst.Gmetricnames["GCSys"]] = inst.Gauge(rtm.GCSys)
		m.Gmetrics[inst.Gmetricnames["HeapAlloc"]] = inst.Gauge(rtm.HeapAlloc)
		m.Gmetrics[inst.Gmetricnames["HeapIdle"]] = inst.Gauge(rtm.HeapIdle)
		m.Gmetrics[inst.Gmetricnames["HeapInuse"]] = inst.Gauge(rtm.HeapInuse)
		m.Gmetrics[inst.Gmetricnames["HeapObjects"]] = inst.Gauge(rtm.HeapObjects)
		m.Gmetrics[inst.Gmetricnames["HeapReleased"]] = inst.Gauge(rtm.HeapReleased)
		m.Gmetrics[inst.Gmetricnames["HeapSys"]] = inst.Gauge(rtm.HeapSys)
		m.Gmetrics[inst.Gmetricnames["LastGC"]] = inst.Gauge(rtm.LastGC)
		m.Gmetrics[inst.Gmetricnames["Lookups"]] = inst.Gauge(rtm.Lookups)
		m.Gmetrics[inst.Gmetricnames["MCacheInuse"]] = inst.Gauge(rtm.MCacheInuse)
		m.Gmetrics[inst.Gmetricnames["MCacheSys"]] = inst.Gauge(rtm.MCacheSys)
		m.Gmetrics[inst.Gmetricnames["MSpanInuse"]] = inst.Gauge(rtm.MSpanInuse)
		m.Gmetrics[inst.Gmetricnames["MSpanSys"]] = inst.Gauge(rtm.MSpanSys)
		m.Gmetrics[inst.Gmetricnames["Mallocs"]] = inst.Gauge(rtm.Mallocs)
		m.Gmetrics[inst.Gmetricnames["NextGC"]] = inst.Gauge(rtm.NextGC)
		m.Gmetrics[inst.Gmetricnames["NumForcedGC"]] = inst.Gauge(rtm.NumForcedGC)
		m.Gmetrics[inst.Gmetricnames["NumGC"]] = inst.Gauge(rtm.NumGC)
		m.Gmetrics[inst.Gmetricnames["OtherSys"]] = inst.Gauge(rtm.OtherSys)
		m.Gmetrics[inst.Gmetricnames["PauseTotalNs"]] = inst.Gauge(rtm.PauseTotalNs)
		m.Gmetrics[inst.Gmetricnames["StackInuse"]] = inst.Gauge(rtm.StackInuse)
		m.Gmetrics[inst.Gmetricnames["StackSys"]] = inst.Gauge(rtm.StackSys)
		m.Gmetrics[inst.Gmetricnames["Sys"]] = inst.Gauge(rtm.Sys)
		m.Gmetrics[inst.Gmetricnames["TotalAlloc"]] = inst.Gauge(rtm.TotalAlloc)
		m.Gmetrics[inst.Gmetricnames["RandomValue"]] = inst.Gauge(rand.Float64())

		m.Cmetrics[inst.Cmetricnames["PoolCount"]]++

		// Send new collected data to the channel
		chanmonitor <- m

	}
}

// send collected metrics to the web API
func SendMetrics(m inst.Monitor) {

	var mc crypt.MetricsCrypt

	// internal.Gauge type send
	for key, element := range inst.Gmetricnames {

		var v inst.Metrics

		v.ID = key
		v.MType = "gauge"
		v.Value = (*float64)(&m.Gmetrics[element])

		if inst.Key != "" {
			mc.M = v
			v.Hash = mc.MakeHashMetrics(inst.Key)
		}

		// inst.Sugar.Infof("agent.SendMetrics v.Hash=%v", v.Hash)

		body, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}

		inst.Sugar.Infof("v=%v", v)
		inst.Sugar.Infof("body=%v", string(body))

		var url = fmt.Sprintf("http://%v/update/",
			inst.ServerAddress)
		inst.Sugar.Infow(url)

		request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			// обработаем ошибку
			inst.Sugar.Infow(err.Error())
		}
		request.Header.Set("Content-Type", inst.ContentType)

		client := &http.Client{}

		// отправляем запрос
		resp, err := client.Do(request)
		if err != nil {
			// обработаем ошибку
			inst.Sugar.Infow(err.Error())
			return
		}
		defer resp.Body.Close()
		inst.Sugar.Infow(resp.Status)
	}

	// counter type send
	for key, element := range inst.Cmetricnames {
		var v inst.Metrics

		v.ID = key
		v.MType = "counter"
		v.Delta = (*int64)(&m.Cmetrics[element])

		if inst.Key != "" {
			mc.M = v
			v.Hash = mc.MakeHashMetrics(inst.Key)
		}

		// inst.Sugar.Infof("agent.SendMetrics v.Hash=%v", v.Hash)

		body, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}

		inst.Sugar.Infof("v=%v", v)
		inst.Sugar.Infof("body=%v", string(body))

		var url = fmt.Sprintf("http://%v/update/",
			inst.ServerAddress)
		inst.Sugar.Infow(url)

		request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			// обработаем ошибку
			inst.Sugar.Infow(err.Error())
		}
		request.Header.Set("Content-Type", inst.ContentType)

		client := &http.Client{}
		// отправляем запрос

		resp, err := client.Do(request)
		if err != nil {
			// обработаем ошибку
			inst.Sugar.Infow(err.Error())
			return
		}
		defer resp.Body.Close()
		inst.Sugar.Infow(resp.Status)
	}
}

// begin waiting metrics from channel and send they to the web APIs
func RunSendMetrics(duration time.Duration, chanmonitor chan inst.Monitor) {

	inst.Sugar.Infow("Agent started gorutine for send metrics")
	for {
		<-time.After(duration)

		c := len(chanmonitor)
		mslice := make([]inst.Monitor, c)

		inst.Sugar.Infof("runSendMetrics -> quantity new elements %v\n", c)
		for i := 0; i < c; i++ {

			m, err := <-chanmonitor
			if !err {
				inst.Sugar.Infow("chan ended... why?")
				break
			}
			if inst.BatchSend {
				//
				// add Metrics to the slice of Monitors
				//
				mslice[i] = m
			} else {
				SendMetrics(m)
			}

		}
		if inst.BatchSend && c > 0 {
			SendBatchMetrics(mslice)
		}
	}
}

//
// create slices []Metrics and send they POST /uodates/
//
func SendBatchMetrics(monitorb []inst.Monitor) {

	var mc crypt.MetricsCrypt
	var metricsb []inst.Metrics
	c := len(monitorb)
	if c == 0 {
		return
	}
	for i := 0; i < c; i++ {
		// Gauge type add to []Metrics
		for key, element := range inst.Gmetricnames {
			var v inst.Metrics

			v.ID = key
			v.MType = "gauge"
			v.Value = (*float64)(&(monitorb[i].Gmetrics[element]))

			if inst.Key != "" {
				mc.M = v
				v.Hash = mc.MakeHashMetrics(inst.Key)
			}
			metricsb = append(metricsb, v)
		}

		// counter type send
		for key, element := range inst.Cmetricnames {
			var v inst.Metrics

			v.ID = key
			v.MType = "counter"
			v.Delta = (*int64)(&(monitorb[i].Cmetrics[element]))

			if inst.Key != "" {
				mc.M = v
				v.Hash = mc.MakeHashMetrics(inst.Key)
			}
			metricsb = append(metricsb, v)
		}
	}
	body, err := json.Marshal(metricsb)
	if err != nil {
		log.Fatal(err)
	}
	inst.Sugar.Infof("SendBatchMetrics -> count=%d metricsb=%v", c, metricsb)
	// inst.Sugar.Infof("SendBatchMetrics -> count=%d metricsb=%v", c, string(body))

	//
	// send json via POST
	//
	var url = fmt.Sprintf("http://%v/updates/",
		inst.ServerAddress)
	inst.Sugar.Infow(url)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		// обработаем ошибку
		inst.Sugar.Infow(err.Error())
	}
	request.Header.Set("Content-Type", inst.ContentType)

	client := &http.Client{}
	// отправляем запрос

	resp, err := client.Do(request)
	if err != nil {
		// обработаем ошибку
		inst.Sugar.Infow(err.Error())
		return
	}
	defer resp.Body.Close()
	inst.Sugar.Infow(resp.Status)
}
