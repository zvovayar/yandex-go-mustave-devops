package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// var chanPStoreMem chan StoreMem

type MemSQLStorage struct {
	sm            StoreMem
	chanPStoreMem chan StoreMem
	DatabaseDSN   string
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
	// TODO: open database
	//

	// infinity loop for{} flash data to file
	for {
		<-time.After(StoreInterval)

	}

	// return nil
}

func (mps *MemSQLStorage) PingSQLserver(ctx context.Context) error {

	log.Print(sql.Drivers())
	db, err := sql.Open("postgres", DatabaseDSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// работаем с базой
	// ...
	// можем продиагностировать соединение
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		log.Println(err)
		return err
	}

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
}

// mirror StoreMem interface + persistance function
func (mps *MemSQLStorage) SetCMvalue(cmname string, cm Counter) {
	mps.sm.SetCMvalue(cmname, cm)
}

func (mps *MemSQLStorage) LoadData() {
}
