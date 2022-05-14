package storage

//var StoreMonitor Monitor

type Storage interface {
	GetGMvalue(gmname string) Gauge
	GetCMvalue(cmname string) Counter
	SetGMvalue(gmname string, gm Gauge)
	SetCMvalue(cmname string, cm Counter)
}

type PersistanceStorage interface {
	Storage
	NewPersistanceStorage(s *Storage) (*PersistanceStorage, error)
	ClosePersistanceStorage() error
}
