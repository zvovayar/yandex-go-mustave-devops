package storage

var StoreMonitor StoreMem

type StoreMem struct {
	monitor Monitor
}

func (sm *StoreMem) GetGMvalue(gmname string) Gauge {
	return sm.monitor.Gmetrics[Gmetricnames[gmname]]
}

func (sm *StoreMem) GetCMvalue(cmname string) Counter {
	return sm.monitor.Cmetrics[Cmetricnames[cmname]]
}

func (sm *StoreMem) SetGMvalue(gmname string, gm Gauge) {
	sm.monitor.Gmetrics[Gmetricnames[gmname]] = gm
}

func (sm *StoreMem) SetCMvalue(cmname string, cm Counter) {
	sm.monitor.Cmetrics[Cmetricnames[cmname]] += cm
}
