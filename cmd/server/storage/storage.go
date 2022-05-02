package storage

import (
	"fmt"

	"github.com/zvovayar/yandex-go-mustave-devops/internal"
)

var StoreMonitor internal.Monitor

type Storage interface {
	GetMonitor() *internal.Monitor
}

type StorageStruct struct {
	ims *internal.Monitor
}

func (im StorageStruct) GetMonitor() *internal.Monitor {
	fmt.Printf("h: %v\n", im)
	im.ims = &StoreMonitor
	return &StoreMonitor
}
