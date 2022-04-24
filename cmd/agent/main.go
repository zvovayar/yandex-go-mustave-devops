package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/zvovayar/yandex-go-mustave-devops/internal"
)

func NewMonitor(duration time.Duration, chanmonitor chan internal.Monitor) {
	var m internal.Monitor
	var rtm runtime.MemStats

	for {
		<-time.After(duration)

		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		// Collect stats
		m.Gmetrics[internal.Gmetricnames["Alloc"]] = internal.Gauge(rtm.Alloc)
		m.Gmetrics[internal.Gmetricnames["BuckHashSys"]] = internal.Gauge(rtm.BuckHashSys)
		m.Gmetrics[internal.Gmetricnames["Frees"]] = internal.Gauge(rtm.Frees)
		m.Gmetrics[internal.Gmetricnames["GCCPUFraction"]] = internal.Gauge(rtm.GCCPUFraction)
		m.Gmetrics[internal.Gmetricnames["GCSys"]] = internal.Gauge(rtm.GCSys)
		m.Gmetrics[internal.Gmetricnames["HeapAlloc"]] = internal.Gauge(rtm.HeapAlloc)
		m.Gmetrics[internal.Gmetricnames["HeapIdle"]] = internal.Gauge(rtm.HeapIdle)
		m.Gmetrics[internal.Gmetricnames["HeapInuse"]] = internal.Gauge(rtm.HeapInuse)
		m.Gmetrics[internal.Gmetricnames["HeapObjects"]] = internal.Gauge(rtm.HeapObjects)
		m.Gmetrics[internal.Gmetricnames["HeapReleased"]] = internal.Gauge(rtm.HeapReleased)
		m.Gmetrics[internal.Gmetricnames["HeapSys"]] = internal.Gauge(rtm.HeapSys)
		m.Gmetrics[internal.Gmetricnames["LastGC"]] = internal.Gauge(rtm.LastGC)
		m.Gmetrics[internal.Gmetricnames["Lookups"]] = internal.Gauge(rtm.Lookups)
		m.Gmetrics[internal.Gmetricnames["MCacheInuse"]] = internal.Gauge(rtm.MCacheInuse)
		m.Gmetrics[internal.Gmetricnames["MCacheSys"]] = internal.Gauge(rtm.MCacheSys)
		m.Gmetrics[internal.Gmetricnames["MSpanInuse"]] = internal.Gauge(rtm.MSpanInuse)
		m.Gmetrics[internal.Gmetricnames["MSpanSys"]] = internal.Gauge(rtm.MSpanSys)
		m.Gmetrics[internal.Gmetricnames["Mallocs"]] = internal.Gauge(rtm.Mallocs)
		m.Gmetrics[internal.Gmetricnames["NextGC"]] = internal.Gauge(rtm.NextGC)
		m.Gmetrics[internal.Gmetricnames["NumForcedGC"]] = internal.Gauge(rtm.NumForcedGC)
		m.Gmetrics[internal.Gmetricnames["NumGC"]] = internal.Gauge(rtm.NumGC)
		m.Gmetrics[internal.Gmetricnames["OtherSys"]] = internal.Gauge(rtm.OtherSys)
		m.Gmetrics[internal.Gmetricnames["PauseTotalNs"]] = internal.Gauge(rtm.PauseTotalNs)
		m.Gmetrics[internal.Gmetricnames["StackInuse"]] = internal.Gauge(rtm.StackInuse)
		m.Gmetrics[internal.Gmetricnames["StackSys"]] = internal.Gauge(rtm.StackSys)
		m.Gmetrics[internal.Gmetricnames["Sys"]] = internal.Gauge(rtm.Sys)
		m.Gmetrics[internal.Gmetricnames["TotalAlloc"]] = internal.Gauge(rtm.TotalAlloc)
		m.Gmetrics[internal.Gmetricnames["RandomValue"]] = internal.Gauge(rand.Float64())

		m.Cmetrics[internal.Cmetricnames["PoolCount"]]++

		// Just encode to json and print
		b, _ := json.Marshal(m)
		log.Printf("NewMonitor - > %v Channel length %v", string(b), len(chanmonitor))

		// Save new collected data to the slice
		chanmonitor <- m

		// m.sendMetrics()

	}
}

func main() {
	chanmonitor := make(chan internal.Monitor, internal.BufferLength)
	chanOS := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(chanOS, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	go NewMonitor(internal.PollInterval, chanmonitor)
	go runSendMetrics(internal.ReportInterval, chanmonitor)

	sig := <-chanOS
	log.Printf("INFO got a signal '%v', start shutting down...\n", sig) // put breakpoint here
	log.Printf("Shutdown complete")
}

func runSendMetrics(duration time.Duration, chanmonitor chan internal.Monitor) {

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
			m.SendMetrics()

		}

	}
}
