package storage

import (
	"time"

	"go.uber.org/zap"
)

type Gauge float64
type Counter int64
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

type Monitor struct {
	Gmetrics []Gauge
	Cmetrics []Counter
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
	"GetSet191":     29,
	"GetSet238":     30,
	"GetSet180":     31,
	"GetSet88":      32,
	"GetSet217":     33,
	"GetSet16":      34,
}
var Cmetricnames = map[string]int{
	"PollCount":    0,
	"testSetGet33": 1,
	"GetSet209":    2,
	"GetSet199":    3,
	"GetSet201":    4,
	"GetSet127":    5,
	"GetSet21":     6,
	"GetSet203":    7,
	"GetSet190":    8,
	"GetSet47":     9,
}
var PollInterval = time.Second * 2    // 2
var ReportInterval = time.Second * 10 //10
var BufferLength = 30
var ServerAddress = "127.0.0.1:8080"
var StoreInterval = time.Second * 300 //300
var StoreFile = "/tmp/devops-metrics-db.json"
var Restore = true
var Key = ""
var DatabaseDSN = ""
var BatchSend = true

//var ContentType = "text/plain"
var ContentType = "application/json"

var Sugar *zap.SugaredLogger
