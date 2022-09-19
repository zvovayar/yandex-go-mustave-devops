package storage

import (
	"context"
	"reflect"

	"go.uber.org/zap"
)

type Storage interface {
	GetGMvalue(gmname string) Gauge
	GetCMvalue(cmname string) Counter
	SetGMvalue(gmname string, gm Gauge)
	SetCMvalue(cmname string, cm Counter)
	SaveBatch(ctx context.Context, batchM []Metrics) error
}

type PersistanceStorage interface {
	Storage
	NewPersistanceStorage() error
	ClosePersistanceStorage() error
}

//
// Uncomment if you want memory storage
//
// var StoreMonitor StoreMem

//
// Uncomment if you want text/json file storage
//
// var StoreMonitor MemPStorage

// StoreMonitor main global SQL databse storage
var StoreMonitor MemSQLStorage

func init() {

	// Sugar = zap.NewExample().Sugar()
	// Sugar = zap.NewNop().Sugar().With(zap.NewDevelopmentConfig())
	Sugar = (*zap.SugaredLogger)(zap.NewExample().Sugar().With(zap.NewProductionConfig()))
	defer Sugar.Sync()

	StoreMonitor.GetMonitor().Cmetrics = make([]Counter, len(Cmetricnames))
	StoreMonitor.GetMonitor().Gmetrics = make([]Gauge, len(Gmetricnames))

	Sugar.Infof("Storage init")

	Sugar.Infof("%v", Cmetricnames)
	Sugar.Infof("%v", Gmetricnames)
	Sugar.Infof("%+v", StoreMonitor.GetMonitor())

	Sugar.Infof("Storage type:%s", reflect.TypeOf(StoreMonitor).Name())
}
