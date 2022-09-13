package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

// var chanPStoreMem chan StoreMem

type MemSQLStorage struct {
	sm            StoreMem
	chanPStoreMem chan StoreMem
	DatabaseDSN   string
	db            *sqlx.DB
}

func (mps *MemSQLStorage) GetMonitor() *Monitor {
	// inst.Sugar.Infof("func (mps MemSQLStorage) GetMonitor() *Monitor ")
	return &(mps.sm.monitor)
}

func (mps *MemSQLStorage) OpenDB() error {
	var err error
	mps.db, err = sqlx.Open("postgres", DatabaseDSN)
	if err != nil {
		Sugar.Infow(err.Error())
		return err
	}
	return nil
}

func (mps *MemSQLStorage) InitMemSQLStorage(ch chan StoreMem) chan StoreMem {
	// ch := make(chan StoreMem, BufferLength)
	mps.chanPStoreMem = ch
	return mps.chanPStoreMem
}

func (mps *MemSQLStorage) NewPersistanceStorage() error {

	//
	// open databse and start infinity loop
	//

	//
	// open database
	//
	Sugar.Infof("Drivers available: %v", sql.Drivers())
	var err error
	mps.db, err = sqlx.Open("postgres", DatabaseDSN)
	if err != nil {
		Sugar.Infow(err.Error())
		mps.db = nil
	}
	defer mps.db.Close()

	if err := mps.CheckAndCreateMDatabase(context.Background()); err != nil {
		Sugar.Infow(err.Error())
	}

	// infinity loop for{} open database
	for {
		<-time.After(StoreInterval)

		if mps.db == nil {
			mps.db, err = sqlx.Open("postgres", DatabaseDSN)
			if err != nil {
				Sugar.Infow(err.Error())
				mps.db = nil
			}
		}
	}

	// return nil
}

func (mps *MemSQLStorage) PingSQLserver(ctx context.Context) error {

	if mps.db == nil {
		var err error
		mps.db, err = sqlx.Open("postgres", DatabaseDSN)
		if err != nil {
			Sugar.Infow(err.Error())
			mps.db = nil
			return err
		}
		defer mps.db.Close()
	}

	Sugar.Infof("PingSQLserver try for open connections=%v", mps.db.DB.Stats().OpenConnections)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := mps.db.PingContext(ctx); err != nil {
		Sugar.Infow("PingSQLserver PingSQLserver error " + err.Error())
		return err
	}
	Sugar.Infof("PingSQLserver success DSN=%v", DatabaseDSN)
	return nil
}

func (mps *MemSQLStorage) ClosePersistanceStorage() error {

	if mps.db == nil {
		return nil
	}
	if err := mps.db.MustBegin().Commit(); err != nil {
		return fmt.Errorf("MemSQLStorage.ClosePersistanceStorage: %v", err.Error())
	}
	if err := mps.db.Close(); err != nil {
		return fmt.Errorf("MemSQLStorage.ClosePersistanceStorage: %v", err.Error())
	}
	return nil
}

// GetGMvalue mirror StoreMem interface
func (mps *MemSQLStorage) GetGMvalue(gmname string) Gauge {

	var g = mps.sm.GetGMvalue(gmname)
	return g
}

// GetCMvalue mirror StoreMem interface
func (mps *MemSQLStorage) GetCMvalue(cmname string) Counter {

	var c = mps.sm.GetCMvalue(cmname)
	return c
}

// SetGMvalue mirror StoreMem interface + persistance function
func (mps *MemSQLStorage) SetGMvalue(gmname string, gm Gauge) {

	mps.sm.SetGMvalue(gmname, gm)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var m Metrics
	g := float64(mps.sm.GetGMvalue(gmname))
	m.Value = &g
	m.ID = gmname
	m.MType = "gauge"

	// if _, err := mps.db.NamedExec(`INSERT INTO metrics (id, mtype, delta, value)
	if _, err := mps.db.NamedExecContext(ctx, `INSERT INTO metrics (id, mtype, delta, value)
	VALUES (:id, :mtype, :delta, :value)`, m); err != nil {
		Sugar.Infow("SetGMvalue " + err.Error())
		return
	}
	Sugar.Infof("SetGMvalue value %v=%f saved", gmname, mps.sm.GetGMvalue(gmname))
}

// SetCMvalue mirror StoreMem interface + persistance function
func (mps *MemSQLStorage) SetCMvalue(cmname string, cm Counter) {

	mps.sm.SetCMvalue(cmname, cm)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var m Metrics
	c := int64(mps.sm.GetCMvalue(cmname))
	m.Delta = &c
	m.ID = cmname
	m.MType = "counter"

	// if _, err := mps.db.NamedExec(`INSERT INTO metrics (id, mtype, delta, value)
	if _, err := mps.db.NamedExecContext(ctx, `INSERT INTO metrics (id, mtype, delta, value)
		VALUES (:id, :mtype, :delta, :value)`, m); err != nil {
		Sugar.Infow("SetCMvalue " + err.Error())
		return
	}
	Sugar.Infof("SetCMvalue value %v=%d saved", cmname, mps.sm.GetCMvalue(cmname))
}

func (mps *MemSQLStorage) LoadData() {
	//
	// load data from SQL database
	//
	Sugar.Infof("LoadData from DSN=%v", DatabaseDSN)
	db, err := sqlx.Open("postgres", DatabaseDSN)
	if err != nil {
		Sugar.Infow("LoadData " + err.Error())
		return
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var m []Metrics
	err = db.SelectContext(ctx, &m, `select id, mtype, delta, value
						from metrics
						where idrec in (select max(idrec)  
						from metrics
						group by id)`)
	if err != nil {
		Sugar.Infow("LoadData " + err.Error())
		return
	}

	c := len(m)
	for i := 0; i < c; i++ {
		if m[i].MType == "gauge" {
			mps.sm.SetGMvalue(m[i].ID, Gauge(*m[i].Value))
		} else if m[i].MType == "counter" {
			mps.sm.SetCMvalue(m[i].ID, Counter(*m[i].Delta))
		}
	}
	Sugar.Infof("LoadData loaded %d metrics", c)
}

func (mps *MemSQLStorage) CheckAndCreateMDatabase(ctx context.Context) error {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := mps.db.ExecContext(ctx,
		"CREATE TABLE IF NOT EXISTS metrics (idrec BIGSERIAL, id VARCHAR(50), mtype VARCHAR(50), delta BIGINT, value NUMERIC)")

	if err != nil {
		Sugar.Infow("CheckAndCreateMDatabase " + err.Error())
		return err
	}
	Sugar.Infof("CheckAndCreateMDatabase table metrics created")

	return nil
}

func (mps *MemSQLStorage) SaveBatch(ctx context.Context, batchM []Metrics) error {

	c := len(batchM)
	for i := 0; i < c; i++ {
		if batchM[i].MType == "gauge" {
			mps.sm.SetGMvalue(batchM[i].ID, Gauge(*batchM[i].Value))
		} else if batchM[i].MType == "counter" {
			mps.sm.SetCMvalue(batchM[i].ID, Counter(*batchM[i].Delta))
		}
	}

	ctxt, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := mps.db.NamedExecContext(ctxt, `INSERT INTO metrics (id, mtype, delta, value)
			VALUES (:id, :mtype, :delta, :value)`, batchM); err != nil {
		Sugar.Infow("SaveBatch " + err.Error())
		return err
	}

	return nil
}
