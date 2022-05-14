package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"time"
)

// var chanPStoreMem chan StoreMem

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

func (mps *MemPStorage) InitMemPStorage(ch chan StoreMem) chan StoreMem {
	// ch := make(chan StoreMem, BufferLength)
	mps.chanPStoreMem = ch
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

		// c := len(chanPStoreMem)
		log.Printf("Save data to the file %v\n", StoreFile)
		// for i := 0; i < c; i++ {

		// m, ok := <-chanPStoreMem
		// if !ok {
		// 	log.Println(err)
		// 	break
		// }

		data, err = json.Marshal(mps.sm.monitor)
		if err != nil {
			return err
		}

		log.Printf("Save data object: %v", mps.sm.monitor)
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

		// }
		// записываем буфер в файл
		if err := mps.writer.Flush(); err != nil {
			log.Fatal(err)
			return err
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

	var g = mps.sm.GetGMvalue(gmname)
	return g
}

// mirror StoreMem interface
func (mps *MemPStorage) GetCMvalue(cmname string) Counter {

	var c = mps.sm.GetCMvalue(cmname)
	return c
}

// mirror StoreMem interface + persistance function
func (mps *MemPStorage) SetGMvalue(gmname string, gm Gauge) {
	mps.sm.SetGMvalue(gmname, gm)
	// mps.chanPStoreMem = make(chan StoreMem, BufferLength)
	// chanPStoreMem <- mps.sm
	// log.Printf("chanPStoreMem <-: %v", mps.sm.monitor)

}

// mirror StoreMem interface + persistance function
func (mps *MemPStorage) SetCMvalue(cmname string, cm Counter) {
	mps.sm.SetCMvalue(cmname, cm)

	// chanPStoreMem <- mps.sm
	// log.Printf("chanPStoreMem <-: %v", mps.sm.monitor)
}

func (mps *MemPStorage) LoadData() {

	//
	var err error
	var data []byte
	var file *os.File
	var scanner *bufio.Scanner

	file, err = os.OpenFile(StoreFile, os.O_RDONLY, 0777)
	if err != nil {
		log.Println(err)
		return
	}

	scanner = bufio.NewScanner(file)

	// одиночное сканирование до следующей строки
	for scanner.Scan() {
		// читаем данные из scanner
		data = scanner.Bytes()

		err := json.Unmarshal(data, &mps.sm.monitor)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Load data: %v", mps.sm.monitor)
	}
	log.Println(scanner.Err())

}
