// Package agent contain core functions agent metrics
// thread safety functions
package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/zvovayar/yandex-go-mustave-devops/internal/crypt"
	"github.com/zvovayar/yandex-go-mustave-devops/internal/proto"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var AgentEncrypter crypt.Encrypter

var GrpcUsersClient *proto.UsersClient

// NewMonitor begin collect metrics infinitly and send they to the channel
func NewMonitor(WG *sync.WaitGroup, duration time.Duration, chanmonitor chan inst.Monitor, chansync chan string) {
	var m inst.Monitor
	var rtm runtime.MemStats

	WG.Add(1)
	defer WG.Done()

	m.Cmetrics = make([]inst.Counter, len(inst.Cmetricnames))
	m.Gmetrics = make([]inst.Gauge, len(inst.Gmetricnames))

	for {
		select {
		case t := <-time.After(duration):
			fmt.Printf("duration event time:%v", t)
		case sig := <-chansync:
			fmt.Printf("sync signal:%v", sig)
			return
		}

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

// NewMonitorGopsutil begin collect metrics infinitly and send they to the channel
func NewMonitorGopsutil(WG *sync.WaitGroup, duration time.Duration, chanmonitor chan inst.Monitor, chansync chan string) {
	var m inst.Monitor

	WG.Add(1)
	defer WG.Done()

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
		select {
		case t := <-time.After(duration):
			fmt.Printf("duration event time:%v", t)
		case sig := <-chansync:
			fmt.Printf("sync signal:%v", sig)
			return
		}

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

		inst.Sugar.Debugf("NewMonitorGopsutil: m=%v send it to chanmonitor", m)

		// Send new collected data to the channel
		chanmonitor <- m

	}
}

// SendMetrics send collected metrics to the web API
func SendMetrics(m inst.Monitor) {

	var mc crypt.MetricsCrypt

	// internal.Gauge type send
	for key, element := range inst.Gmetricnames {

		var v inst.Metrics
		var mv proto.Metric

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

		if inst.ServerAddress != "" {
			err := SendHttpMetric(v)
			if err != nil {
				inst.Sugar.Debug(err)
			}
		}

		if inst.GrpcSrvAddr == "" {
			break
		}
		// send gRPC request
		mv.ID = v.ID
		mv.MType = v.MType
		mv.Value = *v.Value
		mv.Hash = v.Hash

		SendGrpcMetric(&mv)

	}

	// counter type send
	for key, element := range inst.Cmetricnames {

		var v inst.Metrics
		var mv proto.Metric

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

		if inst.ServerAddress != "" {
			err := SendHttpMetric(v)
			if err != nil {
				inst.Sugar.Debug(err)
			}
		}

		if inst.GrpcSrvAddr == "" {
			break
		}
		// send gRPC request
		mv.ID = v.ID
		mv.MType = v.MType
		mv.Delta = *v.Delta
		mv.Hash = v.Hash

		SendGrpcMetric(&mv)
	}
}

func SendHttpMetric(v inst.Metrics) error {
	body, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	var url = fmt.Sprintf("http://%v/update/",
		inst.ServerAddress)

	if inst.PublicKeyFileName != "" {
		AgentEncrypter.Init()
		body, _ = AgentEncrypter.EncryptBytes(body)
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {

		inst.Sugar.Error(err.Error())
		return err
	}
	request.Header.Set("Content-Type", inst.ContentType)

	request.Header.Set("X-Real-IP", inst.UseIp)

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {

		inst.Sugar.Error(err.Error())
		return err
	}
	defer resp.Body.Close()
	inst.Sugar.Debug(resp.Status)
	return nil
}

// RunSendMetrics begin waiting metrics from channel and send they to the web APIs
func RunSendMetrics(WG *sync.WaitGroup, duration time.Duration, chanmonitor1 chan inst.Monitor, chanmonitor2 chan inst.Monitor) {

	inst.Sugar.Infow("Agent started gorutine for send metrics")

	var M, m inst.Monitor
	var err bool

	WG.Add(1)
	defer WG.Done()

	M.Cmetrics = make([]inst.Counter, len(inst.Cmetricnames))
	M.Gmetrics = make([]inst.Gauge, len(inst.Gmetricnames))

	slipcount := 0

	for {
		<-time.After(duration)

		// read first chan
		c1 := len(chanmonitor1)
		c2 := len(chanmonitor2)

		if c1 == 0 && c2 == 0 {
			slipcount++
		}
		if slipcount > 3 {
			inst.Sugar.Infof("runSendMetrics -> 3 times empty channel, exit %v", slipcount)
			return
		}

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
				// add only first FirstGoRutineVariables metrics, they sended in first go rutine
				CopyPartSliceG(M.Gmetrics, m.Gmetrics, 0, inst.FirstGoRutineVariables)
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
				// add only last from FirstGoRutineVariables metrics, they sended in second go rutine
				CopyPartSliceG(M.Gmetrics, m.Gmetrics, inst.FirstGoRutineVariables, len(M.Gmetrics))
			}
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

	if inst.ServerAddress != "" {
		err := SendHttpBatch(metricsb)
		if err != nil {
			inst.Sugar.Infow(err.Error())
			return
		}
	}

	if inst.GrpcSrvAddr == "" {
		return
	}

	err := SendGrpcBatch(metricsb)
	if err != nil {
		inst.Sugar.Infow(err.Error())
		return
	}

}

func SendHttpBatch(metricsb []inst.Metrics) error {
	body, err := json.Marshal(metricsb)
	if err != nil {
		log.Fatal(err)
	}

	var url = fmt.Sprintf("http://%v/updates/",
		inst.ServerAddress)
	inst.Sugar.Infow(url)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {

		inst.Sugar.Infow(err.Error())
	}
	request.Header.Set("Content-Type", inst.ContentType)

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {

		inst.Sugar.Infow(err.Error())
		return err
	}
	defer resp.Body.Close()
	inst.Sugar.Infow(resp.Status)
	return nil
}

func SendGrpcMetric(m *proto.Metric) {

	if GrpcUsersClient == nil {
		// устанавливаем соединение с сервером
		conn, err := grpc.Dial(inst.GrpcSrvAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			inst.Sugar.Fatal(err)
		}
		// defer conn.Close()

		// получаем переменную интерфейсного типа UsersClient,
		// через которую будем отправлять сообщения
		c := proto.NewUsersClient(conn)
		GrpcUsersClient = &c
	}

	ret, err := (*GrpcUsersClient).UpdateMetric(context.Background(), m)
	if err != nil {
		inst.Sugar.Fatal(err)
	}
	inst.Sugar.Debug(ret.GetText())
}

func SendGrpcBatch(metricsb []inst.Metrics) error {

	if GrpcUsersClient == nil {
		// устанавливаем соединение с сервером
		conn, err := grpc.Dial(inst.GrpcSrvAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			inst.Sugar.Debug(err)
			return err
		}
		// defer conn.Close()

		// получаем переменную интерфейсного типа UsersClient,
		// через которую будем отправлять сообщения
		c := proto.NewUsersClient(conn)
		GrpcUsersClient = &c
	}

	var mb proto.BatchMetrics

	for i := range metricsb {

		var pm = proto.Metric{
			ID:    metricsb[i].ID,
			MType: metricsb[i].MType,
			Hash:  metricsb[i].Hash,
		}

		if pm.MType == "counter" {
			pm.Delta = *metricsb[i].Delta
		}
		if pm.MType == "gauge" {
			pm.Value = *metricsb[i].Value
		}

		mb.Metrics = append(mb.Metrics, &pm)
		mb.Count = int32(i)
	}

	ret, err := (*GrpcUsersClient).UpdateBatchMetrics(context.Background(), &mb)
	if err != nil {
		inst.Sugar.Debug(err)
		return err
	}
	inst.Sugar.Debug(ret.GetText())
	return nil
}
