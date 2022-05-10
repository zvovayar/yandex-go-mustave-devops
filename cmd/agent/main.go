package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

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

		// Just encode to json and print
		b, _ := json.Marshal(m)
		log.Printf("NewMonitor - > %v Channel length %v", string(b), len(chanmonitor))

		// Save new collected data to the slice
		chanmonitor <- m

		// m.sendMetrics()

	}
}

func SendMetrics(m inst.Monitor) {

	//Just encode to json and print
	b, _ := json.Marshal(m)
	log.Println("SendMetrics -> " + string(b))
	// var body = []byte(b)

	// internal.Gauge type send
	for key, element := range inst.Gmetricnames {

		var v inst.Metrics

		v.ID = key
		v.MType = "gauge"
		v.Value = (*float64)(&m.Gmetrics[element])
		body, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("v=%v", v)
		log.Printf("body=%v", string(body))

		// var url = fmt.Sprintf("http://%v/update/gauge/%v/%f",
		// 	inst.ServerAddress, key, m.Gmetrics[element])
		var url = fmt.Sprintf("http://%v/update/",
			inst.ServerAddress)
		log.Println(url)

		request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			// обработаем ошибку
			log.Println(err)
		}
		request.Header.Set("Content-Type", inst.ContentType)

		client := &http.Client{}
		// отправляем запрос

		resp, err := client.Do(request)
		if err != nil {
			// обработаем ошибку
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		log.Println(resp)
	}

	// counter type send
	for key, element := range inst.Cmetricnames {
		var v inst.Metrics

		v.ID = key
		v.MType = "counter"
		v.Delta = (*int64)(&m.Cmetrics[element])
		body, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("v=%v", v)
		log.Printf("body=%v", string(body))

		// var url = fmt.Sprintf("http://%v/update/counter/%v/%d",
		// 	inst.ServerAddress, key, m.Cmetrics[element])
		var url = fmt.Sprintf("http://%v/update/",
			inst.ServerAddress)
		log.Println(url)

		request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			// обработаем ошибку
			log.Println(err)
		}
		request.Header.Set("Content-Type", inst.ContentType)

		client := &http.Client{}
		// отправляем запрос

		resp, err := client.Do(request)
		if err != nil {
			// обработаем ошибку
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		log.Println(resp)
	}
}
func runSendMetrics(duration time.Duration, chanmonitor chan inst.Monitor) {

	for {
		<-time.After(duration)

		c := len(chanmonitor)
		log.Printf("runSendMetrics -> quantity new elements %v\n", c)
		for i := 0; i < c; i++ {

			m, err := <-chanmonitor
			if !err {
				fmt.Println(err)
				break
			}
			SendMetrics(m)
			// inhttp.SendMetrics()

		}

	}
}

func main() {
	log.Println("Agent started")
	chanmonitor := make(chan inst.Monitor, inst.BufferLength)
	chanOS := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(chanOS, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	go NewMonitor(inst.PollInterval, chanmonitor)
	go runSendMetrics(inst.ReportInterval, chanmonitor)

	sig := <-chanOS
	log.Printf("INFO got a signal '%v', start shutting down...\n", sig) // put breakpoint here
	log.Printf("Shutdown complete")
}
