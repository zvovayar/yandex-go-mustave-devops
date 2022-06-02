package storage

import (
	"log"
)

type StoreMem struct {
	monitor Monitor
}

type StoreMemInt interface {
	GetMonitor() *Monitor
	SetMonitor(m *Monitor)
}

func (sm *StoreMem) GetMonitor() *Monitor {
	log.Print("func (sm *StoreMem) GetMonitor() *Monitor ")
	return &sm.monitor
}

func (sm *StoreMem) GetGMvalue(gmname string) Gauge {
	if _, ok := Gmetricnames[gmname]; !ok {
		// не нашли название метрики, были ошибки
		log.Printf("GetGMvalue Error Can not find gauge metric name:%v", gmname)
		return 0
	}
	return sm.monitor.Gmetrics[Gmetricnames[gmname]]
}

func (sm *StoreMem) GetCMvalue(cmname string) Counter {
	if _, ok := Cmetricnames[cmname]; !ok {
		// не нашли название метрики, были ошибки
		log.Printf("GetCMvalue Error Can not find counter metric name:%v", cmname)
		return 0
	}
	return sm.monitor.Cmetrics[Cmetricnames[cmname]]
}

func (sm *StoreMem) SetGMvalue(gmname string, gm Gauge) {

	if _, ok := Gmetricnames[gmname]; !ok {
		// не нашли название метрики, были ошибки
		log.Printf("SetGMvalue Can not find gauge metric, add new name:%v", gmname)
		//
		// сделать добавление новой метрики
		//
		Gmetricnames[gmname] = len(Gmetricnames)
		StoreMonitor.GetMonitor().Gmetrics = append(StoreMonitor.GetMonitor().Gmetrics, gm)
		return
	}
	sm.monitor.Gmetrics[Gmetricnames[gmname]] = gm
}

func (sm *StoreMem) SetCMvalue(cmname string, cm Counter) {
	if _, ok := Cmetricnames[cmname]; !ok {
		// не нашли название метрики, были ошибки
		log.Printf("SetCMvalue Can not find counter metric, add new name:%v", cmname)
		//
		// сделать добавление новой метрики
		//
		Cmetricnames[cmname] = len(Cmetricnames)
		StoreMonitor.GetMonitor().Cmetrics = append(StoreMonitor.GetMonitor().Cmetrics, cm)
		return
	}
	// sm.monitor.Cmetrics[Cmetricnames[cmname]] += cm
	sm.monitor.Cmetrics[Cmetricnames[cmname]] = cm
}
