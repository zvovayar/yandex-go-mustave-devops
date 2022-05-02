package storage

import (
	"fmt"
)

var StoreMonitor Monitor

type Storage interface {
	GetMonitor() *Monitor
}

type StorageStruct struct {
	ims *Monitor
}

func (im StorageStruct) GetMonitor() *Monitor {
	fmt.Printf("h: %v\n", im)
	im.ims = &StoreMonitor
	return &StoreMonitor
}
