package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zvovayar/yandex-go-mustave-devops/internal/agent"
	"github.com/zvovayar/yandex-go-mustave-devops/internal/config"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

func main() {

	config.ConfigAgentInit()

	chanm := make(chan inst.Monitor, inst.BufferLength)
	chanOS := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(chanOS, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	go agent.NewMonitor(inst.PollInterval, chanm)
	go agent.RunSendMetrics(inst.ReportInterval, chanm)

	sig := <-chanOS
	log.Printf("INFO got a signal '%v', start shutting down...\n", sig) // put breakpoint here
	log.Printf("Shutdown complete")
}
