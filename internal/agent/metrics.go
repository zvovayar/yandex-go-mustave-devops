// Package agent contain core functions agent metrics
// thread safety functions
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

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/zvovayar/yandex-go-mustave-devops/internal/crypt"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

// NewMonitor begin collect metrics infinitly and send they to the channel
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

		//
		// TODO: make MonitorAgent struct and send it
		//
		chanmonitor <- m

	}
}

// NewMonitorGopsutil begin collect metrics infinitly and send they to the channel
func NewMonitorGopsutil(duration time.Duration, chanmonitor chan inst.Monitor) {
	var m inst.Monitor

	cpuCounts, _ := cpu.Counts(false)

	m.Cmetrics = make([]inst.Counter, len(inst.Cmetricnames))

	startLenGmetricnames := len(inst.Gmetricnames)
	// fill CPUs metrics names
	for i := 0; i < cpuCounts; i++ {
		inst.Gmetricnames["CPUutilization"+fmt.Sprint(i+1)] = startLenGmetricnames + i
	}
	m.Gmetrics = make([]inst.Gauge, len(inst.Gmetricnames)+cpuCounts)

	inst.Sugar.Debugf("NewMonitorGopsutil: Gmetricnames=%v", inst.Gmetricnames)

	for {
		<-time.After(duration)

		vmem, _ := mem.VirtualMemory()

		// Collect stats
		m.Gmetrics[inst.Gmetricnames["TotalMemory"]] = inst.Gauge(vmem.Total)
		m.Gmetrics[inst.Gmetricnames["FreeMemory"]] = inst.Gauge(vmem.Free)

		// Collect cpu utilization stats for all CPUs
		cpuutil, err := cpu.Percent(duration, true)
		if err != nil {
			inst.Sugar.Error(err.Error())
		}

		for i := 0; i < cpuCounts; i++ {

			onecpuutil := cpuutil[i]
			m.Gmetrics[inst.Gmetricnames["CPUutilization"+fmt.Sprint(i+1)]] = inst.Gauge(onecpuutil)
		}
		// Send new collected data to the channel

		//
		// TODO for plan B: make MonitorAgent struct and send it
		//
		inst.Sugar.Debugf("NewMonitorGopsutil: m=%v send it to chanmonitor", m)

		chanmonitor <- m

	}
}

// SendMetrics send collected metrics to the web API
func SendMetrics(m inst.Monitor) {

	var mc crypt.MetricsCrypt

	// internal.Gauge type send
	for key, element := range inst.Gmetricnames {

		var v inst.Metrics

		v.ID = key
		v.MType = "gauge"
		if len(m.Gmetrics) < element+1 {
			break
		}
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

		inst.Sugar.Debugf("v=%v", v)
		inst.Sugar.Debugf("body=%v", string(body))

		var url = fmt.Sprintf("http://%v/update/",
			inst.ServerAddress)
		inst.Sugar.Debugf(url)

		request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			// обработаем ошибку
			inst.Sugar.Error(err.Error())
		}
		request.Header.Set("Content-Type", inst.ContentType)

		client := &http.Client{}

		// отправляем запрос
		resp, err := client.Do(request)
		if err != nil {
			// обработаем ошибку
			inst.Sugar.Error(err.Error())
			return
		}
		defer resp.Body.Close()
		inst.Sugar.Debug(resp.Status)
	}

	// counter type send
	for key, element := range inst.Cmetricnames {
		var v inst.Metrics

		if len(m.Cmetrics) < element+1 {
			break
		}

		v.ID = key
		v.MType = "counter"
		v.Delta = (*int64)(&m.Cmetrics[element])

		if inst.Key != "" {
			mc.M = v
			v.Hash = mc.MakeHashMetrics(inst.Key)
		}

		body, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}

		inst.Sugar.Debugf("v=%v", v)
		inst.Sugar.Debugf("body=%v", string(body))

		var url = fmt.Sprintf("http://%v/update/",
			inst.ServerAddress)
		inst.Sugar.Debug(url)

		request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			// обработаем ошибку
			inst.Sugar.Error(err.Error())
		}
		request.Header.Set("Content-Type", inst.ContentType)

		client := &http.Client{}
		// отправляем запрос

		resp, err := client.Do(request)
		if err != nil {
			// обработаем ошибку
			inst.Sugar.Error(err.Error())
			return
		}
		defer resp.Body.Close()
		inst.Sugar.Debug(resp.Status)
	}
}

// RunSendMetrics begin waiting metrics from channel and send they to the web APIs
func RunSendMetrics(duration time.Duration, chanmonitor1 chan inst.Monitor, chanmonitor2 chan inst.Monitor) {

	inst.Sugar.Infow("Agent started gorutine for send metrics")

	var M, m inst.Monitor
	var err bool

	M.Cmetrics = make([]inst.Counter, len(inst.Cmetricnames))
	M.Gmetrics = make([]inst.Gauge, len(inst.Gmetricnames))

	for {
		<-time.After(duration)

		// read first chan
		c1 := len(chanmonitor1)
		c2 := len(chanmonitor2)
		var C int
		if c1 > c2 {
			C = c1
		} else {
			C = c2
		}

		mslice := make([]inst.Monitor, C)

		inst.Sugar.Infof("runSendMetrics -> chanmonitor1 quantity new elements %v", c1)
		inst.Sugar.Infof("runSendMetrics -> chanmonitor2 quantity new elements %v", c2)

		for i := 0; i < C; i++ {

			// read next Monitor from channels
			if i < c1 {
				inst.Sugar.Debugf("RunSendMetrics: read i=%d element from chanmonitor1 ", i)

				m, err = <-chanmonitor1
				if !err {
					inst.Sugar.Infow("RunSendMetrics: chanmonitor1 ended... why?")
					break
				}
				// add only first 28 metrics
				// allocate memory if m > M
				if len(M.Cmetrics) < len(m.Cmetrics) {
					x := m.Cmetrics[len(M.Cmetrics)-1 : len(m.Cmetrics)]
					M.Cmetrics = append(M.Cmetrics, x...)
				}
				if len(M.Gmetrics) < len(m.Gmetrics) {
					x := m.Gmetrics[len(M.Gmetrics)-1 : len(m.Gmetrics)]
					M.Gmetrics = append(M.Gmetrics, x...)
				}
				copy(M.Cmetrics, m.Cmetrics)
				CopyPartSliceG(M.Gmetrics, m.Gmetrics, 0, 28)
			}
			if i < c2 {
				inst.Sugar.Debugf("RunSendMetrics: read i=%d element from chanmonitor2 ", i)
				m, err = <-chanmonitor2
				if !err {
					inst.Sugar.Infow("RunSendMetrics: chanmonitor2 ended... why?")
					break
				}
				// allocate memory if m > M
				if len(M.Cmetrics) < len(m.Cmetrics) {
					x := m.Cmetrics[len(M.Cmetrics)-1 : len(m.Cmetrics)]
					M.Cmetrics = append(M.Cmetrics, x...)
				}
				if len(M.Gmetrics) < len(m.Gmetrics) {
					x := m.Gmetrics[len(M.Gmetrics)-1 : len(m.Gmetrics)]
					M.Gmetrics = append(M.Gmetrics, x...)
				}
				// copy(M.Cmetrics, m.Cmetrics)
				CopyPartSliceG(M.Gmetrics, m.Gmetrics, 28, len(M.Gmetrics))
			}
			inst.Sugar.Debugf("RunSendMetrics Gmetricnames=%v", inst.Gmetricnames)
			inst.Sugar.Debugf("RunSendMetrics Cmetricnames=%v", inst.Cmetricnames)
			inst.Sugar.Debugf("RunSendMetrics M=%v", M)
			if inst.BatchSend {
				// add Metrics to the slice of Monitors
				mslice[i] = M
			} else {
				SendMetrics(M)
			}

		}
		if inst.BatchSend && C > 0 {
			SendBatchMetrics(mslice)
		}

	}
}

func CopyPartSliceG(sdst []inst.Gauge, ssource []inst.Gauge, begin int, end int) {

	inst.Sugar.Debugf("CopyPartSliceG sdst=%v", sdst)
	inst.Sugar.Debugf("CopyPartSliceG ssource=%v", ssource)

	if begin < 0 || begin >= len(sdst) {
		return
	}
	if end < begin {
		return
	}
	if end > len(sdst) {
		end = len(sdst)
	}
	for i := begin; i < end && i < len(ssource); i++ {
		sdst[i] = ssource[i]
	}
}

// SendBatchMetrics create slices []Metrics and send they POST /updates/
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

			if len(monitorb[i].Gmetrics) < element+1 {
				break
			}
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

			if len(monitorb[i].Cmetrics) < element+1 {
				break
			}
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
	inst.Sugar.Debugf("SendBatchMetrics -> count=%d metricsb=%v", c, metricsb)

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
