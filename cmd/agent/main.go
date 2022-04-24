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
)

type gauge float64
type counter int64
type monitor struct {
	// warning: remember about quantity metrics!
	Gmetrics [28]gauge
	Cmetrics [1]counter
}

var gmetricnames = map[string]int{
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
}
var cmetricnames = map[string]int{
	"PollCount": 0,
}
var pollInterval = time.Second * 1   // 3
var reportInterval = time.Second * 5 //10
var bufferLength = 20
var serverAddress = "127.0.0.1:8080"
var contentType = "text/plain"

//
// добавить слайс для хранения monitor
//

func NewMonitor(duration time.Duration, chanmonitor chan monitor) {
	var m monitor
	var rtm runtime.MemStats
	//
	// добавить сохранение собранных данных в слайс
	//
	for {
		<-time.After(duration)

		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		// Collect stats
		m.Gmetrics[gmetricnames["Alloc"]] = gauge(rtm.Alloc)
		m.Gmetrics[gmetricnames["BuckHashSys"]] = gauge(rtm.BuckHashSys)
		m.Gmetrics[gmetricnames["Frees"]] = gauge(rtm.Frees)
		m.Gmetrics[gmetricnames["GCCPUFraction"]] = gauge(rtm.GCCPUFraction)
		m.Gmetrics[gmetricnames["GCSys"]] = gauge(rtm.GCSys)
		m.Gmetrics[gmetricnames["HeapAlloc"]] = gauge(rtm.HeapAlloc)
		m.Gmetrics[gmetricnames["HeapIdle"]] = gauge(rtm.HeapIdle)
		m.Gmetrics[gmetricnames["HeapInuse"]] = gauge(rtm.HeapInuse)
		m.Gmetrics[gmetricnames["HeapObjects"]] = gauge(rtm.HeapObjects)
		m.Gmetrics[gmetricnames["HeapReleased"]] = gauge(rtm.HeapReleased)
		m.Gmetrics[gmetricnames["HeapSys"]] = gauge(rtm.HeapSys)
		m.Gmetrics[gmetricnames["LastGC"]] = gauge(rtm.LastGC)
		m.Gmetrics[gmetricnames["Lookups"]] = gauge(rtm.Lookups)
		m.Gmetrics[gmetricnames["MCacheInuse"]] = gauge(rtm.MCacheInuse)
		m.Gmetrics[gmetricnames["MCacheSys"]] = gauge(rtm.MCacheSys)
		m.Gmetrics[gmetricnames["MSpanInuse"]] = gauge(rtm.MSpanInuse)
		m.Gmetrics[gmetricnames["MSpanSys"]] = gauge(rtm.MSpanSys)
		m.Gmetrics[gmetricnames["Mallocs"]] = gauge(rtm.Mallocs)
		m.Gmetrics[gmetricnames["NextGC"]] = gauge(rtm.NextGC)
		m.Gmetrics[gmetricnames["NumForcedGC"]] = gauge(rtm.NumForcedGC)
		m.Gmetrics[gmetricnames["NumGC"]] = gauge(rtm.NumGC)
		m.Gmetrics[gmetricnames["OtherSys"]] = gauge(rtm.OtherSys)
		m.Gmetrics[gmetricnames["PauseTotalNs"]] = gauge(rtm.PauseTotalNs)
		m.Gmetrics[gmetricnames["StackInuse"]] = gauge(rtm.StackInuse)
		m.Gmetrics[gmetricnames["StackSys"]] = gauge(rtm.StackSys)
		m.Gmetrics[gmetricnames["Sys"]] = gauge(rtm.Sys)
		m.Gmetrics[gmetricnames["TotalAlloc"]] = gauge(rtm.TotalAlloc)
		m.Gmetrics[gmetricnames["RandomValue"]] = gauge(rand.Float64())

		m.Cmetrics[cmetricnames["PoolCount"]]++

		// Just encode to json and print
		b, _ := json.Marshal(m)
		fmt.Println("NewMonitor - > " + string(b))

		// Save new collected data to the slice
		chanmonitor <- m

		// m.sendMetrics()

	}
}

func main() {
	chanmonitor := make(chan monitor, bufferLength)
	chanOS := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(chanOS, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	go NewMonitor(pollInterval, chanmonitor)
	go runSendMetrics(reportInterval, chanmonitor)

	sig := <-chanOS
	log.Printf("INFO got a signal '%v', start shutting down...\n", sig) // put breakpoint here
	log.Printf("Shutdown complete")
}

func runSendMetrics(duration time.Duration, chanmonitor chan monitor) {

	for {
		<-time.After(duration)

		c := len(chanmonitor)
		fmt.Printf("runSendMetrics -> %v \n quantity new elements %v\n", time.Now(), c)
		for i := 0; i < c; i++ {

			m, err := <-chanmonitor
			if !err {
				fmt.Println(err)
				break
			}
			m.SendMetrics()

		}

		// for i := 0; i < len(slcMonitor); i++ {

		// 	slcMonitor[i].SendMetrics()
		// }
	}
}

func (m monitor) SendMetrics() {

	//Just encode to json and print
	b, _ := json.Marshal(m)
	fmt.Println("SendMetrics -> " + string(b))
	var body = []byte(b)

	// gauge type send
	for key, element := range gmetricnames {
		var url = "http://" + serverAddress + "/update/gauge/" + key + "/" + fmt.Sprint(m.Gmetrics[element])
		fmt.Println(url)

		request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			// обработаем ошибку
			fmt.Println(err)
		}
		request.Header.Set("Content-Type", contentType)

		// client := &http.Client{}
		// отправляем запрос
		//
		// resp, err := client.Do(request)
		// if err != nil {
		// 	// обработаем ошибку
		// 	fmt.Println(err)
		// 	return
		// }
		// defer resp.Body.Close()
		// fmt.Println(resp)
	}

	// counter type send
	for key, element := range cmetricnames {
		var url = "http://" + serverAddress + "/update/counter/" + key + "/" + fmt.Sprint(m.Cmetrics[element])
		fmt.Println(url)

		request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			// обработаем ошибку
			fmt.Println(err)
		}
		request.Header.Set("Content-Type", contentType)

		// client := &http.Client{}
		// отправляем запрос
		//
		// resp, err := client.Do(request)
		// if err != nil {
		// 	// обработаем ошибку
		// 	fmt.Println(err)
		// 	return
		// }
		// defer resp.Body.Close()
		// fmt.Println(resp)
	}
}
