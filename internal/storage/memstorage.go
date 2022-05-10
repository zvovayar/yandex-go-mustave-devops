package storage

import (
	"log"
)

var StoreMonitor StoreMem

type StoreMem struct {
	monitor Monitor
}

func (sm *StoreMem) GetGMvalue(gmname string) Gauge {
	if _, ok := Gmetricnames[gmname]; !ok {
		// не нашли название метрики, были ошибки
		log.Printf("Error Can not find gauge metric name:%v", gmname)
		return 0
	}
	return sm.monitor.Gmetrics[Gmetricnames[gmname]]
}

func (sm *StoreMem) GetCMvalue(cmname string) Counter {
	if _, ok := Cmetricnames[cmname]; !ok {
		// не нашли название метрики, были ошибки
		log.Printf("Error Can not find counter metric name:%v", cmname)
		return 0
	}
	return sm.monitor.Cmetrics[Cmetricnames[cmname]]
}

func (sm *StoreMem) SetGMvalue(gmname string, gm Gauge) {

	if _, ok := Gmetricnames[gmname]; !ok {
		// не нашли название метрики, были ошибки
		log.Printf("Error Can not find gauge metric name:%v", gmname)
		return
	}
	sm.monitor.Gmetrics[Gmetricnames[gmname]] = gm
}

func (sm *StoreMem) SetCMvalue(cmname string, cm Counter) {
	if _, ok := Cmetricnames[cmname]; !ok {
		// не нашли название метрики, были ошибки
		log.Printf("Error Can not find counter metric name:%v", cmname)
		return
	}
	sm.monitor.Cmetrics[Cmetricnames[cmname]] += cm
}
