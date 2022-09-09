// Package config fill config vars from environment and command flags
package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

type AgentConfig struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
}

func ConfigAgentInit() {
	var cfg AgentConfig
	var cfgFromFlags AgentConfig
	// загрузим переменные среды
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	inst.Sugar.Infow("Agent started")
	inst.Sugar.Infof("Config environment:%+v", cfg)

	// load flags
	flag.StringVar(&cfgFromFlags.Address, "a", inst.ServerAddress, "address to bind on")
	flag.StringVar(&cfgFromFlags.Key, "k", "", "key for hash calculate")
	flag.DurationVar(&cfgFromFlags.ReportInterval, "r", inst.ReportInterval, "report interval")
	flag.DurationVar(&cfgFromFlags.PollInterval, "p", inst.PollInterval, "poll interval")
	flag.BoolVar(&inst.BatchSend, "B", true, "batch send data")
	flag.Parse()
	inst.Sugar.Infof("Agent Config flags:%+v", cfgFromFlags)

	// assign work parameters
	if cfg.Address != "" {
		inst.ServerAddress = cfg.Address
	} else {
		inst.ServerAddress = cfgFromFlags.Address
	}
	if cfg.Key != "" {
		inst.Key = cfg.Key
	} else {
		inst.Key = cfgFromFlags.Key
	}
	if cfg.PollInterval > 0 {
		inst.PollInterval = cfg.PollInterval
	} else {
		inst.PollInterval = cfgFromFlags.PollInterval
	}
	if cfg.ReportInterval > 0 {
		inst.ReportInterval = cfg.ReportInterval
	} else {
		inst.ReportInterval = cfgFromFlags.ReportInterval
	}

	inst.Sugar.Infof("Agent Strated with variables: address=%v, poll interval=%v, report interval=%v, key=%v, batch send=%v",
		inst.ServerAddress, inst.PollInterval, inst.ReportInterval, inst.Key, inst.BatchSend)
}
