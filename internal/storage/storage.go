package storage

import (
	"log"
	"reflect"
)

type Storage interface {
	GetGMvalue(gmname string) Gauge
	GetCMvalue(cmname string) Counter
	SetGMvalue(gmname string, gm Gauge)
	SetCMvalue(cmname string, cm Counter)
}

type PersistanceStorage interface {
	Storage
	NewPersistanceStorage() error
	ClosePersistanceStorage() error
}

// var StoreMonitor StoreMem

var StoreMonitor MemPStorage

func init() {
	StoreMonitor.GetMonitor().Cmetrics = make([]Counter, len(Cmetricnames))
	StoreMonitor.GetMonitor().Gmetrics = make([]Gauge, len(Gmetricnames))
	log.Println("StoreMonitor.GetMonitor().Cmetrics = make([]Counter, len(Cmetricnames))")
	log.Printf("%v", Cmetricnames)
	log.Printf("%v", Gmetricnames)
	log.Printf("%+v", StoreMonitor.GetMonitor())

	log.Printf("Storage type:%s", reflect.TypeOf(StoreMonitor).Name())
	// if reflect.TypeOf(StoreMonitor).Name() == "MemStorage" {
	// 	StoreMonitor.InitMemPStorage(make(chan StoreMem, BufferLength))
	// }
	// chanPStoreMem = make(chan StoreMem, BufferLength)
}
