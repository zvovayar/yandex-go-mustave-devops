package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	// log.Printf("func (mps MemSQLStorage) GetMonitor() *Monitor ")
	return &(mps.sm.monitor)
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
	if err := mps.CheckAndCreateMDatabase(context.Background(), DatabaseDSN); err != nil {
		panic(err)
	}

	var err error
	mps.db, err = sqlx.Open("postgres", DatabaseDSN)
	if err != nil {
		panic(err)
	}
	defer mps.db.Close()

	// infinity loop for{} save data to SQL database
	for {
		<-time.After(StoreInterval)

		// for key := range Gmetricnames {

		// 	log.Printf("INSERT INTO gmetrics (gauge, name ) VALUES(%f, '%v')", mps.sm.GetGMvalue(key), key)
		// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// 	defer cancel()
		// 	_, err = mps.db.ExecContext(ctx,
		// 		"INSERT INTO gmetrics (gauge, name ) VALUES($1, $2)", mps.sm.GetGMvalue(key), key)

		// 	if err != nil {
		// 		log.Println("NewPersistanceStorage " + err.Error())
		// 		return err
		// 	}
		// 	log.Printf("NewPersistanceStorage value %v=%f saved", key, mps.sm.GetGMvalue(key))
		// }

		// for key := range Cmetricnames {

		// 	log.Printf("INSERT INTO cmetrics (counter, name ) VALUES(%d, '%v')", mps.sm.GetCMvalue(key), key)
		// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// 	defer cancel()
		// 	_, err = mps.db.ExecContext(ctx,
		// 		"INSERT INTO cmetrics (counter, name ) VALUES($1, $2)", mps.sm.GetCMvalue(key), key)

		// 	if err != nil {
		// 		log.Println("NewPersistanceStorage " + err.Error())
		// 		return err
		// 	}
		// 	log.Printf("NewPersistanceStorage value %v=%d saved", key, mps.sm.GetCMvalue(key))
		// }
	}

	// return nil
}

func (mps *MemSQLStorage) PingSQLserver(ctx context.Context) error {

	log.Printf("SQL drivers available:%v", sql.Drivers())
	db, err := sql.Open("postgres", DatabaseDSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("PingSQLserver success DSN=%v", DatabaseDSN)
	return nil
}

func (mps *MemSQLStorage) ClosePersistanceStorage() error {

	//
	// TODO close database
	//

	return fmt.Errorf("MemSQLStorage.ClosePersistanceStorage not implemented")
}

// mirror StoreMem interface
func (mps *MemSQLStorage) GetGMvalue(gmname string) Gauge {

	var g = mps.sm.GetGMvalue(gmname)
	return g
}

// mirror StoreMem interface
func (mps *MemSQLStorage) GetCMvalue(cmname string) Counter {

	var c = mps.sm.GetCMvalue(cmname)
	return c
}

// mirror StoreMem interface + persistance function
func (mps *MemSQLStorage) SetGMvalue(gmname string, gm Gauge) {
	mps.sm.SetGMvalue(gmname, gm)

	var err error
	log.Printf("INSERT INTO gmetrics (gauge, name ) VALUES(%f, '%v')", mps.sm.GetGMvalue(gmname), gmname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = mps.db.ExecContext(ctx,
		"INSERT INTO gmetrics (gauge, name ) VALUES($1, $2)", mps.sm.GetGMvalue(gmname), gmname)

	if err != nil {
		log.Println("NewPersistanceStorage " + err.Error())
		return
	}
	log.Printf("NewPersistanceStorage value %v=%f saved", gmname, mps.sm.GetGMvalue(gmname))
}

// mirror StoreMem interface + persistance function
func (mps *MemSQLStorage) SetCMvalue(cmname string, cm Counter) {
	mps.sm.SetCMvalue(cmname, cm)
	var err error
	log.Printf("INSERT INTO cmetrics (counter, name ) VALUES(%d, '%v')", mps.sm.GetCMvalue(cmname), cmname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = mps.db.ExecContext(ctx,
		"INSERT INTO cmetrics (counter, name ) VALUES($1, $2)", mps.sm.GetCMvalue(cmname), cmname)

	if err != nil {
		log.Println("NewPersistanceStorage " + err.Error())
		return
	}
	log.Printf("NewPersistanceStorage value %v=%d saved", cmname, mps.sm.GetCMvalue(cmname))
}

func (mps *MemSQLStorage) LoadData() {
	//
	// TODO load data from SQL database
	//
}

func (mps *MemSQLStorage) CheckAndCreateMDatabase(ctx context.Context, DSN string) error {

	err := mps.PingSQLserver(ctx)
	if err != nil {
		return err
	}

	//
	// create database and tables
	//
	db, err := sql.Open("postgres", DatabaseDSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx,
		"CREATE TABLE IF NOT EXISTS gmetrics (id BIGSERIAL, gauge NUMERIC, name VARCHAR(50))")

	if err != nil {
		log.Println("CheckAndCreateMDatabase " + err.Error())
		return err
	}
	log.Printf("CheckAndCreateMDatabase table gmetrics created")

	_, err = db.ExecContext(ctx,
		"CREATE TABLE IF NOT EXISTS cmetrics (id BIGSERIAL, counter BIGINT, name VARCHAR(50))")

	if err != nil {
		log.Println("CheckAndCreateMDatabase " + err.Error())
		return err
	}
	log.Printf("CheckAndCreateMDatabase table cmetrics created")

	_, err = db.ExecContext(ctx,
		"CREATE TABLE IF NOT EXISTS metrics (idrec BIGSERIAL, id VARCHAR(50), mtype VARCHAR(50), delta BIGINT, value NUMERIC)")

	if err != nil {
		log.Println("CheckAndCreateMDatabase " + err.Error())
		return err
	}
	log.Printf("CheckAndCreateMDatabase table metrics created")

	return nil
}

func (mps *MemSQLStorage) SaveBatch(ctx context.Context, batchM []Metrics) error {
	//
	// use: sqlx
	// cars := []Cars{
	// 	{},
	// }
	if _, err := mps.db.NamedExec(`INSERT INTO metrics (id, mtype, delta, value)
			VALUES (:id, :mtype, :delta, :value)`, batchM); err != nil {
		log.Println("SaveBatch " + err.Error())
		return err
	}

	return nil
}
