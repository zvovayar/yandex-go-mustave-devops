package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/zvovayar/yandex-go-mustave-devops/internal/agent"
	"github.com/zvovayar/yandex-go-mustave-devops/internal/config"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"

	"go.uber.org/zap"
)

const pprofAddr = ":8082"

func main() {

	inst.Sugar = zap.NewExample().Sugar()

	defer inst.Sugar.Sync()

	config.ConfigAgentInit()

	chanOS := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(chanOS, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	chanm := make(chan inst.Monitor, inst.BufferLength)
	chanmGopsutil := make(chan inst.Monitor, inst.BufferLength)
	go agent.NewMonitor(inst.PollInterval, chanm)
	go agent.NewMonitorGopsutil(inst.PollInterval, chanmGopsutil)

	go agent.RunSendMetrics(inst.ReportInterval, chanm, chanmGopsutil)

	// start profiler
	// go http.ListenAndServe(pprofAddr, nil)
	go func() {
		log.Println(http.ListenAndServe(pprofAddr, nil))
	}()

	sig := <-chanOS
	inst.Sugar.Infof("INFO got a signal '%v', start shutting down...", sig) // put breakpoint here
	inst.Sugar.Infow("Shutdown complete")

}
