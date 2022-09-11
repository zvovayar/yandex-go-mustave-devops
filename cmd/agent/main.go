// main for agent
// parameters:
// environment variables
// ADDRESS=192.168.0.23:44751
// REPORT_INTERVAL=1s
// POLL_INTERVAL=1s
// KEY=/tmp/zK3deX0
// flags:
// -a address to bind on
// -k key for hash calculate
// -r report interval
// -p poll interval
// -B batch send data
// Agent collect metrics and send to the server
package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/zvovayar/yandex-go-mustave-devops/internal/agent"
	"github.com/zvovayar/yandex-go-mustave-devops/internal/config"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"

	"go.uber.org/zap"
)

const pprofAddr = ":8082"

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {

	var WG sync.WaitGroup
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version:%s\n", buildVersion)
	fmt.Printf("Build date:%s\n", buildDate)
	fmt.Printf("Build commit:%s\n", buildCommit)

	inst.Sugar = zap.NewExample().Sugar()

	defer inst.Sugar.Sync()

	config.ConfigAgentInit()

	chanOS := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(chanOS, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	chanm := make(chan inst.Monitor, inst.BufferLength)
	chanmGopsutil := make(chan inst.Monitor, inst.BufferLength)

	chanSync := make(chan string, 1)

	go agent.NewMonitor(&WG, inst.PollInterval, chanm, chanSync)
	go agent.NewMonitorGopsutil(&WG, inst.PollInterval, chanmGopsutil, chanSync)

	go agent.RunSendMetrics(&WG, inst.ReportInterval, chanm, chanmGopsutil)

	sig := <-chanOS
	inst.Sugar.Infof("INFO got a signal '%v', start shutting down...", sig) // put breakpoint here

	close(chanSync)

	WG.Wait()

	inst.Sugar.Infow("Shutdown complete")

}
