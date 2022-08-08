package storage

import (
	"bufio"
	"context"
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
	// inst.Sugar.Infof("func (mps MemPStorage) GetMonitor() *Monitor ")
	return &(mps.sm.monitor)
}

func (mps *MemPStorage) InitMemPStorage(ch chan StoreMem) chan StoreMem {
	// ch := make(chan StoreMem, BufferLength)
	mps.chanPStoreMem = ch
	return mps.chanPStoreMem
}

func (mps *MemPStorage) NewPersistanceStorage() error {

	//
	// open file and start gorutine saving data
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

		// inst.Sugar.Infof("Save data to the file %v\n", StoreFile)

		//
		// save names maps
		//
		data, err = json.Marshal(Gmetricnames)
		if err != nil {
			return err
		}

		// inst.Sugar.Infof("Save data object: %v", Gmetricnames)
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

		data, err = json.Marshal(Cmetricnames)
		if err != nil {
			return err
		}

		// inst.Sugar.Infof("Save data object: %v", Cmetricnames)
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

		//
		// save data
		//
		data, err = json.Marshal(mps.sm.monitor)
		if err != nil {
			return err
		}

		// inst.Sugar.Infof("Save data object: %v", mps.sm.monitor)
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
	// close file
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
	// inst.Sugar.Infof("chanPStoreMem <-: %v", mps.sm.monitor)

}

// mirror StoreMem interface + persistance function
func (mps *MemPStorage) SetCMvalue(cmname string, cm Counter) {
	mps.sm.SetCMvalue(cmname, cm)

	// inst.Sugar.Infof("chanPStoreMem <-: %v", mps.sm.monitor)
}

func (mps *MemPStorage) LoadData() {

	//
	var err error
	var data []byte
	var file *os.File
	var scanner *bufio.Scanner

	file, err = os.OpenFile(StoreFile, os.O_RDONLY, 0777)
	if err != nil {
		Sugar.Infow(err.Error())
		return
	}

	scanner = bufio.NewScanner(file)

	// одиночное сканирование до следующей строки
	for scanner.Scan() {
		// metrics names читаем данные из scanner
		data = scanner.Bytes()

		err := json.Unmarshal(data, &Gmetricnames)
		if err != nil {
			log.Fatal(err)
		}

		// inst.Sugar.Infof("Load Gmetricnames: %v", Gmetricnames)

		scanner.Scan()
		data = scanner.Bytes()

		err = json.Unmarshal(data, &Cmetricnames)
		if err != nil {
			log.Fatal(err)
		}

		// inst.Sugar.Infof("Load Cmetricnames: %v", Cmetricnames)

		// читаем данные из scanner
		scanner.Scan()
		data = scanner.Bytes()

		err = json.Unmarshal(data, &mps.sm.monitor)
		if err != nil {
			log.Fatal(err)
		}

		// inst.Sugar.Infof("Load data: %v", mps.sm.monitor)
	}
	// inst.Sugar.Infow(scanner.Err())

}

func (mps *MemPStorage) SaveBatch(ctx context.Context, batchM []Metrics) error {

	c := len(batchM)
	for i := 0; i < c; i++ {
		if batchM[i].MType == "gauge" {
			mps.sm.SetGMvalue(batchM[i].ID, Gauge(*batchM[i].Value))
		} else if batchM[i].MType == "counter" {
			mps.sm.SetCMvalue(batchM[i].ID, Counter(*batchM[i].Delta))
		}
	}

	return nil
}

func (mps *MemPStorage) PingSQLserver(ctx context.Context) error {
	return nil
}
