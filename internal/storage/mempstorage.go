package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"time"
)

type MemPStorage struct {
	sm            StoreMem
	chanPStoreMem chan StoreMem
	file          *os.File // файл для записи
	writer        *bufio.Writer
}

func (mps *MemPStorage) GetMonitor() *Monitor {
	log.Printf("func (mps MemPStorage) GetMonitor() *Monitor ")
	return &(mps.sm.monitor)
}

func (mps *MemPStorage) InitMemPStorage() chan StoreMem {
	mps.chanPStoreMem = make(chan StoreMem, BufferLength)
	return mps.chanPStoreMem
}

func (mps *MemPStorage) NewPersistanceStorage() error {

	//
	// TODO open file and start gorutine saving data
	//
	var err error
	var data []byte

	mps.file, err = os.OpenFile(StoreFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
		return err
	}

	mps.writer = bufio.NewWriter(mps.file)

	// infinity for{} flash data to file
	for {
		<-time.After(StoreInterval)

		c := len(mps.chanPStoreMem)
		log.Printf("Flash data to the file %v\n", StoreFile)
		for i := 0; i < c; i++ {

			m, ok := <-mps.chanPStoreMem
			if !ok {
				log.Println(err)
				break
			}

			data, err = json.Marshal(m.monitor)
			if err != nil {
				return err
			}

			// записываем в буфер
			if _, err := mps.writer.Write(data); err != nil {
				log.Fatal(err)
				return err
			}

			// добавляем перенос строки
			if err := mps.writer.WriteByte('\n'); err != nil {
				log.Fatal(err)
				return err
			}

			// записываем буфер в файл
			if err := mps.writer.Flush(); err != nil {
				log.Fatal(err)
				return err
			}
		}
	}

	// return nil
}

func (mps *MemPStorage) ClosePersistanceStorage() error {

	//
	// TODO close file
	//

	return mps.file.Close()
}

// mirror StoreMem interface
func (mps *MemPStorage) GetGMvalue(gmname string) Gauge {
	var g Gauge = mps.sm.GetGMvalue(gmname)
	return g
}

// mirror StoreMem interface
func (mps *MemPStorage) GetCMvalue(cmname string) Counter {
	var c Counter = mps.sm.GetCMvalue(cmname)
	return c
}

// mirror StoreMem interface + persistance function
func (mps *MemPStorage) SetGMvalue(gmname string, gm Gauge) {
	mps.sm.SetGMvalue(gmname, gm)

	mps.chanPStoreMem <- mps.sm
}

// mirror StoreMem interface + persistance function
func (mps *MemPStorage) SetCMvalue(cmname string, cm Counter) {
	mps.sm.SetCMvalue(cmname, cm)

	mps.chanPStoreMem <- mps.sm
}

func init() {
	StoreMonitor.GetMonitor().Cmetrics = make([]Counter, len(Cmetricnames))
	StoreMonitor.GetMonitor().Gmetrics = make([]Gauge, len(Gmetricnames))
	log.Println("StoreMonitor.GetMonitor().Cmetrics = make([]Counter, len(Cmetricnames))")
	log.Printf("%v", Cmetricnames)
	log.Printf("%v", Gmetricnames)
	log.Printf("%+v", StoreMonitor.GetMonitor())
}
