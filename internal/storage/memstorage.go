package storage

type StoreMem struct {
	monitor Monitor
}

type StoreMemInt interface {
	GetMonitor() *Monitor
	SetMonitor(m *Monitor)
}

func (sm *StoreMem) GetMonitor() *Monitor {
	Sugar.Infow("func (sm *StoreMem) GetMonitor() *Monitor ")
	return &sm.monitor
}

func (sm *StoreMem) GetGMvalue(gmname string) Gauge {
	if _, ok := Gmetricnames[gmname]; !ok {
		// не нашли название метрики, были ошибки
		Sugar.Infof("GetGMvalue Error Can not find gauge metric name:%v", gmname)
		return 0
	}
	return sm.monitor.Gmetrics[Gmetricnames[gmname]]
}

func (sm *StoreMem) GetCMvalue(cmname string) Counter {
	if _, ok := Cmetricnames[cmname]; !ok {
		// не нашли название метрики, были ошибки
		Sugar.Infof("GetCMvalue Error Can not find counter metric name:%v", cmname)
		return 0
	}
	return sm.monitor.Cmetrics[Cmetricnames[cmname]]
}

func (sm *StoreMem) SetGMvalue(gmname string, gm Gauge) {

	if _, ok := Gmetricnames[gmname]; !ok {
		// не нашли название метрики, были ошибки
		Sugar.Infof("SetGMvalue Can not find gauge metric, add new name:%v", gmname)

		// сделать добавление новой метрики
		Gmetricnames[gmname] = len(Gmetricnames)
		StoreMonitor.GetMonitor().Gmetrics = append(StoreMonitor.GetMonitor().Gmetrics, gm)
		// sm.monitor.Gmetrics = append(sm.monitor.Gmetrics, gm)
		return
	}
	sm.monitor.Gmetrics[Gmetricnames[gmname]] = gm
}

func (sm *StoreMem) SetCMvalue(cmname string, cm Counter) {
	if _, ok := Cmetricnames[cmname]; !ok {
		// не нашли название метрики, были ошибки
		Sugar.Infof("SetCMvalue Can not find counter metric, add new name:%v", cmname)

		// сделать добавление новой метрики
		Cmetricnames[cmname] = len(Cmetricnames)
		StoreMonitor.GetMonitor().Cmetrics = append(StoreMonitor.GetMonitor().Cmetrics, cm)
		// sm.monitor.Cmetrics = append(sm.monitor.Cmetrics, cm)
		return
	}
	sm.monitor.Cmetrics[Cmetricnames[cmname]] += cm
}
