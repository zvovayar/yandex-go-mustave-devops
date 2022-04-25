package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Gauge float64
type Counter int64
type Monitor struct {
	// warning: remember about quantity metrics!
	Gmetrics [29]Gauge
	Cmetrics [2]Counter
}

var Gmetricnames = map[string]int{
	"Alloc":         0,
	"BuckHashSys":   1,
	"Frees":         2,
	"GCCPUFraction": 3,
	"GCSys":         4,
	"HeapAlloc":     5,
	"HeapIdle":      6,
	"HeapInuse":     7,
	"HeapObjects":   8,
	"HeapReleased":  9,
	"HeapSys":       10,
	"LastGC":        11,
	"Lookups":       12,
	"MCacheInuse":   13,
	"MCacheSys":     14,
	"MSpanInuse":    15,
	"MSpanSys":      16,
	"Mallocs":       17,
	"NextGC":        18,
	"NumForcedGC":   19,
	"NumGC":         20,
	"OtherSys":      21,
	"PauseTotalNs":  22,
	"StackInuse":    23,
	"StackSys":      24,
	"Sys":           25,
	"TotalAlloc":    26,
	"RandomValue":   27,
	"testSetGet134": 28,
}
var Cmetricnames = map[string]int{
	"PollCount":    0,
	"testSetGet33": 1,
}
var PollInterval = time.Second * 2    // 2
var ReportInterval = time.Second * 10 //10
var BufferLength = 30
var ServerAddress = "127.0.0.1:8080"
var ContentType = "text/plain"

func (m Monitor) SendMetrics() {

	//Just encode to json and print
	b, _ := json.Marshal(m)
	log.Println("SendMetrics -> " + string(b))
	var body = []byte(b)

	// internal.Gauge type send
	for key, element := range Gmetricnames {
		var url = "http://" + ServerAddress + "/update/gauge/" + key + "/" + fmt.Sprint(m.Gmetrics[element])
		log.Println(url)

		request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			// обработаем ошибку
			log.Println(err)
		}
		request.Header.Set("Content-Type", ContentType)

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
	for key, element := range Cmetricnames {
		var url = "http://" + ServerAddress + "/update/counter/" + key + "/" + fmt.Sprint(m.Cmetrics[element])
		log.Println(url)

		request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			// обработаем ошибку
			log.Println(err)
		}
		request.Header.Set("Content-Type", ContentType)

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
